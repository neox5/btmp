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

		if !g.IsFree(5, 0, 1, 10) {
			t.Error("expected true for free row")
		}
	})

	t.Run("full row occupied", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 0, 1, 10) // Fill row 5

		if g.IsFree(5, 0, 1, 10) {
			t.Error("expected false for occupied row")
		}
	})

	t.Run("full column free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.IsFree(0, 5, 10, 1) {
			t.Error("expected true for free column")
		}
	})

	t.Run("full column occupied", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 5, 10, 1) // Fill column 5

		if g.IsFree(0, 5, 10, 1) {
			t.Error("expected false for occupied column")
		}
	})

	t.Run("multi-row rectangle free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.IsFree(3, 2, 4, 5) {
			t.Error("expected true for free rectangle")
		}
	})

	t.Run("multi-row rectangle with bit in middle", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(5, 4)) // Middle of (3,2,4,5)

		if g.IsFree(3, 2, 4, 5) {
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
		g.B.SetBit(g.Index(5, 1)) // Just outside (3,2,4,5)

		if !g.IsFree(3, 2, 4, 5) {
			t.Error("expected true, bit is outside rectangle")
		}
	})

	t.Run("detects bit outside on right", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(5, 7)) // Just outside (3,2,4,5) right edge

		if !g.IsFree(3, 2, 4, 5) {
			t.Error("expected true, bit is outside rectangle")
		}
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(-1, 5, 3, 3)
	})

	t.Run("panics on negative c", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative c")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, -1, 3, 3)
	})

	t.Run("panics on negative h", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative h")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 5, -1, 3)
	})

	t.Run("panics on negative w", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative w")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 5, 3, -1)
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 5, 0, 3)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(5, 5, 3, 0)
	})

	t.Run("panics when r+h exceeds rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for r+h > rows")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.IsFree(8, 5, 3, 3) // 8+3=11 > 10
	})

	t.Run("panics when c+w exceeds cols", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for c+w > cols")
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
		g.SetRect(3, 2, 4, 3) // Rectangle at (3,2,4,3)

		if !g.CanShiftRight(3, 2, 4, 3) {
			t.Error("expected true when target column (5) is free")
		}
	})

	t.Run("returns false when target column has bit", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 2, 4, 3)     // Rectangle at (3,2,4,3)
		g.B.SetBit(g.Index(4, 5)) // Set bit in target column (c+w=5)

		if g.CanShiftRight(3, 2, 4, 3) {
			t.Error("expected false when target column has set bit")
		}
	})

	t.Run("returns false when any bit in target column set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 2, 4, 3)     // Rectangle at (3,2,4,3)
		g.B.SetBit(g.Index(6, 5)) // Last row of target column

		if g.CanShiftRight(3, 2, 4, 3) {
			t.Error("expected false when any bit in target column set")
		}
	})

	t.Run("single row shift validation", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 2, 1, 3) // Single row at r=5

		if !g.CanShiftRight(5, 2, 1, 3) {
			t.Error("expected true for single row shift")
		}
	})

	t.Run("single row with occupied target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 2, 1, 3)
		g.B.SetBit(g.Index(5, 5)) // Target column occupied

		if g.CanShiftRight(5, 2, 1, 3) {
			t.Error("expected false when single row target occupied")
		}
	})

	t.Run("multi-row shift with free target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 1, 5, 4) // Multi-row rectangle

		if !g.CanShiftRight(2, 1, 5, 4) {
			t.Error("expected true for multi-row shift with free target")
		}
	})

	t.Run("shift to rightmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 0, 3, 9) // Rectangle ending at col 8, can shift to 9

		if !g.CanShiftRight(0, 0, 3, 9) {
			t.Error("expected true shifting to rightmost column")
		}
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftRight(5, 5, 0, 3)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
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
		g.CanShiftRight(5, 8, 3, 3) // Source rect exceeds bounds (8+3=11)
	})
}

