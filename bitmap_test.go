package btmp_test

import (
	"math"
	"testing"

	"github.com/neox5/btmp"
)

/*************************
 * helpers / reference
 *************************/

type refBits struct {
	len  int
	bits []uint8 // 0/1
}

func newRef(n int) *refBits {
	r := &refBits{}
	r.ensure(n)
	return r
}

func (r *refBits) ensure(n int) {
	if n <= r.len {
		return
	}
	old := r.len
	r.len = n
	r.bits = append(r.bits, make([]uint8, n-old)...)
}

func (r *refBits) setRange(start, count int) {
	if count <= 0 {
		return
	}
	end := start + count
	r.ensure(end)
	for i := start; i < end; i++ {
		r.bits[i] = 1
	}
}

func (r *refBits) clearRange(start, count int) {
	if count <= 0 {
		return
	}
	end := start + count
	if start < 0 || end > r.len {
		panic("ref: clear oob")
	}
	for i := start; i < end; i++ {
		r.bits[i] = 0
	}
}

func (r *refBits) copyRange(src *refBits, srcStart, dstStart, count int) {
	if count <= 0 {
		return
	}
	srcEnd := srcStart + count
	if srcStart < 0 || srcEnd > src.len {
		panic("ref: copy source oob")
	}
	dstEnd := dstStart + count
	r.ensure(dstEnd)

	if r == src && rangesOverlap(srcStart, srcEnd, dstStart, dstEnd) {
		if dstStart < srcStart {
			for i := range count {
				r.bits[dstStart+i] = src.bits[srcStart+i]
			}
		} else {
			for i := count - 1; i >= 0; i-- {
				r.bits[dstStart+i] = src.bits[srcStart+i]
			}
		}
		return
	}
	for i := range count {
		r.bits[dstStart+i] = src.bits[srcStart+i]
	}
}

func (r *refBits) at(i int) uint8 { return r.bits[i] }

func (r *refBits) count() int {
	c := 0
	for _, b := range r.bits {
		if b == 1 {
			c++
		}
	}
	return c
}

func (r *refBits) nextSetBit(from int) int {
	if from < 0 || from > r.len {
		panic("ref: nextSetBit from oob")
	}
	for i := from; i < r.len; i++ {
		if r.bits[i] == 1 {
			return i
		}
	}
	return -1
}

func rangesOverlap(a0, a1, b0, b1 int) bool { return a0 < b1 && b0 < a1 }

func eqBitmapRef(t *testing.T, got *btmp.Bitmap, ref *refBits) {
	t.Helper()
	if got.Len() != ref.len {
		t.Fatalf("Len mismatch: got=%d want=%d", got.Len(), ref.len)
	}
	for i := range ref.len {
		gb := got.Test(i)
		rb := ref.at(i) == 1
		if gb != rb {
			t.Fatalf("bit mismatch at %d: got=%v want=%v", i, gb, rb)
		}
	}
	// Tail must be clean beyond Len().
	w := got.Words()
	if got.Len()%64 != 0 && len(w) > 0 {
		last := len(w) - 1
		tailBits := got.Len() & 63
		var mask uint64
		if tailBits == 0 {
			mask = ^uint64(0)
		} else {
			mask = (uint64(1) << uint(tailBits)) - 1
		}
		if (w[last] &^ mask) != 0 {
			t.Fatalf("tail not masked: words[%d]=0x%x mask=0x%x", last, w[last], mask)
		}
	}
}

func mustPanic(t *testing.T, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic")
		}
	}()
	f()
}

/*************************
 * tests
 *************************/

func TestSetClear_BasicAndBoundaries(t *testing.T) {
	t.Parallel()
	cases := []struct {
		preLen     int
		start      int
		count      int
		clearAfter bool
	}{
		{0, 0, 0, true},
		{0, 0, 1, true},
		{0, 1, 1, true},
		{0, 63, 1, true},
		{0, 64, 1, true},
		{0, 65, 1, true},
		{0, 10, 53, true},
		{0, 10, 64, true},
		{0, 10, 65, true},
		{0, 0, 128, true},
		{0, 127, 129, true},
		{256, 5, 5, true},   // entirely in-bounds
		{130, 64, 64, true}, // two full words
	}
	for _, tc := range cases {
		b := btmp.New(0)
		if tc.preLen > 0 {
			b = b.EnsureBits(tc.preLen)
		}
		r := newRef(b.Len())

		b = b.SetRange(tc.start, tc.count)
		r.setRange(tc.start, tc.count)
		r.ensure(tc.start + tc.count)
		eqBitmapRef(t, b, r)

		if tc.clearAfter && tc.start+tc.count <= b.Len() {
			b = b.ClearRange(tc.start, tc.count)
			r.clearRange(tc.start, tc.count)
			eqBitmapRef(t, b, r)
		}
	}
}

