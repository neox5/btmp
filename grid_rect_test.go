package btmp_test

import (
	"testing"

	"github.com/neox5/btmp"
)

// TestGridIsFree validates Grid.IsFree() query operation behavior.
func TestGridIsFree(t *testing.T) {
	t.Run("returns true when all bits are zero", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		// Don't set any bits

		if !g.IsFree(2, 2, 5, 5) {
			t.Error("expected true for all zero bits")
		}
	})

	t.Run("returns false when single bit is set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(5, 5)) // Set one bit in middle

		if g.IsFree(3, 3, 5, 5) {
			t.Error("expected false when bit is set in rectangle")
		}
	})

	t.Run("returns false when first bit is set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(2, 2)) // Set first bit of query rect

		if g.IsFree(2, 2, 5, 5) {
			t.Error("expected false when first bit is set")
		}
	})

	t.Run("returns false when last bit is set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(6, 6)) // Set last bit of query rect (2,2,5,5)

		if g.IsFree(2, 2, 5, 5) {
			t.Error("expected false when last bit is set")
		}
	})

	t.Run("single cell free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.IsFree(5, 5, 1, 1) {
			t.Error("expected true for single free cell")
		}
	})

	t.Run("single cell occupied", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(5, 5))

		if g.IsFree(5, 5, 1, 1) {
			t.Error("expected false for single occupied cell")
		}
	})

	t.Run("full row free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.IsFree(0, 5, 10, 1) {
			t.Error("expected true for free row")
		}
	})

	t.Run("full row occupied", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 5, 10, 1) // Fill row 5

		if g.IsFree(0, 5, 10, 1) {
			t.Error("expected false for occupied row")
		}
	})

	t.Run("full column free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.IsFree(5, 0, 1, 10) {
			t.Error("expected true for free column")
		}
	})

	t.Run("full column occupied", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 0, 1, 10) // Fill column 5

		if g.IsFree(5, 0, 1, 10) {
			t.Error("expected false for occupied column")
		}
	})

	t.Run("multi-row rectangle free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.IsFree(2, 3, 5, 4) {
			t.Error("expected true for free rectangle")
		}
	})

	t.Run("multi-row rectangle with bit in middle", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(4, 5)) // Middle of (2,3,5,4)

		if g.IsFree(2, 3, 5, 4) {
			t.Error("expected false with bit in middle")
		}
	})

	t.Run("at top-left corner", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.IsFree(0, 0, 3, 3) {
			t.Error("expected true at top-left corner")
		}
	})

	t.Run("at bottom-right corner", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.IsFree(7, 7, 3, 3) {
			t.Error("expected true at bottom-right corner")
		}
	})

	t.Run("detects bit outside on left", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(1, 5)) // Just outside (2,3,5,4)

		if !g.IsFree(2, 3, 5, 4) {
			t.Error("expected true, bit is outside rectangle")
		}
	})

	t.Run("detects bit outside on right", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(7, 5)) // Just outside (2,3,5,4) right edge

		if !g.IsFree(2, 3, 5, 4) {
			t.Error("expected true, bit is outside rectangle")
		}
	})

	t.Run("panics on negative x", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative x")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(-1, 5, 3, 3)
	})

	t.Run("panics on negative y", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative y")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, -1, 3, 3)
	})

	t.Run("panics on negative w", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative w")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 5, -1, 3)
	})

	t.Run("panics on negative h", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative h")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 5, 3, -1)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 5, 0, 3)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 5, 3, 0)
	})

	t.Run("panics when x+w exceeds cols", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for x+w > cols")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(8, 5, 3, 3) // 8+3=11 > 10
	})

	t.Run("panics when y+h exceeds rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for y+h > rows")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 8, 3, 3) // 8+3=11 > 10
	})
}