// TestGridCanShiftLeft validates Grid.CanShiftLeft() query operation behavior.
func TestGridCanShiftLeft(t *testing.T) {
	t.Run("returns true when target column is free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 4, 3) // Rectangle at (3,5,4,3)

		if !g.CanShiftLeft(3, 5, 4, 3) {
			t.Error("expected true when target column (4) is free")
		}
	})

	t.Run("returns false when target column has bit", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 4, 3)
		g.B.SetBit(g.Index(4, 4)) // Set bit in target column (c-1=4)

		if g.CanShiftLeft(3, 5, 4, 3) {
			t.Error("expected false when target column has set bit")
		}
	})

	t.Run("returns false when any bit in target column set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 4, 3)
		g.B.SetBit(g.Index(6, 4)) // Last row of target column

		if g.CanShiftLeft(3, 5, 4, 3) {
			t.Error("expected false when any bit in target column set")
		}
	})

	t.Run("single row shift validation", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 1, 3) // Single row at r=5

		if !g.CanShiftLeft(5, 5, 1, 3) {
			t.Error("expected true for single row shift")
		}
	})

	t.Run("single row with occupied target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 1, 3)
		g.B.SetBit(g.Index(5, 4)) // Target column occupied

		if g.CanShiftLeft(5, 5, 1, 3) {
			t.Error("expected false when single row target occupied")
		}
	})

	t.Run("multi-row shift with free target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 5, 5, 4) // Multi-row rectangle

		if !g.CanShiftLeft(2, 5, 5, 4) {
			t.Error("expected true for multi-row shift with free target")
		}
	})

	t.Run("shift to leftmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 1, 3, 5) // Rectangle starting at col 1, can shift to 0

		if !g.CanShiftLeft(0, 1, 3, 5) {
			t.Error("expected true shifting to leftmost column")
		}
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftLeft(5, 5, 0, 3)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
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
		g.CanShiftLeft(5, 9, 3, 3) // Source rect exceeds bounds (9+3=12)
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftLeft(9, 5, 3, 3) // Source rect exceeds bounds (9+3=12)
	})

	t.Run("panics on negative c", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative c")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftLeft(5, -1, 3, 3)
	})
}

// TestGridCanShiftUp validates Grid.CanShiftUp() query operation behavior.
func TestGridCanShiftUp(t *testing.T) {
	t.Run("returns true when target row is free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 3, 4) // Rectangle at (5,3,3,4)

		if !g.CanShiftUp(5, 3, 3, 4) {
			t.Error("expected true when target row (4) is free")
		}
	})

	t.Run("returns false when target row has bit", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 3, 4)
		g.B.SetBit(g.Index(4, 4)) // Set bit in target row (r-1=4)

		if g.CanShiftUp(5, 3, 3, 4) {
			t.Error("expected false when target row has set bit")
		}
	})

	t.Run("returns false when any bit in target row set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 3, 4)
		g.B.SetBit(g.Index(4, 6)) // Last column of target row

		if g.CanShiftUp(5, 3, 3, 4) {
			t.Error("expected false when any bit in target row set")
		}
	})

	t.Run("single column shift validation", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 3, 1) // Single column at c=5

		if !g.CanShiftUp(5, 5, 3, 1) {
			t.Error("expected true for single column shift")
		}
	})

	t.Run("single column with occupied target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 3, 1)
		g.B.SetBit(g.Index(4, 5)) // Target row occupied

		if g.CanShiftUp(5, 5, 3, 1) {
			t.Error("expected false when single column target occupied")
		}
	})

	t.Run("multi-column shift with free target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 2, 4, 5) // Multi-column rectangle

		if !g.CanShiftUp(5, 2, 4, 5) {
			t.Error("expected true for multi-column shift with free target")
		}
	})

	t.Run("shift to topmost row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(1, 0, 3, 5) // Rectangle starting at row 1, can shift to 0

		if !g.CanShiftUp(1, 0, 3, 5) {
			t.Error("expected true shifting to topmost row")
		}
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftUp(5, 5, 0, 3)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
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
		g.CanShiftUp(5, 9, 3, 3) // Source rect exceeds bounds (9+3=12)
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftUp(9, 5, 3, 3) // Source rect exceeds bounds (9+3=12)
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftUp(-1, 5, 3, 3)
	})
}

