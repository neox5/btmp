package btmp_test

import (
	"testing"

	"github.com/neox5/btmp"
)

// TestGridShiftRectRight validates Grid.ShiftRectRight() shift operation.
func TestGridShiftRectRight(t *testing.T) {
	t.Run("valid shift with free target column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 3, 2, 2) // 2x2 rectangle at (3,3)

		g.ShiftRectRight(3, 3, 2, 2)

		// Verify rectangle moved to (3,4)
		if !g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4)")
		}
		if !g.B.Test(g.Index(3, 5)) {
			t.Error("expected bit at (3,5)")
		}
		if !g.B.Test(g.Index(4, 4)) {
			t.Error("expected bit at (4,4)")
		}
		if !g.B.Test(g.Index(4, 5)) {
			t.Error("expected bit at (4,5)")
		}

		// Verify leftmost column cleared (c=3)
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

	t.Run("shift clears leftmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 2, 2)

		// Verify initial state
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5) before shift")
		}
		if !g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5) before shift")
		}

		g.ShiftRectRight(5, 5, 2, 2)

		// Verify leftmost column (c=5) cleared
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

		g.ShiftRectRight(2, 2, 2, 2)

		// Verify pattern preserved at new location (2,3)
		if !g.B.Test(g.Index(2, 3)) {
			t.Error("expected bit at (2,3) - top-left preserved")
		}
		if !g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4) - bottom-right preserved")
		}

		// Verify other cells in new rectangle remain clear
		if g.B.Test(g.Index(2, 4)) {
			t.Error("expected bit at (2,4) to remain clear")
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
		g.SetRect(5, 1, 2, 2)

		// Shift right 3 times
		g.ShiftRectRight(5, 1, 2, 2)
		g.ShiftRectRight(5, 2, 2, 2)
		g.ShiftRectRight(5, 3, 2, 2)

		// Verify rectangle at (5,4)
		if !g.B.Test(g.Index(5, 4)) {
			t.Error("expected bit at (5,4)")
		}
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5)")
		}
		if !g.B.Test(g.Index(6, 4)) {
			t.Error("expected bit at (6,4)")
		}
		if !g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5)")
		}

		// Verify original positions cleared
		for c := range 4 {
			for r := 5; r < 7; r++ {
				if g.B.Test(g.Index(r, c)) {
					t.Errorf("expected bit at (%d,%d) to be cleared", r, c)
				}
			}
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift to rightmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 6, 2, 2) // Rectangle ending at col 7, can shift to 8-9

		g.ShiftRectRight(3, 6, 2, 2)

		// Verify rectangle at (3,7)
		if !g.B.Test(g.Index(3, 7)) {
			t.Error("expected bit at (3,7)")
		}
		if !g.B.Test(g.Index(3, 8)) {
			t.Error("expected bit at (3,8)")
		}
		if !g.B.Test(g.Index(4, 7)) {
			t.Error("expected bit at (4,7)")
		}
		if !g.B.Test(g.Index(4, 8)) {
			t.Error("expected bit at (4,8)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(5, 5, 0, 2)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
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
		g.B.SetBit(g.Index(3, 5)) // Occupy target column (c+w=5)

		g.ShiftRectRight(3, 3, 2, 2)
	})

	t.Run("panics when target column out of bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when target column out of bounds")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 8, 2, 2) // c+w=10, cannot shift right

		g.ShiftRectRight(3, 8, 2, 2)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(3, 9, 2, 2) // c+w=11 exceeds bounds
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(9, 3, 2, 2) // r+h=11 exceeds bounds
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectRight(-1, 3, 2, 2)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 1, 2, 2)

		result := g.ShiftRectRight(5, 1, 2, 2)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(10, 10)
		g2.SetRect(5, 1, 2, 2).
			ShiftRectRight(5, 1, 2, 2).
			ShiftRectRight(5, 2, 2, 2)

		if g2.B.Count() != 4 {
			t.Errorf("expected count=4 after chaining, got %d", g2.B.Count())
		}
	})
}