// TestGridCanShiftRight validates Grid.CanShiftRight() query operation behavior.
func TestGridCanShiftRight(t *testing.T) {
	t.Run("returns true when target column is free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 3, 3, 4) // Rectangle at (2,3,3,4)

		if !g.CanShiftRight(2, 3, 3, 4) {
			t.Error("expected true when target column (5) is free")
		}
	})

	t.Run("returns false when target column has bit", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 3, 3, 4)     // Rectangle at (2,3,3,4)
		g.B.SetBit(g.Index(5, 4)) // Set bit in target column (x+w=5)

		if g.CanShiftRight(2, 3, 3, 4) {
			t.Error("expected false when target column has set bit")
		}
	})

	t.Run("returns false when any bit in target column set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 3, 3, 4)     // Rectangle at (2,3,3,4)
		g.B.SetBit(g.Index(5, 6)) // Last row of target column

		if g.CanShiftRight(2, 3, 3, 4) {
			t.Error("expected false when any bit in target column set")
		}
	})

	t.Run("single row shift validation", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 5, 3, 1) // Single row at y=5

		if !g.CanShiftRight(2, 5, 3, 1) {
			t.Error("expected true for single row shift")
		}
	})

	t.Run("single row with occupied target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 5, 3, 1)
		g.B.SetBit(g.Index(5, 5)) // Target column occupied

		if g.CanShiftRight(2, 5, 3, 1) {
			t.Error("expected false when single row target occupied")
		}
	})

	t.Run("multi-row shift with free target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(1, 2, 4, 5) // Multi-row rectangle

		if !g.CanShiftRight(1, 2, 4, 5) {
			t.Error("expected true for multi-row shift with free target")
		}
	})

	t.Run("shift to rightmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 0, 9, 3) // Rectangle ending at col 8, can shift to 9

		if !g.CanShiftRight(0, 0, 9, 3) {
			t.Error("expected true shifting to rightmost column")
		}
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftRight(5, 5, 0, 3)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftRight(5, 5, 3, 0)
	})

	t.Run("panics on invalid source rectangle", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftRight(8, 5, 3, 3) // Source rect exceeds bounds (8+3=11)
	})
}

// TestGridCanShiftLeft validates Grid.CanShiftLeft() query operation behavior.
func TestGridCanShiftLeft(t *testing.T) {
	t.Run("returns true when target column is free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 3, 4) // Rectangle at (5,3,3,4)

		if !g.CanShiftLeft(5, 3, 3, 4) {
			t.Error("expected true when target column (4) is free")
		}
	})

	t.Run("returns false when target column has bit", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 3, 4)
		g.B.SetBit(g.Index(4, 4)) // Set bit in target column (x-1=4)

		if g.CanShiftLeft(5, 3, 3, 4) {
			t.Error("expected false when target column has set bit")
		}
	})

	t.Run("returns false when any bit in target column set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 3, 4)
		g.B.SetBit(g.Index(4, 6)) // Last row of target column

		if g.CanShiftLeft(5, 3, 3, 4) {
			t.Error("expected false when any bit in target column set")
		}
	})

	t.Run("single row shift validation", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 3, 1) // Single row at y=5

		if !g.CanShiftLeft(5, 5, 3, 1) {
			t.Error("expected true for single row shift")
		}
	})

	t.Run("single row with occupied target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 3, 1)
		g.B.SetBit(g.Index(4, 5)) // Target column occupied

		if g.CanShiftLeft(5, 5, 3, 1) {
			t.Error("expected false when single row target occupied")
		}
	})

	t.Run("multi-row shift with free target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 2, 4, 5) // Multi-row rectangle

		if !g.CanShiftLeft(5, 2, 4, 5) {
			t.Error("expected true for multi-row shift with free target")
		}
	})

	t.Run("shift to leftmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(1, 0, 5, 3) // Rectangle starting at col 1, can shift to 0

		if !g.CanShiftLeft(1, 0, 5, 3) {
			t.Error("expected true shifting to leftmost column")
		}
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftLeft(5, 5, 0, 3)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftLeft(5, 5, 3, 0)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftLeft(9, 5, 3, 3) // Source rect exceeds bounds (9+3=12)
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftLeft(5, 9, 3, 3) // Source rect exceeds bounds (9+3=12)
	})

	t.Run("panics on negative x", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative x")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftLeft(-1, 5, 3, 3)
	})
}

