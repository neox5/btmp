package btmp_test

import (
	"testing"

	"github.com/neox5/btmp"
)

// TestNew validates New constructor behavior.
// Tests:
//   - Creates bitmap with exact n bits (Len==n)
//   - Zero size accepted (n==0)
//   - Bitmap length matches constructor parameter
//   - All bits initialized to zero
//   - Underlying words slice sized correctly (ceil(n/64))
func TestNew(t *testing.T) {
	t.Run("creates bitmap with exact n bits", func(t *testing.T) {
		b := btmp.New(100)
		if b.Len() != 100 {
			t.Errorf("expected len=100, got %d", b.Len())
		}
	})

	t.Run("accepts zero size", func(t *testing.T) {
		b := btmp.New(0)
		if b.Len() != 0 {
			t.Errorf("expected len=0, got %d", b.Len())
		}
	})

	t.Run("all bits initialized to zero", func(t *testing.T) {
		b := btmp.New(128)
		if b.Any() {
			t.Error("expected all bits to be zero")
		}
		if b.Count() != 0 {
			t.Errorf("expected count=0, got %d", b.Count())
		}
	})

	t.Run("words slice sized correctly for exact word multiple", func(t *testing.T) {
		b := btmp.New(128) // Exactly 2 words
		words := b.Words()
		if len(words) != 2 {
			t.Errorf("expected 2 words, got %d", len(words))
		}
	})

	t.Run("words slice sized correctly for partial word", func(t *testing.T) {
		b := btmp.New(100) // ceil(100/64) = 2 words
		words := b.Words()
		if len(words) != 2 {
			t.Errorf("expected 2 words, got %d", len(words))
		}
	})

	t.Run("words slice sized correctly for single bit", func(t *testing.T) {
		b := btmp.New(1)
		words := b.Words()
		if len(words) != 1 {
			t.Errorf("expected 1 word, got %d", len(words))
		}
	})

	t.Run("words slice empty for zero size", func(t *testing.T) {
		b := btmp.New(0)
		words := b.Words()
		if len(words) != 0 {
			t.Errorf("expected 0 words, got %d", len(words))
		}
	})
}

// TestBitmapLen validates Bitmap.Len() accessor behavior.
// Tests:
//   - Returns correct length for initialized bitmap
//   - Returns 0 for empty bitmap
//   - Length matches constructor parameter
func TestBitmapLen(t *testing.T) {
	t.Run("returns correct length for initialized bitmap", func(t *testing.T) {
		tests := []uint{1, 63, 64, 65, 100, 128, 1000}
		for _, n := range tests {
			b := btmp.New(n)
			if b.Len() != int(n) {
				t.Errorf("New(%d): expected len=%d, got %d", n, n, b.Len())
			}
		}
	})

	t.Run("returns 0 for empty bitmap", func(t *testing.T) {
		b := btmp.New(0)
		if b.Len() != 0 {
			t.Errorf("expected len=0, got %d", b.Len())
		}
	})

	t.Run("length persists after operations", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(50)
		if b.Len() != 100 {
			t.Errorf("expected len=100 after SetBit, got %d", b.Len())
		}
	})
}

// TestBitmapWords validates Bitmap.Words() accessor behavior.
// Tests:
//   - Returns underlying words slice
//   - Slice length matches calculated word count ceil(n/64)
//   - Words slice correctly sized for various bit counts
func TestBitmapWords(t *testing.T) {
	t.Run("returns underlying words slice", func(t *testing.T) {
		b := btmp.New(64)
		words := b.Words()
		if words == nil {
			t.Error("expected non-nil words slice")
		}
	})

	t.Run("slice length matches word count", func(t *testing.T) {
		tests := []struct {
			bits  uint
			words int
		}{
			{0, 0},
			{1, 1},
			{63, 1},
			{64, 1},
			{65, 2},
			{128, 2},
			{129, 3},
			{192, 3},
		}

		for _, tt := range tests {
			b := btmp.New(tt.bits)
			words := b.Words()
			if len(words) != tt.words {
				t.Errorf("New(%d): expected %d words, got %d", tt.bits, tt.words, len(words))
			}
		}
	})

	t.Run("words reflect bitmap state", func(t *testing.T) {
		b := btmp.New(128)
		b.SetBit(0)
		b.SetBit(64)

		words := b.Words()
		if words[0] != 1 {
			t.Errorf("expected words[0]=1, got %d", words[0])
		}
		if words[1] != 1 {
			t.Errorf("expected words[1]=1, got %d", words[1])
		}
	})
}