// TestGridCanShiftDown validates Grid.CanShiftDown() query operation behavior.
func TestGridCanShiftDown(t *testing.T) {
	t.Run("returns true when target row is free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 3, 4) // Rectangle at (3,3,3,4), target row is 6

		if !g.CanShiftDown(3, 3, 3, 4) {
			t.Error("expected true when target row (6) is free")
		}
	})

	t.Run("returns false when target row has bit", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 3, 4)
		g.B.SetBit(g.Index(6, 4)) // Set bit in target row (r+h=6)

		if g.CanShiftDown(3, 3, 3, 4) {
			t.Error("expected false when target row has set bit")
		}
	})

	t.Run("returns false when any bit in target row set", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 3, 4)
		g.B.SetBit(g.Index(6, 6)) // Last column of target row

		if g.CanShiftDown(3, 3, 3, 4) {
			t.Error("expected false when any bit in target row set")
		}
	})

	t.Run("single column shift validation", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 3, 1) // Single column at c=5

		if !g.CanShiftDown(3, 5, 3, 1) {
			t.Error("expected true for single column shift")
		}
	})

	t.Run("single column with occupied target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 3, 1)
		g.B.SetBit(g.Index(6, 5)) // Target row occupied (r+h=6)

		if g.CanShiftDown(3, 5, 3, 1) {
			t.Error("expected false when single column target occupied")
		}
	})

	t.Run("multi-column shift with free target", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 2, 4, 5) // Multi-column rectangle

		if !g.CanShiftDown(2, 2, 4, 5) {
			t.Error("expected true for multi-column shift with free target")
		}
	})

	t.Run("shift to bottommost row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 0, 9, 5) // Rectangle ending at row 8, can shift to 9

		if !g.CanShiftDown(0, 0, 9, 5) {
			t.Error("expected true shifting to bottommost row")
		}
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanShiftDown(5, 5, 0, 3)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
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
		g.CanShiftDown(5, 8, 3, 3) // Source rect exceeds bounds (8+3=11)
	})
}