// TestGridShiftRectLeft validates Grid.ShiftRectLeft() shift operation.
func TestGridShiftRectLeft(t *testing.T) {
	t.Run("valid shift with free target column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 5, 2, 2) // 2x2 rectangle at (3,5)

		g.ShiftRectLeft(3, 5, 2, 2)

		// Verify rectangle moved to (3,4)
		if !g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4)")
		}
		if !g.B.Test(g.Index(3, 5)) {
			t.Error("expected bit at (3,5)")
		}
		if !g.B.Test(g.Index(4, 4)) {
			t.Error("expected bit at (4,4)")
		}
		if !g.B.Test(g.Index(4, 5)) {
			t.Error("expected bit at (4,5)")
		}

		// Verify rightmost column cleared (c+w-1=6)
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

	t.Run("shift clears rightmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 2, 2)

		// Verify initial state
		if !g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6) before shift")
		}
		if !g.B.Test(g.Index(6, 6)) {
			t.Error("expected bit at (6,6) before shift")
		}

		g.ShiftRectLeft(5, 5, 2, 2)

		// Verify rightmost column (c+w-1=6) cleared
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

		g.ShiftRectLeft(2, 5, 2, 2)

		// Verify pattern preserved at new location (2,4)
		if !g.B.Test(g.Index(2, 4)) {
			t.Error("expected bit at (2,4) - top-left preserved")
		}
		if !g.B.Test(g.Index(3, 5)) {
			t.Error("expected bit at (3,5) - bottom-right preserved")
		}

		// Verify other cells in new rectangle remain clear
		if g.B.Test(g.Index(2, 5)) {
			t.Error("expected bit at (2,5) to remain clear")
		}
		if g.B.Test(g.Index(3, 4)) {
			t.Error("expected bit at (3,4) to remain clear")
		}

		if g.B.Count() != 2 {
			t.Errorf("expected count=2, got %d", g.B.Count())
		}
	})

	t.Run("multiple consecutive shifts", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 7, 2, 2)

		// Shift left 3 times
		g.ShiftRectLeft(5, 7, 2, 2)
		g.ShiftRectLeft(5, 6, 2, 2)
		g.ShiftRectLeft(5, 5, 2, 2)

		// Verify rectangle at (5,4)
		if !g.B.Test(g.Index(5, 4)) {
			t.Error("expected bit at (5,4)")
		}
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5)")
		}
		if !g.B.Test(g.Index(6, 4)) {
			t.Error("expected bit at (6,4)")
		}
		if !g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5)")
		}

		// Verify original positions cleared
		for c := 6; c < 9; c++ {
			for r := 5; r < 7; r++ {
				if g.B.Test(g.Index(r, c)) {
					t.Errorf("expected bit at (%d,%d) to be cleared", r, c)
				}
			}
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift to leftmost column", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 2, 2, 2) // Rectangle starting at col 2, can shift to 0-1

		g.ShiftRectLeft(3, 2, 2, 2)

		// Verify rectangle at (3,1)
		if !g.B.Test(g.Index(3, 1)) {
			t.Error("expected bit at (3,1)")
		}
		if !g.B.Test(g.Index(3, 2)) {
			t.Error("expected bit at (3,2)")
		}
		if !g.B.Test(g.Index(4, 1)) {
			t.Error("expected bit at (4,1)")
		}
		if !g.B.Test(g.Index(4, 2)) {
			t.Error("expected bit at (4,2)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(5, 5, 0, 2)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
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
		g.SetRect(3, 5, 2, 2)
		g.B.SetBit(g.Index(3, 4)) // Occupy target column (c-1=4)

		g.ShiftRectLeft(3, 5, 2, 2)
	})

	t.Run("panics at left edge", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic at left edge (c=0)")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(3, 0, 2, 2) // c=0, cannot shift left

		g.ShiftRectLeft(3, 0, 2, 2)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(3, 9, 2, 2) // c+w=11 exceeds bounds
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(9, 5, 2, 2) // r+h=11 exceeds bounds
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectLeft(-1, 3, 2, 2)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 7, 2, 2)

		result := g.ShiftRectLeft(5, 7, 2, 2)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(10, 10)
		g2.SetRect(5, 7, 2, 2).
			ShiftRectLeft(5, 7, 2, 2).
			ShiftRectLeft(5, 6, 2, 2)

		if g2.B.Count() != 4 {
			t.Errorf("expected count=4 after chaining, got %d", g2.B.Count())
		}
	})
}