// TestGridCanShiftUp validates Grid.CanShiftUp() query operation behavior.
func TestGridCanShiftUp(t *testing.T) {
	t.Run("returns true when target row is free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 4, 3) // Rectangle at (3,5,4,3)

		if !g.CanShiftUp(3, 5, 4, 3) {
			t.Error("expected true when target row (4) is free")
		}
	})

	t.Run("returns false when target row has bit", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 4, 3)
		g.B.SetBit(g.Index(4, 4)) // Set bit in target row (y-1=4)

		if g.CanShiftUp(3, 5, 4, 3) {
			t.Error("expected false when target row has set bit")
		}
	})

	t.Run("returns false when any bit in target row set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 4, 3)
		g.B.SetBit(g.Index(6, 4)) // Last column of target row

		if g.CanShiftUp(3, 5, 4, 3) {
			t.Error("expected false when any bit in target row set")
		}
	})

	t.Run("single column shift validation", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 1, 3) // Single column at x=5

		if !g.CanShiftUp(5, 5, 1, 3) {
			t.Error("expected true for single column shift")
		}
	})

	t.Run("single column with occupied target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 1, 3)
		g.B.SetBit(g.Index(5, 4)) // Target row occupied

		if g.CanShiftUp(5, 5, 1, 3) {
			t.Error("expected false when single column target occupied")
		}
	})

	t.Run("multi-column shift with free target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 5, 5, 4) // Multi-column rectangle

		if !g.CanShiftUp(2, 5, 5, 4) {
			t.Error("expected true for multi-column shift with free target")
		}
	})

	t.Run("shift to topmost row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 1, 5, 3) // Rectangle starting at row 1, can shift to 0

		if !g.CanShiftUp(0, 1, 5, 3) {
			t.Error("expected true shifting to topmost row")
		}
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftUp(5, 5, 0, 3)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftUp(5, 5, 3, 0)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftUp(9, 5, 3, 3) // Source rect exceeds bounds (9+3=12)
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftUp(5, 9, 3, 3) // Source rect exceeds bounds (9+3=12)
	})

	t.Run("panics on negative y", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative y")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftUp(5, -1, 3, 3)
	})
}

// TestGridCanShiftDown validates Grid.CanShiftDown() query operation behavior.
func TestGridCanShiftDown(t *testing.T) {
	t.Run("returns true when target row is free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 4, 3) // Rectangle at (3,3,4,3), target row is 6

		if !g.CanShiftDown(3, 3, 4, 3) {
			t.Error("expected true when target row (6) is free")
		}
	})

	t.Run("returns false when target row has bit", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 4, 3)
		g.B.SetBit(g.Index(4, 6)) // Set bit in target row (y+h=6)

		if g.CanShiftDown(3, 3, 4, 3) {
			t.Error("expected false when target row has set bit")
		}
	})

	t.Run("returns false when any bit in target row set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 4, 3)
		g.B.SetBit(g.Index(6, 6)) // Last column of target row

		if g.CanShiftDown(3, 3, 4, 3) {
			t.Error("expected false when any bit in target row set")
		}
	})

	t.Run("single column shift validation", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 1, 3) // Single column at x=5

		if !g.CanShiftDown(5, 3, 1, 3) {
			t.Error("expected true for single column shift")
		}
	})

	t.Run("single column with occupied target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 1, 3)
		g.B.SetBit(g.Index(5, 6)) // Target row occupied (y+h=6)

		if g.CanShiftDown(5, 3, 1, 3) {
			t.Error("expected false when single column target occupied")
		}
	})

	t.Run("multi-column shift with free target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 2, 5, 4) // Multi-column rectangle

		if !g.CanShiftDown(2, 2, 5, 4) {
			t.Error("expected true for multi-column shift with free target")
		}
	})

	t.Run("shift to bottommost row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 0, 5, 9) // Rectangle ending at row 8, can shift to 9

		if !g.CanShiftDown(0, 0, 5, 9) {
			t.Error("expected true shifting to bottommost row")
		}
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftDown(5, 5, 0, 3)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftDown(5, 5, 3, 0)
	})

	t.Run("panics on invalid source rectangle", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftDown(8, 5, 3, 3) // Source rect exceeds bounds (8+3=11)
	})
}