// TestGridNextFreeCol validates Grid.NextFreeCol() query operation behavior.
func TestGridNextFreeCol(t *testing.T) {
	t.Run("returns first column when row is empty", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		col := g.NextFreeCol(5, 0)
		if col != 0 {
			t.Errorf("expected col=0, got %d", col)
		}
	})

	t.Run("returns starting column when it is free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 0, 1, 3) // Occupy columns 0-2

		col := g.NextFreeCol(5, 3)
		if col != 3 {
			t.Errorf("expected col=3, got %d", col)
		}
	})

	t.Run("skips occupied columns", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 2, 1, 3) // Occupy columns 2-4

		col := g.NextFreeCol(5, 0)
		if col != 0 {
			t.Errorf("expected col=0, got %d", col)
		}

		col = g.NextFreeCol(5, 2)
		if col != 5 {
			t.Errorf("expected col=5, got %d", col)
		}
	})

	t.Run("finds free column after occupied region", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 1, 1, 4) // Row 3, columns 1-4 occupied

		col := g.NextFreeCol(3, 1)
		if col != 5 {
			t.Errorf("expected col=5, got %d", col)
		}
	})

	t.Run("returns -1 when no free column exists", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 0, 1, 10) // Fill entire row

		col := g.NextFreeCol(5, 0)
		if col != -1 {
			t.Errorf("expected col=-1, got %d", col)
		}
	})

	t.Run("returns -1 when starting beyond last free column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 7, 1, 3) // Occupy columns 7-9

		col := g.NextFreeCol(5, 7)
		if col != -1 {
			t.Errorf("expected col=-1, got %d", col)
		}
	})

	t.Run("works at row boundaries", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 5, 1, 5) // First row, columns 5-9 occupied

		col := g.NextFreeCol(0, 0)
		if col != 0 {
			t.Errorf("expected col=0, got %d", col)
		}

		col = g.NextFreeCol(0, 5)
		if col != -1 {
			t.Errorf("expected col=-1, got %d", col)
		}
	})

	t.Run("handles rowSpan scenario", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		// Simulate cell with rowSpan=3 at (0,5,3,2)
		g.SetRect(0, 5, 3, 2) // Rows 0-2, columns 5-6 occupied

		// In row 1, columns 5-6 are occupied by rowSpan from row 0
		col := g.NextFreeCol(1, 5)
		if col != 7 {
			t.Errorf("expected col=7, got %d", col)
		}

		// Starting from column 0 should return 0
		col = g.NextFreeCol(1, 0)
		if col != 0 {
			t.Errorf("expected col=0, got %d", col)
		}
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeCol(-1, 0)
	})

	t.Run("panics on negative c", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative c")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeCol(5, -1)
	})

	t.Run("panics when r >= Rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for r >= Rows")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeCol(10, 0)
	})

	t.Run("panics when c >= Cols", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for c >= Cols")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeCol(5, 10)
	})
}

// TestGridNextFreeColInRange validates Grid.NextFreeColInRange() query operation behavior.
func TestGridNextFreeColInRange(t *testing.T) {
	t.Run("finds free column within range", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 2, 1, 2) // Occupy columns 2-3

		col := g.NextFreeColInRange(5, 0, 10)
		if col != 0 {
			t.Errorf("expected col=0, got %d", col)
		}

		col = g.NextFreeColInRange(5, 2, 10)
		if col != 4 {
			t.Errorf("expected col=4, got %d", col)
		}
	})

	t.Run("returns -1 when no free column in range", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 2, 1, 5) // Occupy columns 2-6

		col := g.NextFreeColInRange(5, 2, 5)
		if col != -1 {
			t.Errorf("expected col=-1, got %d", col)
		}
	})

	t.Run("limits search to available columns", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		// Request count=20 but only 5 columns remain
		col := g.NextFreeColInRange(5, 5, 20)
		if col != 5 {
			t.Errorf("expected col=5, got %d", col)
		}
	})

	t.Run("returns -1 when count exceeds available columns", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 0, 1, 10) // Fill entire row

		col := g.NextFreeColInRange(5, 8, 5) // Only 2 columns available
		if col != -1 {
			t.Errorf("expected col=-1, got %d", col)
		}
	})

	t.Run("finds column in exact range", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 0, 1, 5) // Occupy columns 0-4
		g.SetRect(5, 8, 1, 2) // Occupy columns 8-9

		col := g.NextFreeColInRange(5, 5, 3) // Search columns 5-7
		if col != 5 {
			t.Errorf("expected col=5, got %d", col)
		}
	})

	t.Run("handles rowSpan with limited range", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 3, 3, 4) // Rows 0-2, columns 3-6 occupied

		// In row 1, search columns 3-7 (count=5)
		col := g.NextFreeColInRange(1, 3, 5)
		if col != 7 {
			t.Errorf("expected col=7, got %d", col)
		}

		// Search only within occupied region
		col = g.NextFreeColInRange(1, 3, 4)
		if col != -1 {
			t.Errorf("expected col=-1, got %d", col)
		}
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeColInRange(-1, 0, 5)
	})

	t.Run("panics on negative c", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative c")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeColInRange(5, -1, 5)
	})

	t.Run("panics on zero count", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero count")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeColInRange(5, 0, 0)
	})

	t.Run("panics on negative count", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative count")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeColInRange(5, 0, -1)
	})

	t.Run("panics when r >= Rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for r >= Rows")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeColInRange(10, 0, 5)
	})

	t.Run("panics when c >= Cols", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for c >= Cols")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.NextFreeColInRange(5, 10, 5)
	})
}