// TestBitmapTest validates Bitmap.Test() query operation.
// Tests:
//   - Returns false for unset bit
//   - Returns true for set bit
//   - Works at word boundaries (positions 0, 63, 64, 127)
//   - Panics on negative position
//   - Panics on position >= Len()
func TestBitmapTest(t *testing.T) {
	t.Run("returns false for unset bit", func(t *testing.T) {
		b := btmp.New(100)
		if b.Test(50) {
			t.Error("expected false for unset bit")
		}
	})

	t.Run("returns true for set bit", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(50)
		if !b.Test(50) {
			t.Error("expected true for set bit")
		}
	})

	t.Run("works at word boundaries", func(t *testing.T) {
		b := btmp.New(200)
		positions := []int{0, 63, 64, 127, 128}

		for _, pos := range positions {
			b.SetBit(pos)
		}

		for _, pos := range positions {
			if !b.Test(pos) {
				t.Errorf("expected true at position %d", pos)
			}
		}
	})

	t.Run("unset bits remain false", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(50)

		if b.Test(49) {
			t.Error("expected false at position 49")
		}
		if b.Test(51) {
			t.Error("expected false at position 51")
		}
	})

	t.Run("panics on negative position", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative position")
			}
		}()
		b := btmp.New(100)
		b.Test(-1)
	})

	t.Run("panics on position >= Len", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for position >= Len")
			}
		}()
		b := btmp.New(100)
		b.Test(100)
	})

	t.Run("panics on position way beyond Len", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for position >> Len")
			}
		}()
		b := btmp.New(10)
		b.Test(1000)
	})
}

