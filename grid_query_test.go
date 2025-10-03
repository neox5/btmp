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