// TestGridShiftRectRight validates Grid.ShiftRectRight() shift operation.
func TestGridShiftRectRight(t *testing.T) {
	t.Run("valid shift with free target column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 2, 2) // 2x2 rectangle at (3,3)

		g.ShiftRectRight(3, 3, 2, 2)

		// Verify rectangle moved to (4,3)
		if !g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3)")
		}
		if !g.B.Test(g.Index(5, 3)) {
			t.Error("expected bit at (5,3)")
		}
		if !g.B.Test(g.Index(4, 4)) {
			t.Error("expected bit at (4,4)")
		}
		if !g.B.Test(g.Index(5, 4)) {
			t.Error("expected bit at (5,4)")
		}

		// Verify leftmost column cleared (x=3)
		if g.B.Test(g.Index(3, 3)) {
			t.Error("expected bit at (3,3) to be cleared")
		}
		if g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4) to be cleared")
		}

		// Verify count unchanged
		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift clears leftmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 2, 2)

		// Verify initial state
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5) before shift")
		}
		if !g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6) before shift")
		}

		g.ShiftRectRight(5, 5, 2, 2)

		// Verify leftmost column (x=5) cleared
		if g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5) cleared after shift")
		}
		if g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6) cleared after shift")
		}
	})

	t.Run("shift preserves rectangle data", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		// Create specific pattern in 2x2 rectangle
		g.B.SetBit(g.Index(2, 2)) // top-left
		g.B.SetBit(g.Index(3, 3)) // bottom-right

		g.ShiftRectRight(2, 2, 2, 2)

		// Verify pattern preserved at new location (3,2)
		if !g.B.Test(g.Index(3, 2)) {
			t.Error("expected bit at (3,2) - top-left preserved")
		}
		if !g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3) - bottom-right preserved")
		}

		// Verify other cells in new rectangle remain clear
		if g.B.Test(g.Index(4, 2)) {
			t.Error("expected bit at (4,2) to remain clear")
		}
		if g.B.Test(g.Index(3, 3)) {
			t.Error("expected bit at (3,3) to remain clear")
		}

		if g.B.Count() != 2 {
			t.Errorf("expected count=2, got %d", g.B.Count())
		}
	})

	t.Run("multiple consecutive shifts", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(1, 5, 2, 2)

		// Shift right 3 times
		g.ShiftRectRight(1, 5, 2, 2)
		g.ShiftRectRight(2, 5, 2, 2)
		g.ShiftRectRight(3, 5, 2, 2)

		// Verify rectangle at (4,5)
		if !g.B.Test(g.Index(4, 5)) {
			t.Error("expected bit at (4,5)")
		}
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5)")
		}
		if !g.B.Test(g.Index(4, 6)) {
			t.Error("expected bit at (4,6)")
		}
		if !g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6)")
		}

		// Verify original positions cleared
		for x := range 4 {
			for y := 5; y < 7; y++ {
				if g.B.Test(g.Index(x, y)) {
					t.Errorf("expected bit at (%d,%d) to be cleared", x, y)
				}
			}
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift to rightmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(6, 3, 2, 2) // Rectangle ending at col 7, can shift to 8-9

		g.ShiftRectRight(6, 3, 2, 2)

		// Verify rectangle at (7,3)
		if !g.B.Test(g.Index(7, 3)) {
			t.Error("expected bit at (7,3)")
		}
		if !g.B.Test(g.Index(8, 3)) {
			t.Error("expected bit at (8,3)")
		}
		if !g.B.Test(g.Index(7, 4)) {
			t.Error("expected bit at (7,4)")
		}
		if !g.B.Test(g.Index(8, 4)) {
			t.Error("expected bit at (8,4)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(5, 5, 0, 2)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(5, 5, 2, 0)
	})

	t.Run("panics when target column occupied", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when target column occupied")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 2, 2)
		g.B.SetBit(g.Index(5, 3)) // Occupy target column (x+w=5)

		g.ShiftRectRight(3, 3, 2, 2)
	})

	t.Run("panics when target column out of bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when target column out of bounds")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(8, 3, 2, 2) // x+w=10, cannot shift right

		g.ShiftRectRight(8, 3, 2, 2)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(9, 3, 2, 2) // x+w=11 exceeds bounds
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(3, 9, 2, 2) // y+h=11 exceeds bounds
	})

	t.Run("panics on negative x", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative x")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(-1, 3, 2, 2)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(1, 5, 2, 2)

		result := g.ShiftRectRight(1, 5, 2, 2)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(10, 10)
		g2.SetRect(1, 5, 2, 2).
			ShiftRectRight(1, 5, 2, 2).
			ShiftRectRight(2, 5, 2, 2)

		if g2.B.Count() != 4 {
			t.Errorf("expected count=4 after chaining, got %d", g2.B.Count())
		}
	})
}

