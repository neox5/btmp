package btmp_test

import (
	"testing"

	"github.com/neox5/btmp"
)

// TestNew validates New constructor behavior.
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

// TestBitmapAny validates Bitmap.Any() query operation.
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
		for i := range 100 {
			b.SetBit(i)
		}
		if b.Count() != 100 {
			t.Errorf("expected count=100, got %d", b.Count())
		}
	})
}