// TestGridFreeColsFrom validates Grid.FreeColsFrom() query operation behavior.
func TestGridFreeColsFrom(t *testing.T) {
	t.Run("returns full row width when empty", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		count := g.FreeColsFrom(5, 0)
		if count != 10 {
			t.Errorf("expected count=10, got %d", count)
		}
	})

	t.Run("returns 0 when starting cell is occupied", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 1, 4) // Occupy columns 3-6

		count := g.FreeColsFrom(5, 3)
		if count != 0 {
			t.Errorf("expected count=0, got %d", count)
		}
	})

	t.Run("counts consecutive free columns", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 1, 3) // Occupy columns 5-7

		count := g.FreeColsFrom(5, 0)
		if count != 5 {
			t.Errorf("expected count=5, got %d", count)
		}
	})

	t.Run("stops at first occupied column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(5, 3)) // Single occupied cell at column 3

		count := g.FreeColsFrom(5, 0)
		if count != 3 {
			t.Errorf("expected count=3, got %d", count)
		}
	})

	t.Run("counts to end of row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 0, 1, 4) // Occupy columns 0-3

		count := g.FreeColsFrom(5, 4)
		if count != 6 {
			t.Errorf("expected count=6, got %d", count)
		}
	})

	t.Run("handles rowSpan scenario", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 3, 3, 4) // Rows 0-2, columns 3-6 occupied

		// In row 1, starting from column 0
		count := g.FreeColsFrom(1, 0)
		if count != 3 {
			t.Errorf("expected count=3, got %d", count)
		}

		// Starting from column 7 (after occupied region)
		count = g.FreeColsFrom(1, 7)
		if count != 3 {
			t.Errorf("expected count=3, got %d", count)
		}
	})

	t.Run("returns remaining columns when all free", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		count := g.FreeColsFrom(5, 7)
		if count != 3 {
			t.Errorf("expected count=3, got %d", count)
		}
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.FreeColsFrom(-1, 0)
	})

	t.Run("panics on negative c", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative c")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.FreeColsFrom(5, -1)
	})

	t.Run("panics when r >= Rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for r >= Rows")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.FreeColsFrom(10, 0)
	})

	t.Run("panics when c >= Cols", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for c >= Cols")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.FreeColsFrom(5, 10)
	})
}

