package btmp_test

import (
	"testing"

	"github.com/neox5/btmp"
)

// TestBitmapCountZerosFromInRange validates Bitmap.CountZerosFromInRange() query operation behavior.
// This test suite specifically targets the bug in countBitsFromInRange where TrailingZeros64
// (absolute position) is compared to OnesCount64 (bit count) on line 1743.
func TestBitmapCountZerosFromInRange(t *testing.T) {
	t.Run("returns 0 when starting bit is set", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(50)

		count := b.CountZerosFromInRange(50, 10)
		if count != 0 {
			t.Errorf("expected count=0, got %d", count)
		}
	})

	t.Run("counts all zeros when range is completely clear", func(t *testing.T) {
		b := btmp.New(100)

		count := b.CountZerosFromInRange(50, 10)
		if count != 10 {
			t.Errorf("expected count=10, got %d", count)
		}
	})

	t.Run("stops at first set bit in single word range", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(55) // Set bit in middle of range [50, 60)

		count := b.CountZerosFromInRange(50, 10)
		if count != 5 {
			t.Errorf("expected count=5 (bits 50-54 clear), got %d", count)
		}
	})

	t.Run("stops at first set bit immediately after start", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(51) // Second bit in range

		count := b.CountZerosFromInRange(50, 10)
		if count != 1 {
			t.Errorf("expected count=1 (only bit 50 clear), got %d", count)
		}
	})

	t.Run("stops at first set bit near end of range", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(58) // Near end of range [50, 60)

		count := b.CountZerosFromInRange(50, 10)
		if count != 8 {
			t.Errorf("expected count=8 (bits 50-57 clear), got %d", count)
		}
	})

	t.Run("stops at last bit in range", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(59) // Last bit in range [50, 60)

		count := b.CountZerosFromInRange(50, 10)
		if count != 9 {
			t.Errorf("expected count=9 (bits 50-58 clear), got %d", count)
		}
	})

	t.Run("stops with multiple consecutive set bits", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(55)
		b.SetBit(56)
		b.SetBit(57)

		count := b.CountZerosFromInRange(50, 10)
		if count != 5 {
			t.Errorf("expected count=5 (bits 50-54 clear), got %d", count)
		}
	})

	t.Run("handles range starting at bit 0", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(5)

		count := b.CountZerosFromInRange(0, 10)
		if count != 5 {
			t.Errorf("expected count=5 (bits 0-4 clear), got %d", count)
		}
	})

	t.Run("handles range at word boundary start", func(t *testing.T) {
		b := btmp.New(200)
		b.SetBit(68) // In range [64, 74)

		count := b.CountZerosFromInRange(64, 10)
		if count != 4 {
			t.Errorf("expected count=4 (bits 64-67 clear), got %d", count)
		}
	})

	t.Run("handles range crossing word boundary", func(t *testing.T) {
		b := btmp.New(200)
		b.SetBit(65) // In second word, range [60, 70) crosses boundary at 64

		count := b.CountZerosFromInRange(60, 10)
		if count != 5 {
			t.Errorf("expected count=5 (bits 60-64 clear in first word), got %d", count)
		}
	})

	t.Run("counts full word when all clear", func(t *testing.T) {
		b := btmp.New(200)

		count := b.CountZerosFromInRange(0, 64)
		if count != 64 {
			t.Errorf("expected count=64, got %d", count)
		}
	})

	t.Run("counts to end of range when no set bits", func(t *testing.T) {
		b := btmp.New(200)

		count := b.CountZerosFromInRange(50, 20)
		if count != 20 {
			t.Errorf("expected count=20, got %d", count)
		}
	})

	t.Run("handles small count value", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(52)

		count := b.CountZerosFromInRange(50, 5)
		if count != 2 {
			t.Errorf("expected count=2 (bits 50-51 clear), got %d", count)
		}
	})

	t.Run("handles count of 1 with clear bit", func(t *testing.T) {
		b := btmp.New(100)

		count := b.CountZerosFromInRange(50, 1)
		if count != 1 {
			t.Errorf("expected count=1, got %d", count)
		}
	})

	t.Run("handles count of 1 with set bit", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(50)

		count := b.CountZerosFromInRange(50, 1)
		if count != 0 {
			t.Errorf("expected count=0, got %d", count)
		}
	})

	// BUG EXPOSURE: Grid scenario that fails
	t.Run("BUG: grid scenario with 10-column row", func(t *testing.T) {
		b := btmp.New(100)
		// Simulate 10x10 grid, row 5 starts at bit 50
		// Columns 5-7 occupied (bits 55-57)
		b.SetBit(55)
		b.SetBit(56)
		b.SetBit(57)

		count := b.CountZerosFromInRange(50, 10)

		// Should count bits 50-54 (5 zeros), then stop at bit 55
		if count != 5 {
			t.Errorf("BUG REPRODUCED: expected count=5, got %d", count)
			t.Logf("This is the exact scenario from TestGridFreeColsFrom failure")
			t.Logf("Range [50, 60) in single word, first set bit at 55")
			t.Logf("TrailingZeros64(inverted)=55, bitsInMask=10, comparison 55<10 fails")
		}
	})

	t.Run("BUG: grid scenario with earlier occupied column", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(53) // Column 3 occupied

		count := b.CountZerosFromInRange(50, 10)

		if count != 3 {
			t.Errorf("BUG REPRODUCED: expected count=3, got %d", count)
		}
	})

	t.Run("panics on negative pos", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative pos")
			}
		}()
		b := btmp.New(100)
		b.CountZerosFromInRange(-1, 10)
	})

	t.Run("panics on negative count", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative count")
			}
		}()
		b := btmp.New(100)
		b.CountZerosFromInRange(50, -1)
	})

	t.Run("panics when range exceeds bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for range exceeding bounds")
			}
		}()
		b := btmp.New(100)
		b.CountZerosFromInRange(95, 10)
	})
}