func TestCopyRange_Variants(t *testing.T) {
	t.Parallel()

	// Aligned large copy.
	{
		b := btmp.New(0).EnsureBits(4096)
		r := newRef(4096)
		b = b.SetRange(64, 1024)
		r.setRange(64, 1024)

		b = b.CopyRange(b, 64, 2048, 1024)
		r.copyRange(r, 64, 2048, 1024)

		eqBitmapRef(t, b, r)
	}

	// Misaligned copy.
	{
		b := btmp.New(0).EnsureBits(3000)
		r := newRef(3000)
		b = b.SetRange(3, 200).SetRange(777, 129).SetRange(1999, 1)
		r.setRange(3, 200)
		r.setRange(777, 129)
		r.setRange(1999, 1)

		b = b.CopyRange(b, 5, 113, 777)
		r.copyRange(r, 5, 113, 777)

		eqBitmapRef(t, b, r)
	}

	// Overlap forward (dst<src).
	{
		b := btmp.New(0).EnsureBits(2048)
		r := newRef(2048)
		b = b.SetRange(100, 300)
		r.setRange(100, 300)

		b = b.CopyRange(b, 100, 50, 300)
		r.copyRange(r, 100, 50, 300)

		eqBitmapRef(t, b, r)
	}

	// Overlap backward (dst>src).
	{
		b := btmp.New(0).EnsureBits(2048)
		r := newRef(2048)
		b = b.SetRange(200, 400)
		r.setRange(200, 400)

		b = b.CopyRange(b, 200, 350, 400)
		r.copyRange(r, 200, 350, 400)

		eqBitmapRef(t, b, r)
	}
}

func TestNextSetBit_And_Count(t *testing.T) {
	t.Parallel()
	b := btmp.New(0)
	r := newRef(0)

	sets := [][2]int{{2, 1}, {5, 3}, {130, 2}, {512, 200}, {1023, 2}}
	for _, s := range sets {
		b = b.SetRange(s[0], s[1])
		r.setRange(s[0], s[1])
	}
	eqBitmapRef(t, b, r)

	var got []int
	for i := b.NextSetBit(0); i >= 0; i = b.NextSetBit(i + 1) {
		got = append(got, i)
	}
	want := make([]int, 0, r.count())
	for i := r.nextSetBit(0); i >= 0; i = r.nextSetBit(i + 1) {
		want = append(want, i)
	}
	if len(got) != len(want) {
		t.Fatalf("scan len mismatch: got=%d want=%d", len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("scan mismatch at %d: got=%d want=%d", i, got[i], want[i])
		}
	}
	if b.Count() != r.count() {
		t.Fatalf("count mismatch: got=%d want=%d", b.Count(), r.count())
	}
}

func TestAny(t *testing.T) {
	t.Parallel()

	// Empty bitmap
	b := btmp.New(0)
	if b.Any() {
		t.Fatal("empty bitmap should return false for Any()")
	}

	// Zero-length bitmap
	b = btmp.New(64)
	if b.Any() {
		t.Fatal("zero bitmap should return false for Any()")
	}

	// Set single bit
	b = b.SetRange(32, 1)
	if !b.Any() {
		t.Fatal("bitmap with set bit should return true for Any()")
	}

	// Clear all bits
	b = b.ClearRange(32, 1)
	if b.Any() {
		t.Fatal("cleared bitmap should return false for Any()")
	}

	// Multiple words, set bit in middle word
	b = btmp.New(192)
	b = b.SetRange(128, 1)
	if !b.Any() {
		t.Fatal("bitmap with middle word set should return true for Any()")
	}

	// Set bit in last word
	b = btmp.New(192)
	b = b.SetRange(191, 1)
	if !b.Any() {
		t.Fatal("bitmap with last bit set should return true for Any()")
	}
}

func TestAny_EarlyReturn(t *testing.T) {
	t.Parallel()

	// Create bitmap with bit set in first word (covers line 73)
	b := btmp.New(128)
	b = b.SetRange(0, 1) // Set first bit

	// This should return true immediately without checking other words
	if !b.Any() {
		t.Fatal("Any() should return true for bit in first word")
	}
}

func TestCount_Empty(t *testing.T) {
	t.Parallel()

	// Empty bitmap (covers line 83)
	b := btmp.New(0)
	if b.Count() != 0 {
		t.Fatalf("Count() on empty bitmap should return 0, got %d", b.Count())
	}
}

