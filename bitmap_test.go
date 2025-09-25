package btmp_test

import (
	"math"
	"testing"

	btmp "github.com/neox5/btmp"
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

	// Clear OOB
	mustPanic(t, func() { _ = b.ClearRange(90, 20) })
	// Negative args
	mustPanic(t, func() { _ = b.SetRange(-1, 1) })
	mustPanic(t, func() { _ = b.SetRange(1, -1) })
	mustPanic(t, func() { _ = b.ClearRange(-1, 1) })
	mustPanic(t, func() { _ = b.CopyRange(b, -1, 0, 1) })
	// Overflow
	mustPanic(t, func() { _ = b.SetRange(math.MaxInt-10, 20) })
	// Test/NextSetBit OOB
	mustPanic(t, func() { _ = b.Test(-1) })
	mustPanic(t, func() { _ = b.Test(100) })
	mustPanic(t, func() { _ = b.NextSetBit(-1) })
	mustPanic(t, func() { _ = b.NextSetBit(101) })
	// Copy nil src
	mustPanic(t, func() { _ = b.CopyRange(nil, 0, 0, 1) })
}