// TestGridShiftRectLeft validates Grid.ShiftRectLeft() shift operation.
func TestGridShiftRectLeft(t *testing.T) {
	t.Run("valid shift with free target column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 2, 2) // 2x2 rectangle at (5,3)

		g.ShiftRectLeft(5, 3, 2, 2)

		// Verify rectangle moved to (4,3)
		if !g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3)")
		}
		if !g.B.Test(g.Index(5, 3)) {
			t.Error("expected bit at (5,3)")
		}
		if !g.B.Test(g.Index(4, 4)) {
			t.Error("expected bit at (4,4)")
		}
		if !g.B.Test(g.Index(5, 4)) {
			t.Error("expected bit at (5,4)")
		}

		// Verify rightmost column cleared (x+w-1=6)
		if g.B.Test(g.Index(6, 3)) {
			t.Error("expected bit at (6,3) to be cleared")
		}
		if g.B.Test(g.Index(6, 4)) {
			t.Error("expected bit at (6,4) to be cleared")
		}

		// Verify count unchanged
		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift clears rightmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 2, 2)

		// Verify initial state
		if !g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5) before shift")
		}
		if !g.B.Test(g.Index(6, 6)) {
			t.Error("expected bit at (6,6) before shift")
		}

		g.ShiftRectLeft(5, 5, 2, 2)

		// Verify rightmost column (x+w-1=6) cleared
		if g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5) cleared after shift")
		}
		if g.B.Test(g.Index(6, 6)) {
			t.Error("expected bit at (6,6) cleared after shift")
		}
	})

	t.Run("shift preserves rectangle data", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		// Create specific pattern in 2x2 rectangle
		g.B.SetBit(g.Index(5, 2)) // top-left
		g.B.SetBit(g.Index(6, 3)) // bottom-right

		g.ShiftRectLeft(5, 2, 2, 2)

		// Verify pattern preserved at new location (4,2)
		if !g.B.Test(g.Index(4, 2)) {
			t.Error("expected bit at (4,2) - top-left preserved")
		}
		if !g.B.Test(g.Index(5, 3)) {
			t.Error("expected bit at (5,3) - bottom-right preserved")
		}

		// Verify other cells in new rectangle remain clear
		if g.B.Test(g.Index(5, 2)) {
			t.Error("expected bit at (5,2) to remain clear")
		}
		if g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3) to remain clear")
		}

		if g.B.Count() != 2 {
			t.Errorf("expected count=2, got %d", g.B.Count())
		}
	})

	t.Run("multiple consecutive shifts", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(7, 5, 2, 2)

		// Shift left 3 times
		g.ShiftRectLeft(7, 5, 2, 2)
		g.ShiftRectLeft(6, 5, 2, 2)
		g.ShiftRectLeft(5, 5, 2, 2)

		// Verify rectangle at (4,5)
		if !g.B.Test(g.Index(4, 5)) {
			t.Error("expected bit at (4,5)")
		}
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5)")
		}
		if !g.B.Test(g.Index(4, 6)) {
			t.Error("expected bit at (4,6)")
		}
		if !g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6)")
		}

		// Verify original positions cleared
		for x := 6; x < 9; x++ {
			for y := 5; y < 7; y++ {
				if g.B.Test(g.Index(x, y)) {
					t.Errorf("expected bit at (%d,%d) to be cleared", x, y)
				}
			}
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift to leftmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 3, 2, 2) // Rectangle starting at col 2, can shift to 0-1

		g.ShiftRectLeft(2, 3, 2, 2)

		// Verify rectangle at (1,3)
		if !g.B.Test(g.Index(1, 3)) {
			t.Error("expected bit at (1,3)")
		}
		if !g.B.Test(g.Index(2, 3)) {
			t.Error("expected bit at (2,3)")
		}
		if !g.B.Test(g.Index(1, 4)) {
			t.Error("expected bit at (1,4)")
		}
		if !g.B.Test(g.Index(2, 4)) {
			t.Error("expected bit at (2,4)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(5, 5, 0, 2)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(5, 5, 2, 0)
	})

	t.Run("panics when target column occupied", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when target column occupied")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 2, 2)
		g.B.SetBit(g.Index(4, 3)) // Occupy target column (x-1=4)

		g.ShiftRectLeft(5, 3, 2, 2)
	})

	t.Run("panics at left edge", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic at left edge (x=0)")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 3, 2, 2) // x=0, cannot shift left

		g.ShiftRectLeft(0, 3, 2, 2)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(9, 3, 2, 2) // x+w=11 exceeds bounds
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(5, 9, 2, 2) // y+h=11 exceeds bounds
	})

	t.Run("panics on negative x", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative x")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(-1, 3, 2, 2)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(7, 5, 2, 2)

		result := g.ShiftRectLeft(7, 5, 2, 2)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(10, 10)
		g2.SetRect(7, 5, 2, 2).
			ShiftRectLeft(7, 5, 2, 2).
			ShiftRectLeft(6, 5, 2, 2)

		if g2.B.Count() != 4 {
			t.Errorf("expected count=4 after chaining, got %d", g2.B.Count())
		}
	})
}