// TestGridCanFitWidth validates Grid.CanFitWidth() query operation behavior.
func TestGridCanFitWidth(t *testing.T) {
	t.Run("returns true when width fits in empty row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.CanFitWidth(5, 0, 5) {
			t.Error("expected true for width=5 in empty row")
		}
	})

	t.Run("returns false when width exceeds row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if g.CanFitWidth(5, 0, 15) {
			t.Error("expected false for width=15 in 10-column row")
		}
	})

	t.Run("returns false when width exceeds remaining columns", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if g.CanFitWidth(5, 8, 5) {
			t.Error("expected false for width=5 starting at column 8")
		}
	})

	t.Run("returns false when any cell is occupied", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.B.SetBit(g.Index(5, 3)) // Occupy column 3

		if g.CanFitWidth(5, 0, 5) {
			t.Error("expected false when column 3 is occupied")
		}
	})

	t.Run("returns true when exact width fits", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 0, 1, 5) // Occupy columns 0-4

		if !g.CanFitWidth(5, 5, 5) {
			t.Error("expected true for exact fit in columns 5-9")
		}
	})

	t.Run("returns false when width spans occupied region", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 1, 2) // Occupy columns 3-4

		if g.CanFitWidth(5, 2, 4) {
			t.Error("expected false when width spans occupied columns")
		}
	})

	t.Run("handles rowSpan scenario", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 3, 3, 4) // Rows 0-2, columns 3-6 occupied

		// In row 1, width=3 fits in columns 0-2
		if !g.CanFitWidth(1, 0, 3) {
			t.Error("expected true for width=3 in columns 0-2")
		}

		// Width=4 would span into occupied region
		if g.CanFitWidth(1, 0, 4) {
			t.Error("expected false for width=4 spanning occupied region")
		}

		// Width=3 fits after occupied region
		if !g.CanFitWidth(1, 7, 3) {
			t.Error("expected true for width=3 in columns 7-9")
		}
	})

	t.Run("returns true for width=1 at last column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if !g.CanFitWidth(5, 9, 1) {
			t.Error("expected true for width=1 at last column")
		}
	})

	t.Run("returns false for width=2 at last column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		if g.CanFitWidth(5, 9, 2) {
			t.Error("expected false for width=2 at last column")
		}
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFitWidth(-1, 0, 5)
	})

	t.Run("panics on negative c", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative c")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFitWidth(5, -1, 5)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFitWidth(5, 0, 0)
	})

	t.Run("panics on negative width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative width")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFitWidth(5, 0, -1)
	})

	t.Run("panics when r >= Rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for r >= Rows")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFitWidth(10, 0, 5)
	})

	t.Run("panics when c >= Cols", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for c >= Cols")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFitWidth(5, 10, 5)
	})
}

// TestGridCanFit validates Grid.CanFit() boundary checking behavior.
func TestGridCanFit(t *testing.T) {
	t.Run("returns true, true when rectangle fits exactly", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(0, 0, 10, 10)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true for 10x10 rect in 10x10 grid, got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("returns true, true when rectangle fits with space", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(2, 3, 5, 4)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true for 5x4 rect at (2,3) in 10x10 grid, got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("returns true, true for single cell at origin", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(0, 0, 1, 1)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true for 1x1 rect at (0,0), got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("returns true, true for single cell at corner", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(9, 9, 1, 1)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true for 1x1 rect at (9,9), got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("returns true, true for zero height rectangle", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(5, 5, 0, 5)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true for zero height rectangle, got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("returns true, true for zero width rectangle", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(5, 5, 5, 0)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true for zero width rectangle, got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("returns true, true for zero size rectangle", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(5, 5, 0, 0)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true for zero size rectangle, got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("returns false, true when height exceeds bounds", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(5, 3, 6, 4)
		if fitRow {
			t.Error("expected fitRow=false when r+h > Rows (5+6=11 > 10)")
		}
		if !fitCol {
			t.Error("expected fitCol=true when c+w <= Cols (3+4=7 <= 10)")
		}
	})

	t.Run("returns true, false when width exceeds bounds", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(3, 5, 4, 6)
		if !fitRow {
			t.Error("expected fitRow=true when r+h <= Rows (3+4=7 <= 10)")
		}
		if fitCol {
			t.Error("expected fitCol=false when c+w > Cols (5+6=11 > 10)")
		}
	})

	t.Run("returns false, false when both dimensions exceed bounds", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(8, 8, 5, 5)
		if fitRow {
			t.Error("expected fitRow=false when r+h > Rows (8+5=13 > 10)")
		}
		if fitCol {
			t.Error("expected fitCol=false when c+w > Cols (8+5=13 > 10)")
		}
	})

	t.Run("returns false, true when height exactly exceeds by one", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(0, 0, 11, 5)
		if fitRow {
			t.Error("expected fitRow=false when r+h=11 > 10")
		}
		if !fitCol {
			t.Error("expected fitCol=true when c+w=5 <= 10")
		}
	})

	t.Run("returns true, false when width exactly exceeds by one", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(0, 0, 5, 11)
		if !fitRow {
			t.Error("expected fitRow=true when r+h=5 <= 10")
		}
		if fitCol {
			t.Error("expected fitCol=false when c+w=11 > 10")
		}
	})

	t.Run("does not check cell occupancy", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 3, 5, 4) // Occupy the target area

		// Should still return true, true since CanFit only checks boundaries
		fitRow, fitCol := g.CanFit(2, 3, 5, 4)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true regardless of cell occupancy, got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("returns true, true at exact boundary", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)

		fitRow, fitCol := g.CanFit(5, 5, 5, 5)
		if !fitRow || !fitCol {
			t.Errorf("expected true, true when r+h=10 and c+w=10, got %v, %v", fitRow, fitCol)
		}
	})

	t.Run("independent row and column checks", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 15)

		// Height exceeds, width fits
		fitRow, fitCol := g.CanFit(8, 5, 5, 5)
		if fitRow {
			t.Error("expected fitRow=false when r+h=13 > 10")
		}
		if !fitCol {
			t.Error("expected fitCol=true when c+w=10 <= 15")
		}

		// Height fits, width exceeds
		fitRow, fitCol = g.CanFit(5, 12, 5, 5)
		if !fitRow {
			t.Error("expected fitRow=true when r+h=10 <= 10")
		}
		if fitCol {
			t.Error("expected fitCol=false when c+w=17 > 15")
		}
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFit(-1, 5, 3, 3)
	})

	t.Run("panics on negative c", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative c")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFit(5, -1, 3, 3)
	})

	t.Run("panics on negative h", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative h")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFit(5, 5, -1, 3)
	})

	t.Run("panics on negative w", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative w")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFit(5, 5, 3, -1)
	})

	t.Run("panics when r >= Rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for r >= Rows")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFit(10, 5, 3, 3)
	})

	t.Run("panics when c >= Cols", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for c >= Cols")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.CanFit(5, 10, 3, 3)
	})
}