// TestBitmapCountOnesFromInRange validates Bitmap.CountOnesFromInRange() query operation behavior.
func TestBitmapCountOnesFromInRange(t *testing.T) {
	t.Run("returns 0 when starting bit is clear", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(51, 9) // Set bits 51-59, but not 50

		count := b.CountOnesFromInRange(50, 10)
		if count != 0 {
			t.Errorf("expected count=0, got %d", count)
		}
	})

	t.Run("counts all ones when range is completely set", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(50, 10)

		count := b.CountOnesFromInRange(50, 10)
		if count != 10 {
			t.Errorf("expected count=10, got %d", count)
		}
	})

	t.Run("stops at first clear bit in single word range", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(50, 10)
		b.ClearBit(55) // Clear bit in middle

		count := b.CountOnesFromInRange(50, 10)
		if count != 5 {
			t.Errorf("expected count=5 (bits 50-54 set), got %d", count)
		}
	})

	t.Run("stops at first clear bit immediately after start", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(50, 10)
		b.ClearBit(51)

		count := b.CountOnesFromInRange(50, 10)
		if count != 1 {
			t.Errorf("expected count=1 (only bit 50 set), got %d", count)
		}
	})

	t.Run("stops at first clear bit near end of range", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(50, 10)
		b.ClearBit(58)

		count := b.CountOnesFromInRange(50, 10)
		if count != 8 {
			t.Errorf("expected count=8 (bits 50-57 set), got %d", count)
		}
	})

	t.Run("stops at last bit in range", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(50, 10)
		b.ClearBit(59)

		count := b.CountOnesFromInRange(50, 10)
		if count != 9 {
			t.Errorf("expected count=9 (bits 50-58 set), got %d", count)
		}
	})

	t.Run("stops with multiple consecutive clear bits", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(50, 10)
		b.ClearBit(55)
		b.ClearBit(56)
		b.ClearBit(57)

		count := b.CountOnesFromInRange(50, 10)
		if count != 5 {
			t.Errorf("expected count=5 (bits 50-54 set), got %d", count)
		}
	})

	t.Run("handles range starting at bit 0", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(0, 10)
		b.ClearBit(5)

		count := b.CountOnesFromInRange(0, 10)
		if count != 5 {
			t.Errorf("expected count=5 (bits 0-4 set), got %d", count)
		}
	})

	t.Run("handles range at word boundary", func(t *testing.T) {
		b := btmp.New(200)
		b.SetRange(64, 10)
		b.ClearBit(68)

		count := b.CountOnesFromInRange(64, 10)
		if count != 4 {
			t.Errorf("expected count=4 (bits 64-67 set), got %d", count)
		}
	})

	t.Run("handles range crossing word boundary", func(t *testing.T) {
		b := btmp.New(200)
		b.SetRange(60, 10)
		b.ClearBit(65)

		count := b.CountOnesFromInRange(60, 10)
		if count != 5 {
			t.Errorf("expected count=5 (bits 60-64 set), got %d", count)
		}
	})

	t.Run("counts full word when all set", func(t *testing.T) {
		b := btmp.New(200)
		b.SetRange(0, 64)

		count := b.CountOnesFromInRange(0, 64)
		if count != 64 {
			t.Errorf("expected count=64, got %d", count)
		}
	})

	t.Run("counts to end of range when no clear bits", func(t *testing.T) {
		b := btmp.New(200)
		b.SetRange(50, 20)

		count := b.CountOnesFromInRange(50, 20)
		if count != 20 {
			t.Errorf("expected count=20, got %d", count)
		}
	})

	t.Run("handles count of 1 with set bit", func(t *testing.T) {
		b := btmp.New(100)
		b.SetBit(50)

		count := b.CountOnesFromInRange(50, 1)
		if count != 1 {
			t.Errorf("expected count=1, got %d", count)
		}
	})

	t.Run("handles count of 1 with clear bit", func(t *testing.T) {
		b := btmp.New(100)

		count := b.CountOnesFromInRange(50, 1)
		if count != 0 {
			t.Errorf("expected count=0, got %d", count)
		}
	})

	// BUG EXPOSURE: Symmetric scenario for ones
	t.Run("BUG: counts ones with gap in single word range", func(t *testing.T) {
		b := btmp.New(100)
		b.SetRange(50, 10) // Set all bits [50, 60)
		b.ClearBit(55)     // Clear bit in middle
		b.ClearBit(56)
		b.ClearBit(57)

		count := b.CountOnesFromInRange(50, 10)

		if count != 5 {
			t.Errorf("BUG REPRODUCED: expected count=5, got %d", count)
			t.Logf("Symmetric scenario to CountZerosFromInRange bug")
		}
	})

	t.Run("panics on negative pos", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative pos")
			}
		}()
		b := btmp.New(100)
		b.CountOnesFromInRange(-1, 10)
	})

	t.Run("panics on negative count", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative count")
			}
		}()
		b := btmp.New(100)
		b.CountOnesFromInRange(50, -1)
	})

	t.Run("panics when range exceeds bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for range exceeding bounds")
			}
		}()
		b := btmp.New(100)
		b.CountOnesFromInRange(95, 10)
	})
}