// TestGridShiftRectUp validates Grid.ShiftRectUp() shift operation.
func TestGridShiftRectUp(t *testing.T) {
	t.Run("valid shift with free target row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 2, 2) // 2x2 rectangle at (3,5)

		g.ShiftRectUp(3, 5, 2, 2)

		// Verify rectangle moved to (3,4)
		if !g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4)")
		}
		if !g.B.Test(g.Index(4, 4)) {
			t.Error("expected bit at (4,4)")
		}
		if !g.B.Test(g.Index(3, 5)) {
			t.Error("expected bit at (3,5)")
		}
		if !g.B.Test(g.Index(4, 5)) {
			t.Error("expected bit at (4,5)")
		}

		// Verify bottom row cleared (y+h-1=6)
		if g.B.Test(g.Index(3, 6)) {
			t.Error("expected bit at (3,6) to be cleared")
		}
		if g.B.Test(g.Index(4, 6)) {
			t.Error("expected bit at (4,6) to be cleared")
		}

		// Verify count unchanged
		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift clears bottom row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 2, 2)

		// Verify initial state
		if !g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6) before shift")
		}
		if !g.B.Test(g.Index(6, 6)) {
			t.Error("expected bit at (6,6) before shift")
		}

		g.ShiftRectUp(5, 5, 2, 2)

		// Verify bottom row (y+h-1=6) cleared
		if g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6) cleared after shift")
		}
		if g.B.Test(g.Index(6, 6)) {
			t.Error("expected bit at (6,6) cleared after shift")
		}
	})

	t.Run("shift preserves rectangle data", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		// Create specific pattern in 2x2 rectangle
		g.B.SetBit(g.Index(2, 5)) // top-left
		g.B.SetBit(g.Index(3, 6)) // bottom-right

		g.ShiftRectUp(2, 5, 2, 2)

		// Verify pattern preserved at new location (2,4)
		if !g.B.Test(g.Index(2, 4)) {
			t.Error("expected bit at (2,4) - top-left preserved")
		}
		if !g.B.Test(g.Index(3, 5)) {
			t.Error("expected bit at (3,5) - bottom-right preserved")
		}

		// Verify other cells in new rectangle remain clear
		if g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4) to remain clear")
		}
		if g.B.Test(g.Index(2, 5)) {
			t.Error("expected bit at (2,5) to remain clear")
		}

		if g.B.Count() != 2 {
			t.Errorf("expected count=2, got %d", g.B.Count())
		}
	})

	t.Run("multiple consecutive shifts", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 7, 2, 2)

		// Shift up 3 times
		g.ShiftRectUp(5, 7, 2, 2)
		g.ShiftRectUp(5, 6, 2, 2)
		g.ShiftRectUp(5, 5, 2, 2)

		// Verify rectangle at (5,4)
		if !g.B.Test(g.Index(5, 4)) {
			t.Error("expected bit at (5,4)")
		}
		if !g.B.Test(g.Index(6, 4)) {
			t.Error("expected bit at (6,4)")
		}
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5)")
		}
		if !g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5)")
		}

		// Verify original positions cleared
		for x := 5; x < 7; x++ {
			for y := 6; y < 9; y++ {
				if g.B.Test(g.Index(x, y)) {
					t.Errorf("expected bit at (%d,%d) to be cleared", x, y)
				}
			}
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift to topmost row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 2, 2, 2) // Rectangle starting at row 2, can shift to 0-1

		g.ShiftRectUp(3, 2, 2, 2)

		// Verify rectangle at (3,1)
		if !g.B.Test(g.Index(3, 1)) {
			t.Error("expected bit at (3,1)")
		}
		if !g.B.Test(g.Index(4, 1)) {
			t.Error("expected bit at (4,1)")
		}
		if !g.B.Test(g.Index(3, 2)) {
			t.Error("expected bit at (3,2)")
		}
		if !g.B.Test(g.Index(4, 2)) {
			t.Error("expected bit at (4,2)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(5, 5, 0, 2)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(5, 5, 2, 0)
	})

	t.Run("panics when target row occupied", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when target row occupied")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 2, 2)
		g.B.SetBit(g.Index(3, 4)) // Occupy target row (y-1=4)

		g.ShiftRectUp(3, 5, 2, 2)
	})

	t.Run("panics at top edge", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic at top edge (y=0)")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 0, 2, 2) // y=0, cannot shift up

		g.ShiftRectUp(3, 0, 2, 2)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(9, 5, 2, 2) // x+w=11 exceeds bounds
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(5, 9, 2, 2) // y+h=11 exceeds bounds
	})

	t.Run("panics on negative y", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative y")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(5, -1, 2, 2)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 7, 2, 2)

		result := g.ShiftRectUp(5, 7, 2, 2)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(10, 10)
		g2.SetRect(5, 7, 2, 2).
			ShiftRectUp(5, 7, 2, 2).
			ShiftRectUp(5, 6, 2, 2)

		if g2.B.Count() != 4 {
			t.Errorf("expected count=4 after chaining, got %d", g2.B.Count())
		}
	})
}