// TestGridShiftRectUp validates Grid.ShiftRectUp() shift operation.
func TestGridShiftRectUp(t *testing.T) {
	t.Run("valid shift with free target row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 3, 2, 2) // 2x2 rectangle at (5,3)

		g.ShiftRectUp(5, 3, 2, 2)

		// Verify rectangle moved to (4,3)
		if !g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3)")
		}
		if !g.B.Test(g.Index(4, 4)) {
			t.Error("expected bit at (4,4)")
		}
		if !g.B.Test(g.Index(5, 3)) {
			t.Error("expected bit at (5,3)")
		}
		if !g.B.Test(g.Index(5, 4)) {
			t.Error("expected bit at (5,4)")
		}

		// Verify bottom row cleared (r+h-1=6)
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

	t.Run("shift clears bottom row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 2, 2)

		// Verify initial state
		if !g.B.Test(g.Index(6, 5)) {
			t.Error("expected bit at (6,5) before shift")
		}
		if !g.B.Test(g.Index(6, 6)) {
			t.Error("expected bit at (6,6) before shift")
		}

		g.ShiftRectUp(5, 5, 2, 2)

		// Verify bottom row (r+h-1=6) cleared
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

		g.ShiftRectUp(5, 2, 2, 2)

		// Verify pattern preserved at new location (4,2)
		if !g.B.Test(g.Index(4, 2)) {
			t.Error("expected bit at (4,2) - top-left preserved")
		}
		if !g.B.Test(g.Index(5, 3)) {
			t.Error("expected bit at (5,3) - bottom-right preserved")
		}

		// Verify other cells in new rectangle remain clear
		if g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3) to remain clear")
		}
		if g.B.Test(g.Index(5, 2)) {
			t.Error("expected bit at (5,2) to remain clear")
		}

		if g.B.Count() != 2 {
			t.Errorf("expected count=2, got %d", g.B.Count())
		}
	})

	t.Run("multiple consecutive shifts", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(7, 5, 2, 2)

		// Shift up 3 times
		g.ShiftRectUp(7, 5, 2, 2)
		g.ShiftRectUp(6, 5, 2, 2)
		g.ShiftRectUp(5, 5, 2, 2)

		// Verify rectangle at (4,5)
		if !g.B.Test(g.Index(4, 5)) {
			t.Error("expected bit at (4,5)")
		}
		if !g.B.Test(g.Index(4, 6)) {
			t.Error("expected bit at (4,6)")
		}
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5)")
		}
		if !g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6)")
		}

		// Verify original positions cleared
		for c := 5; c < 7; c++ {
			for r := 6; r < 9; r++ {
				if g.B.Test(g.Index(r, c)) {
					t.Errorf("expected bit at (%d,%d) to be cleared", r, c)
				}
			}
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift to topmost row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(2, 3, 2, 2) // Rectangle starting at row 2, can shift to 0-1

		g.ShiftRectUp(2, 3, 2, 2)

		// Verify rectangle at (1,3)
		if !g.B.Test(g.Index(1, 3)) {
			t.Error("expected bit at (1,3)")
		}
		if !g.B.Test(g.Index(1, 4)) {
			t.Error("expected bit at (1,4)")
		}
		if !g.B.Test(g.Index(2, 3)) {
			t.Error("expected bit at (2,3)")
		}
		if !g.B.Test(g.Index(2, 4)) {
			t.Error("expected bit at (2,4)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(5, 5, 0, 2)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
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
		g.SetRect(5, 3, 2, 2)
		g.B.SetBit(g.Index(4, 3)) // Occupy target row (r-1=4)

		g.ShiftRectUp(5, 3, 2, 2)
	})

	t.Run("panics at top edge", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic at top edge (r=0)")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(0, 3, 2, 2) // r=0, cannot shift up

		g.ShiftRectUp(0, 3, 2, 2)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(5, 9, 2, 2) // c+w=11 exceeds bounds
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(9, 5, 2, 2) // r+h=11 exceeds bounds
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectUp(-1, 5, 2, 2)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(7, 5, 2, 2)

		result := g.ShiftRectUp(7, 5, 2, 2)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(10, 10)
		g2.SetRect(7, 5, 2, 2).
			ShiftRectUp(7, 5, 2, 2).
			ShiftRectUp(6, 5, 2, 2)

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

		// Verify rectangle moved to (4,3)
		if !g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3)")
		}
		if !g.B.Test(g.Index(4, 4)) {
			t.Error("expected bit at (4,4)")
		}
		if !g.B.Test(g.Index(5, 3)) {
			t.Error("expected bit at (5,3)")
		}
		if !g.B.Test(g.Index(5, 4)) {
			t.Error("expected bit at (5,4)")
		}

		// Verify top row cleared (r=3)
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

	t.Run("shift clears top row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(5, 5, 2, 2)

		// Verify initial state
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5) before shift")
		}
		if !g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6) before shift")
		}

		g.ShiftRectDown(5, 5, 2, 2)

		// Verify top row (r=5) cleared
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

		g.ShiftRectDown(2, 2, 2, 2)

		// Verify pattern preserved at new location (3,2)
		if !g.B.Test(g.Index(3, 2)) {
			t.Error("expected bit at (3,2) - top-left preserved")
		}
		if !g.B.Test(g.Index(4, 3)) {
			t.Error("expected bit at (4,3) - bottom-right preserved")
		}

		// Verify other cells in new rectangle remain clear
		if g.B.Test(g.Index(3, 3)) {
			t.Error("expected bit at (3,3) to remain clear")
		}
		if g.B.Test(g.Index(4, 2)) {
			t.Error("expected bit at (4,2) to remain clear")
		}

		if g.B.Count() != 2 {
			t.Errorf("expected count=2, got %d", g.B.Count())
		}
	})

	t.Run("multiple consecutive shifts", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(1, 5, 2, 2)

		// Shift down 3 times
		g.ShiftRectDown(1, 5, 2, 2)
		g.ShiftRectDown(2, 5, 2, 2)
		g.ShiftRectDown(3, 5, 2, 2)

		// Verify rectangle at (4,5)
		if !g.B.Test(g.Index(4, 5)) {
			t.Error("expected bit at (4,5)")
		}
		if !g.B.Test(g.Index(4, 6)) {
			t.Error("expected bit at (4,6)")
		}
		if !g.B.Test(g.Index(5, 5)) {
			t.Error("expected bit at (5,5)")
		}
		if !g.B.Test(g.Index(5, 6)) {
			t.Error("expected bit at (5,6)")
		}

		// Verify original positions cleared
		for c := 5; c < 7; c++ {
			for r := range 4 {
				if g.B.Test(g.Index(r, c)) {
					t.Errorf("expected bit at (%d,%d) to be cleared", r, c)
				}
			}
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("shift to bottommost row", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(6, 3, 2, 2) // Rectangle at rows 6-7, can shift to 8-9

		g.ShiftRectDown(6, 3, 2, 2)

		// Verify rectangle at (7,3)
		if !g.B.Test(g.Index(7, 3)) {
			t.Error("expected bit at (7,3)")
		}
		if !g.B.Test(g.Index(7, 4)) {
			t.Error("expected bit at (7,4)")
		}
		if !g.B.Test(g.Index(8, 3)) {
			t.Error("expected bit at (8,3)")
		}
		if !g.B.Test(g.Index(8, 4)) {
			t.Error("expected bit at (8,4)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on zero height", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero height")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(5, 5, 0, 2)
	})

	t.Run("panics on zero width", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for zero width")
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
		g.B.SetBit(g.Index(5, 3)) // Occupy target row (r+h=5)

		g.ShiftRectDown(3, 3, 2, 2)
	})

	t.Run("panics when target row out of bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic when target row out of bounds")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(8, 3, 2, 2) // r+h=10, cannot shift down

		g.ShiftRectDown(8, 3, 2, 2)
	})

	t.Run("panics on invalid source rectangle x bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(3, 9, 2, 2) // c+w=11 exceeds bounds
	})

	t.Run("panics on invalid source rectangle y bounds", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for invalid source rectangle")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(9, 3, 2, 2) // r+h=11 exceeds bounds
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(10, 10)
		g.ShiftRectDown(-1, 5, 2, 2)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 10)
		g.SetRect(1, 5, 2, 2)

		result := g.ShiftRectDown(1, 5, 2, 2)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(10, 10)
		g2.SetRect(1, 5, 2, 2).
			ShiftRectDown(1, 5, 2, 2).
			ShiftRectDown(2, 5, 2, 2)

		if g2.B.Count() != 4 {
			t.Errorf("expected count=4 after chaining, got %d", g2.B.Count())
		}
	})
}