// TestGridAllGrid validates Grid.AllGrid() query operation behavior.
func TestGridAllGrid(t *testing.T) {
	t.Run("returns true when all bits set", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 5)
		g.B.SetAll()

		if !g.AllGrid() {
			t.Error("expected true when all bits set")
		}
	})

	t.Run("returns false when no bits set", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 5)

		if g.AllGrid() {
			t.Error("expected false when no bits set")
		}
	})

	t.Run("returns false when single bit clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 5)
		g.B.SetAll()
		g.B.ClearBit(g.Index(2, 2))

		if g.AllGrid() {
			t.Error("expected false when single bit clear")
		}
	})

	t.Run("returns false when first bit clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 5)
		g.B.SetAll()
		g.B.ClearBit(0)

		if g.AllGrid() {
			t.Error("expected false when first bit clear")
		}
	})

	t.Run("returns false when last bit clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 5)
		g.B.SetAll()
		g.B.ClearBit(g.Index(4, 4))

		if g.AllGrid() {
			t.Error("expected false when last bit clear")
		}
	})

	t.Run("returns false for empty grid zero rows", func(t *testing.T) {
		g := btmp.NewGridWithSize(0, 10)

		if g.AllGrid() {
			t.Error("expected false for empty grid (0 rows)")
		}
	})

	t.Run("returns false for empty grid zero cols", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 0)

		if g.AllGrid() {
			t.Error("expected false for empty grid (0 cols)")
		}
	})

	t.Run("returns false for empty grid zero both", func(t *testing.T) {
		g := btmp.NewGrid()

		if g.AllGrid() {
			t.Error("expected false for empty grid (0x0)")
		}
	})

	t.Run("single cell set", func(t *testing.T) {
		g := btmp.NewGridWithSize(1, 1)
		g.B.SetBit(0)

		if !g.AllGrid() {
			t.Error("expected true for single set cell")
		}
	})

	t.Run("single cell clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(1, 1)

		if g.AllGrid() {
			t.Error("expected false for single clear cell")
		}
	})

	t.Run("large grid all set", func(t *testing.T) {
		g := btmp.NewGridWithSize(100, 100)
		g.B.SetAll()

		if !g.AllGrid() {
			t.Error("expected true for large grid all set")
		}
	})

	t.Run("large grid one clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(100, 100)
		g.B.SetAll()
		g.B.ClearBit(g.Index(50, 50))

		if g.AllGrid() {
			t.Error("expected false for large grid with one clear")
		}
	})
}