// TestGridShiftRectDown validates Grid.ShiftRectDown() shift operation.
func TestGridShiftRectDown(t *testing.T) {
	t.Run("valid shift with free target row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 2, 2) // 2x2 rectangle at (3,3)

		g.ShiftRectDown(3, 3, 2, 2)

		// Verify rectangle moved to (3,4)
		if !g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4)")
		}
		if !g.B.Test(g.Index(4, 4)) {
			t.Error("expected bit at (4,4)")
		}
		if !g.B.Test(g.Index(3, 5)) {
			t.Error("expected bit at (3,5)")
		}
		if !g.B.Test(g.Index(4, 5)) {
			t.Error("expected bit at (4,5)")
		}

		// Verify top row cleared (y=3)
		if g.B.Test(g.Index(3, 3)) {
			t.Error("expected bit at (3,3) to be cleared")
		}
		if g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3) to be cleared")
		}

		// Verify count unchanged
		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift clears top row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 2, 2)

		// Verify initial state
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5) before shift")
		}
		if !g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5) before shift")
		}

		g.ShiftRectDown(5, 5, 2, 2)

		// Verify top row (y=5) cleared
		if g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5) cleared after shift")
		}
		if g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5) cleared after shift")
		}
	})

	t.Run("shift preserves rectangle data", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		// Create specific pattern in 2x2 rectangle
		g.B.SetBit(g.Index(2, 2)) // top-left
		g.B.SetBit(g.Index(3, 3)) // bottom-right

		g.ShiftRectDown(2, 2, 2, 2)

		// Verify pattern preserved at new location (2,3)
		if !g.B.Test(g.Index(2, 3)) {
			t.Error("expected bit at (2,3) - top-left preserved")
		}
		if !g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4) - bottom-right preserved")
		}

		// Verify other cells in new rectangle remain clear
		if g.B.Test(g.Index(3, 3)) {
			t.Error("expected bit at (3,3) to remain clear")
		}
		if g.B.Test(g.Index(2, 4)) {
			t.Error("expected bit at (2,4) to remain clear")
		}

		if g.B.Count() != 2 {
			t.Errorf("expected count=2, got %d", g.B.Count())
		}
	})

	t.Run("multiple consecutive shifts", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 1, 2, 2)

		// Shift down 3 times
		g.ShiftRectDown(5, 1, 2, 2)
		g.ShiftRectDown(5, 2, 2, 2)
		g.ShiftRectDown(5, 3, 2, 2)

		// Verify rectangle at (5,4)
		if !g.B.Test(g.Index(5, 4)) {
			t.Error("expected bit at (5,4)")
		}
		if !g.B.Test(g.Index(6, 4)) {
			t.Error("expected bit at (6,4)")
		}
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5)")
		}
		if !g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5)")
		}

		// Verify original positions cleared
		for x := 5; x < 7; x++ {
			for y := range 4 {
				if g.B.Test(g.Index(x, y)) {
					t.Errorf("expected bit at (%d,%d) to be cleared", x, y)
				}
			}
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift to bottommost row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 6, 2, 2) // Rectangle at rows 6-7, can shift to 8-9

		g.ShiftRectDown(3, 6, 2, 2)

		// Verify rectangle at (3,7)
		if !g.B.Test(g.Index(3, 7)) {
			t.Error("expected bit at (3,7)")
		}
		if !g.B.Test(g.Index(4, 7)) {
			t.Error("expected bit at (4,7)")
		}
		if !g.B.Test(g.Index(3, 8)) {
			t.Error("expected bit at (3,8)")
		}
		if !g.B.Test(g.Index(4, 8)) {
			t.Error("expected bit at (4,8)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(5, 5, 0, 2)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(5, 5, 2, 0)
	})

	t.Run("panics when target row occupied", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when target row occupied")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 2, 2)
		g.B.SetBit(g.Index(3, 5)) // Occupy target row (y+h=5)

		g.ShiftRectDown(3, 3, 2, 2)
	})

	t.Run("panics when target row out of bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when target row out of bounds")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 8, 2, 2) // y+h=10, cannot shift down

		g.ShiftRectDown(3, 8, 2, 2)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(9, 3, 2, 2) // x+w=11 exceeds bounds
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(3, 9, 2, 2) // y+h=11 exceeds bounds
	})

	t.Run("panics on negative y", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative y")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(5, -1, 2, 2)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 1, 2, 2)

		result := g.ShiftRectDown(5, 1, 2, 2)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(10, 10)
		g2.SetRect(5, 1, 2, 2).
			ShiftRectDown(5, 1, 2, 2).
			ShiftRectDown(5, 2, 2, 2)

		if g2.B.Count() != 4 {
			t.Errorf("expected count=4 after chaining, got %d", g2.B.Count())
		}
	})
}