func TestNextSetBit_EdgeCases(t *testing.T) {
	t.Parallel()

	// Empty bitmap
	b := btmp.New(0)
	if got := b.NextSetBit(0); got != -1 {
		t.Fatalf("NextSetBit on empty should return -1, got %d", got)
	}

	// No set bits
	b = btmp.New(128)
	if got := b.NextSetBit(0); got != -1 {
		t.Fatalf("NextSetBit on zero bitmap should return -1, got %d", got)
	}

	// From equals Len()
	b = b.SetRange(64, 1)
	if got := b.NextSetBit(b.Len()); got != -1 {
		t.Fatalf("NextSetBit(Len()) should return -1, got %d", got)
	}

	// Multiple words, bit in last word
	b = btmp.New(192)
	b = b.SetRange(191, 1)
	if got := b.NextSetBit(0); got != 191 {
		t.Fatalf("NextSetBit should find bit in last word, got %d want 191", got)
	}

	// Search from middle of word
	b = btmp.New(128)
	b = b.SetRange(100, 1)
	if got := b.NextSetBit(90); got != 100 {
		t.Fatalf("NextSetBit from middle should find bit, got %d want 100", got)
	}
}

func TestGrowthMethods(t *testing.T) {
	t.Parallel()

	// ReserveCap capacity growth
	b := btmp.New(0)
	initialCap := cap(b.Words())
	b = b.ReserveCap(1024) // Force capacity growth
	newCap := cap(b.Words())
	if newCap <= initialCap {
		t.Fatal("ReserveCap should increase capacity")
	}
	if b.Len() != 0 {
		t.Fatalf("ReserveCap should not change length, got %d", b.Len())
	}

	// ReserveCap no-op when sufficient capacity
	b = btmp.New(64)
	beforeCap := cap(b.Words())
	b = b.ReserveCap(32) // Less than current need
	afterCap := cap(b.Words())
	if afterCap != beforeCap {
		t.Fatal("ReserveCap should be no-op when capacity sufficient")
	}

	// Truncate - ensure it sets words slice to minimum needed length
	b = btmp.New(1024) // 1024 bits = 16 words
	b = b.Truncate()   // Should keep 16 words since we need them all
	if len(b.Words()) != 16 {
		t.Fatalf("Truncate should keep necessary words, got %d want 16", len(b.Words()))
	}

	// Test with a case where we might have excess (though EnsureBits doesn't create excess)
	b = btmp.New(127) // 127 bits = 2 words needed
	b = b.Truncate()
	expectedLen := (127 + 63) / 64 // ceil(127/64) = 2
	if len(b.Words()) != expectedLen {
		t.Fatalf("Truncate should have %d words for 127 bits, got %d", expectedLen, len(b.Words()))
	}

	// Clip
	b = btmp.New(64)
	b = b.ReserveCap(1024) // Create excess capacity
	oldCap := cap(b.Words())
	b = b.Clip()
	newCap = cap(b.Words())
	if newCap >= oldCap {
		t.Fatal("Clip should reduce capacity")
	}
	// Ensure data integrity after clip
	b = b.SetRange(32, 16)
	for i := range 16 {
		if !b.Test(32 + i) {
			t.Fatalf("bit %d should be set after clip", 32+i)
		}
	}
}

func TestTruncate_ActualTruncation(t *testing.T) {
	t.Parallel()

	// Note: The condition need < len(b.words) in Truncate (line 51)
	// may not be reachable through the public API since all mutators
	// call finalize() which keeps the words slice properly sized.
	// This test verifies Truncate works correctly when called.

	b := btmp.New(256) // 4 words
	initialLen := len(b.Words())
	initialCap := cap(b.Words())

	// Truncate should be a no-op when already properly sized
	b = b.Truncate()

	if len(b.Words()) != initialLen {
		t.Fatalf("Truncate changed length when already correct: before=%d after=%d",
			initialLen, len(b.Words()))
	}

	if cap(b.Words()) != initialCap {
		t.Fatalf("Truncate should never change capacity: before=%d after=%d",
			initialCap, cap(b.Words()))
	}
}

func TestExtractWord_Boundary(t *testing.T) {
	t.Parallel()

	// Test extractWord boundary case via CopyRange
	b1 := btmp.New(128)
	b1 = b1.SetRange(120, 8) // Set bits near end

	b2 := btmp.New(0)
	// This should trigger extractWord boundary case where sw+1 >= len(src.words)
	b2 = b2.CopyRange(b1, 120, 0, 8)

	// Verify copy worked correctly
	for i := range 8 {
		if !b2.Test(i) {
			t.Fatalf("bit %d should be set after boundary copy", i)
		}
	}

	// Test boundary with exact word boundary
	b3 := btmp.New(64)
	b3 = b3.SetRange(60, 4)
	b4 := btmp.New(0)
	b4 = b4.CopyRange(b3, 60, 0, 4)

	for i := range 4 {
		if !b4.Test(i) {
			t.Fatalf("bit %d should be set after word boundary copy", i)
		}
	}
}