// TestGridAllRow validates Grid.AllRow() query operation behavior.
func TestGridAllRow(t *testing.T) {
	t.Run("returns true when entire row set", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		g.SetRect(2, 0, 1, 10) // Fill row 2

		if !g.AllRow(2) {
			t.Error("expected true when entire row set")
		}
	})

	t.Run("returns false when entire row clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)

		if g.AllRow(2) {
			t.Error("expected false when entire row clear")
		}
	})

	t.Run("returns false when single bit clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		g.SetRect(2, 0, 1, 10)
		g.B.ClearBit(g.Index(2, 5)) // Clear one bit in row 2

		if g.AllRow(2) {
			t.Error("expected false when single bit clear")
		}
	})

	t.Run("returns false when first bit clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		g.SetRect(2, 0, 1, 10)
		g.B.ClearBit(g.Index(2, 0))

		if g.AllRow(2) {
			t.Error("expected false when first bit clear")
		}
	})

	t.Run("returns false when last bit clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		g.SetRect(2, 0, 1, 10)
		g.B.ClearBit(g.Index(2, 9))

		if g.AllRow(2) {
			t.Error("expected false when last bit clear")
		}
	})

	t.Run("first row all set", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		g.SetRect(0, 0, 1, 10)

		if !g.AllRow(0) {
			t.Error("expected true for first row all set")
		}
	})

	t.Run("last row all set", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		g.SetRect(4, 0, 1, 10)

		if !g.AllRow(4) {
			t.Error("expected true for last row all set")
		}
	})

	t.Run("other rows do not affect result", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		g.B.SetAll()             // Set all bits
		g.ClearRect(2, 0, 1, 10) // Clear row 2

		if g.AllRow(2) {
			t.Error("expected false for cleared row")
		}
		if !g.AllRow(1) {
			t.Error("expected true for adjacent row")
		}
		if !g.AllRow(3) {
			t.Error("expected true for adjacent row")
		}
	})

	t.Run("single column grid", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 1)
		g.B.SetBit(g.Index(2, 0))

		if !g.AllRow(2) {
			t.Error("expected true for single column row set")
		}
	})

	t.Run("single column grid clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 1)

		if g.AllRow(2) {
			t.Error("expected false for single column row clear")
		}
	})

	t.Run("returns false for empty row zero cols", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 0)

		if g.AllRow(0) {
			t.Error("expected false for empty row (0 cols)")
		}
	})

	t.Run("wide row all set", func(t *testing.T) {
		g := btmp.NewGridWithSize(3, 200)
		g.SetRect(1, 0, 1, 200)

		if !g.AllRow(1) {
			t.Error("expected true for wide row all set")
		}
	})

	t.Run("wide row one clear", func(t *testing.T) {
		g := btmp.NewGridWithSize(3, 200)
		g.SetRect(1, 0, 1, 200)
		g.B.ClearBit(g.Index(1, 100))

		if g.AllRow(1) {
			t.Error("expected false for wide row with one clear")
		}
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(5, 10)
		g.AllRow(-1)
	})

	t.Run("panics on r equals rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for r >= Rows()")
			}
		}()
		g := btmp.NewGridWithSize(5, 10)
		g.AllRow(5)
	})

	t.Run("panics on r exceeds rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for r >= Rows()")
			}
		}()
		g := btmp.NewGridWithSize(5, 10)
		g.AllRow(10)
	})
}