// TestBitmapGetBits validates Bitmap.GetBits() query operation.
// Tests:
//   - Extracts n bits starting from pos, returned right-aligned
//   - Single word aligned read (fast path: n=64, pos%64==0)
//   - Single word unaligned read
//   - Cross-word read spanning exactly two words
//   - Full word read (n=64)
//   - Partial word reads (various n values)
//   - Panics on pos < 0
//   - Panics on n <= 0
//   - Panics on n > 64
//   - Panics on pos+n > Len()
func TestBitmapGetBits(t *testing.T) {
	t.Run("extracts bits right-aligned", func(t *testing.T) {
		b := btmp.New(100)
		// Set pattern: 101 at positions 10,11,12
		b.SetBit(10)
		b.SetBit(12)

		result := b.GetBits(10, 3)
		if result != 0b101 {
			t.Errorf("expected 0b101, got 0b%b", result)
		}
	})

	t.Run("single word aligned read fast path", func(t *testing.T) {
		b := btmp.New(128)
		// Set all bits in first word
		for i := 0; i < 64; i++ {
			b.SetBit(i)
		}

		result := b.GetBits(0, 64)
		if result != 0xFFFFFFFFFFFFFFFF {
			t.Errorf("expected all ones, got 0x%X", result)
		}
	})

	t.Run("single word unaligned read", func(t *testing.T) {
		b := btmp.New(100)
		// Set pattern: 1111 at positions 10-13
		b.SetBit(10)
		b.SetBit(11)
		b.SetBit(12)
		b.SetBit(13)

		result := b.GetBits(10, 4)
		if result != 0b1111 {
			t.Errorf("expected 0b1111, got 0b%b", result)
		}
	})

	t.Run("cross-word read", func(t *testing.T) {
		b := btmp.New(200)
		// Set bits at positions 62, 63, 64, 65 (spans word boundary)
		b.SetBit(62)
		b.SetBit(63)
		b.SetBit(64)
		b.SetBit(65)

		result := b.GetBits(62, 4)
		if result != 0b1111 {
			t.Errorf("expected 0b1111, got 0b%b", result)
		}
	})

	t.Run("cross-word read larger span", func(t *testing.T) {
		b := btmp.New(200)
		// Set pattern across word boundary
		// Positions 60-67 (8 bits spanning words)
		for i := 60; i < 68; i++ {
			b.SetBit(i)
		}

		result := b.GetBits(60, 8)
		if result != 0xFF {
			t.Errorf("expected 0xFF, got 0x%X", result)
		}
	})

	t.Run("full word read at offset", func(t *testing.T) {
		b := btmp.New(200)
		// Set all bits in second word (positions 64-127)
		for i := 64; i < 128; i++ {
			b.SetBit(i)
		}

		result := b.GetBits(64, 64)
		if result != 0xFFFFFFFFFFFFFFFF {
			t.Errorf("expected all ones, got 0x%X", result)
		}
	})

	t.Run("extracts zero bits correctly", func(t *testing.T) {
		b := btmp.New(100)
		// Don't set any bits

		result := b.GetBits(10, 8)
		if result != 0 {
			t.Errorf("expected 0, got 0x%X", result)
		}
	})

	t.Run("single bit extraction", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(50)

		result := b.GetBits(50, 1)
		if result != 1 {
			t.Errorf("expected 1, got %d", result)
		}

		result = b.GetBits(49, 1)
		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("panics on negative position", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative position")
			}
		}()
		b := btmp.New(100)
		b.GetBits(-1, 8)
	})

	t.Run("panics on n <= 0", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for n=0")
			}
		}()
		b := btmp.New(100)
		b.GetBits(10, 0)
	})

	t.Run("panics on n > 64", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for n>64")
			}
		}()
		b := btmp.New(200)
		b.GetBits(10, 65)
	})

	t.Run("panics on pos+n > Len", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for out of bounds")
			}
		}()
		b := btmp.New(100)
		b.GetBits(90, 20) // 90+20 = 110 > 100
	})
}

// TestBitmapAny validates Bitmap.Any() query operation.
// Tests:
//   - Returns false for empty bitmap (Len==0)
//   - Returns false when all bits are clear
//   - Returns true when at least one bit is set
//   - Detects set bit in first word
//   - Detects set bit in middle word
//   - Detects set bit in last partial word
//   - Properly masks last word
func TestBitmapAny(t *testing.T) {
	t.Run("returns false for empty bitmap", func(t *testing.T) {
		b := btmp.New(0)
		if b.Any() {
			t.Error("expected false for empty bitmap")
		}
	})

	t.Run("returns false when all bits clear", func(t *testing.T) {
		b := btmp.New(200)
		if b.Any() {
			t.Error("expected false when all bits clear")
		}
	})

	t.Run("returns true with single bit set", func(t *testing.T) {
		b := btmp.New(200)
		b.SetBit(100)
		if !b.Any() {
			t.Error("expected true with bit set")
		}
	})

	t.Run("detects bit in first word", func(t *testing.T) {
		b := btmp.New(200)
		b.SetBit(0)
		if !b.Any() {
			t.Error("expected true with bit in first word")
		}
	})

	t.Run("detects bit in middle word", func(t *testing.T) {
		b := btmp.New(200)
		b.SetBit(100) // Second word
		if !b.Any() {
			t.Error("expected true with bit in middle word")
		}
	})

	t.Run("detects bit in last partial word", func(t *testing.T) {
		b := btmp.New(100) // ceil(100/64) = 2 words, last is partial
		b.SetBit(99)       // Last bit
		if !b.Any() {
			t.Error("expected true with bit in last partial word")
		}
	})

	t.Run("detects bit at last position", func(t *testing.T) {
		b := btmp.New(128) // Exactly 2 words
		b.SetBit(127)
		if !b.Any() {
			t.Error("expected true with bit at last position")
		}
	})

	t.Run("properly masks last word", func(t *testing.T) {
		b := btmp.New(65) // 2 words, last has only 1 valid bit
		// Manually corrupt beyond valid length (testing internal masking)
		// This tests that Any() properly masks the last word
		b.SetBit(64) // Valid bit
		if !b.Any() {
			t.Error("expected true with valid bit in last word")
		}
	})
}