func TestExtractWord_CompletelyOutOfBounds(t *testing.T) {
	t.Parallel()

	// This test covers line 160 where extractWord returns 0 for out of bounds

	// Create small bitmap
	b1 := btmp.New(64)
	b1 = b1.SetRange(0, 64)

	b2 := btmp.New(0)
	// Copy from position that would make sw >= len(src.words)
	// Using count=0 to avoid panic, but still test the path
	b2 = b2.CopyRange(b1, 64, 0, 0)

	// Also test with actual copy from near boundary
	b3 := btmp.New(128)
	b3 = b3.SetRange(0, 64) // Only first 64 bits set
	b4 := btmp.New(0)
	// Copy from position 120 (which is in second word that's all zeros)
	b4 = b4.CopyRange(b3, 120, 0, 8)

	// Verify no bits set (since source was all zeros in that range)
	if b4.Count() != 0 {
		t.Fatal("Should have no bits set when copying zeros from second word")
	}

	// Test copying from exactly at the boundary
	b5 := btmp.New(64)
	b5 = b5.SetRange(63, 1) // Set last bit of first word
	b6 := btmp.New(0)
	// This tests the boundary condition in extractWord
	b6 = b6.CopyRange(b5, 63, 0, 1)

	if !b6.Test(0) {
		t.Fatal("Should have copied the bit at boundary")
	}
}

func TestWordsExposure_TailMask(t *testing.T) {
	t.Parallel()
	b := btmp.New(130) // needWords = 3
	w := b.Words()
	if len(w) != 3 {
		t.Fatalf("unexpected words len: got=%d want=3", len(w))
	}

	// Dirty all words including tail, then trigger a no-op mutator to enforce masking.
	w[0] = ^uint64(0)
	w[1] = ^uint64(0)
	w[2] = ^uint64(0)

	b = b.SetRange(0, 0) // enforce tail mask path

	lastIdx := (b.Len()+63)>>6 - 1
	tailBits := b.Len() & 63
	var mask uint64
	if tailBits == 0 {
		mask = ^uint64(0)
	} else {
		mask = (uint64(1) << uint(tailBits)) - 1
	}
	if (b.Words()[lastIdx] &^ mask) != 0 {
		t.Fatalf("tail not masked after mutator")
	}
}

func TestPanicCases(t *testing.T) {
	t.Parallel()
	b := btmp.New(100)

	// EnsureBits negative length
	mustPanic(t, func() { _ = b.EnsureBits(-1) })

	// ReserveCap negative
	mustPanic(t, func() { _ = b.ReserveCap(-1) })

	// SetRange negative start
	mustPanic(t, func() { _ = b.SetRange(-1, 1) })
	// SetRange negative count
	mustPanic(t, func() { _ = b.SetRange(1, -1) })

	// ClearRange negative start (implicit in checkedEnd)
	mustPanic(t, func() { _ = b.ClearRange(-1, 1) })
	// ClearRange negative count
	mustPanic(t, func() { _ = b.ClearRange(1, -1) })
	// ClearRange out of bounds
	mustPanic(t, func() { _ = b.ClearRange(90, 20) })

	// CopyRange negative start (implicit in checkedEnd)
	mustPanic(t, func() { _ = b.CopyRange(b, -1, 0, 1) })
	// CopyRange negative count (this covers line 135 in checkedEnd)
	mustPanic(t, func() { _ = b.CopyRange(b, 0, 0, -1) })
	// CopyRange nil src
	mustPanic(t, func() { _ = b.CopyRange(nil, 0, 0, 1) })
	// CopyRange source out of bounds
	mustPanic(t, func() { _ = b.CopyRange(b, 90, 0, 20) })

	// Overflow cases
	mustPanic(t, func() { _ = b.SetRange(math.MaxInt-10, 20) })
	mustPanic(t, func() { _ = b.ClearRange(math.MaxInt-10, 20) })
	mustPanic(t, func() { _ = b.CopyRange(b, math.MaxInt-10, 0, 20) })
	mustPanic(t, func() { _ = b.CopyRange(b, 0, math.MaxInt-10, 20) })

	// Test out of bounds
	mustPanic(t, func() { _ = b.Test(-1) })
	mustPanic(t, func() { _ = b.Test(100) })

	// NextSetBit out of bounds
	mustPanic(t, func() { _ = b.NextSetBit(-1) })
	mustPanic(t, func() { _ = b.NextSetBit(101) })
}