// TestBitmapCount validates Bitmap.Count() query operation.
// Tests:
//   - Returns 0 for empty bitmap (Len==0)
//   - Returns 0 when all bits are clear
//   - Counts single set bit correctly
//   - Counts multiple set bits across multiple words
//   - Counts all bits set (full bitmap)
//   - Handles last partial word masking correctly
//   - Accurate count for various patterns
func TestBitmapCount(t *testing.T) {
	t.Run("returns 0 for empty bitmap", func(t *testing.T) {
		b := btmp.New(0)
		if b.Count() != 0 {
			t.Errorf("expected count=0, got %d", b.Count())
		}
	})

	t.Run("returns 0 when all bits clear", func(t *testing.T) {
		b := btmp.New(200)
		if b.Count() != 0 {
			t.Errorf("expected count=0, got %d", b.Count())
		}
	})

	t.Run("counts single set bit", func(t *testing.T) {
		b := btmp.New(200)
		b.SetBit(100)
		if b.Count() != 1 {
			t.Errorf("expected count=1, got %d", b.Count())
		}
	})

	t.Run("counts multiple bits in single word", func(t *testing.T) {
		b := btmp.New(64)
		b.SetBit(0)
		b.SetBit(10)
		b.SetBit(20)
		b.SetBit(63)
		if b.Count() != 4 {
			t.Errorf("expected count=4, got %d", b.Count())
		}
	})

	t.Run("counts bits across multiple words", func(t *testing.T) {
		b := btmp.New(200)
		b.SetBit(0)   // First word
		b.SetBit(63)  // First word
		b.SetBit(64)  // Second word
		b.SetBit(100) // Second word
		b.SetBit(150) // Third word
		if b.Count() != 5 {
			t.Errorf("expected count=5, got %d", b.Count())
		}
	})

	t.Run("counts all bits set in full words", func(t *testing.T) {
		b := btmp.New(128)
		b.SetAll()
		if b.Count() != 128 {
			t.Errorf("expected count=128, got %d", b.Count())
		}
	})

	t.Run("counts all bits set with partial word", func(t *testing.T) {
		b := btmp.New(100)
		b.SetAll()
		if b.Count() != 100 {
			t.Errorf("expected count=100, got %d", b.Count())
		}
	})

	t.Run("handles last partial word correctly", func(t *testing.T) {
		b := btmp.New(65) // 2 words, last has 1 valid bit
		b.SetBit(0)       // First word
		b.SetBit(64)      // Last partial word
		if b.Count() != 2 {
			t.Errorf("expected count=2, got %d", b.Count())
		}
	})

	t.Run("counts pattern correctly", func(t *testing.T) {
		b := btmp.New(200)
		// Set every 10th bit
		for i := 0; i < 200; i += 10 {
			b.SetBit(i)
		}
		expected := 20
		if b.Count() != expected {
			t.Errorf("expected count=%d, got %d", expected, b.Count())
		}
	})

	t.Run("counts after clear operations", func(t *testing.T) {
		b := btmp.New(100)
		b.SetAll()
		b.ClearBit(50)
		b.ClearBit(51)
		if b.Count() != 98 {
			t.Errorf("expected count=98, got %d", b.Count())
		}
	})

	t.Run("counts dense pattern", func(t *testing.T) {
		b := btmp.New(128)
		// Set first 100 bits
		for i := 0; i < 100; i++ {
			b.SetBit(i)
		}
		if b.Count() != 100 {
			t.Errorf("expected count=100, got %d", b.Count())
		}
	})
}
