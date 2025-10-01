package btmp_test

import (
	"testing"

	"github.com/neox5/btmp"
)

// TestGridEnsureCols validates Grid.EnsureCols() behavior.
// Tests:
//   - No-op when cols <= current Cols()
//   - Grows columns and repositions existing data
//   - Grows from zero columns
//   - Validates negative input panics
//   - Returns *Grid for chaining
func TestGridEnsureCols(t *testing.T) {
	t.Run("no-op when cols <= current", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 5)
		g.SetRect(0, 0, 10, 5) // Fill grid

		g.EnsureCols(5)

		if g.Cols() != 10 {
			t.Errorf("expected cols=10, got %d", g.Cols())
		}
		if g.Rows() != 5 {
			t.Errorf("expected rows=5, got %d", g.Rows())
		}
		// Verify data unchanged
		if g.B.Count() != 50 {
			t.Errorf("expected count=50, got %d", g.B.Count())
		}
	})

	t.Run("no-op when cols == current", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 5)
		g.SetRect(0, 0, 10, 5)

		g.EnsureCols(10)

		if g.Cols() != 10 {
			t.Errorf("expected cols=10, got %d", g.Cols())
		}
		if g.B.Count() != 50 {
			t.Errorf("expected count=50, got %d", g.B.Count())
		}
	})

	t.Run("grows columns and repositions data", func(t *testing.T) {
		g := btmp.NewGridWithSize(3, 2)
		// Set pattern:
		// Row 0: [1,0,1]
		// Row 1: [0,1,0]
		g.B.SetBit(0) // (0,0)
		g.B.SetBit(2) // (2,0)
		g.B.SetBit(4) // (1,1)

		g.EnsureCols(5)

		if g.Cols() != 5 {
			t.Errorf("expected cols=5, got %d", g.Cols())
		}
		if g.Rows() != 2 {
			t.Errorf("expected rows=2, got %d", g.Rows())
		}

		// Verify data repositioned correctly:
		// Row 0: [1,0,1,0,0]
		// Row 1: [0,1,0,0,0]
		if !g.B.Test(g.Index(0, 0)) {
			t.Error("expected bit at (0,0)")
		}
		if !g.B.Test(g.Index(2, 0)) {
			t.Error("expected bit at (2,0)")
		}
		if !g.B.Test(g.Index(1, 1)) {
			t.Error("expected bit at (1,1)")
		}

		// Verify new columns are zero
		if g.B.Test(g.Index(3, 0)) {
			t.Error("expected zero at (3,0)")
		}
		if g.B.Test(g.Index(4, 0)) {
			t.Error("expected zero at (4,0)")
		}
		if g.B.Test(g.Index(3, 1)) {
			t.Error("expected zero at (3,1)")
		}
		if g.B.Test(g.Index(4, 1)) {
			t.Error("expected zero at (4,1)")
		}

		// Total count should remain 3
		if g.B.Count() != 3 {
			t.Errorf("expected count=3, got %d", g.B.Count())
		}
	})

	t.Run("grows from zero columns", func(t *testing.T) {
		g := btmp.NewGrid()

		g.EnsureCols(5)

		if g.Cols() != 5 {
			t.Errorf("expected cols=5, got %d", g.Cols())
		}
		if g.Rows() != 0 {
			t.Errorf("expected rows=0, got %d", g.Rows())
		}
		if g.B.Len() != 0 {
			t.Errorf("expected bitmap len=0, got %d", g.B.Len())
		}
	})

	t.Run("grows preserves multiple rows", func(t *testing.T) {
		g := btmp.NewGridWithSize(2, 3)
		// Row 0: [1,0]
		// Row 1: [0,1]
		// Row 2: [1,1]
		g.B.SetBit(g.Index(0, 0))
		g.B.SetBit(g.Index(1, 1))
		g.B.SetBit(g.Index(0, 2))
		g.B.SetBit(g.Index(1, 2))

		g.EnsureCols(4)

		if g.Cols() != 4 {
			t.Errorf("expected cols=4, got %d", g.Cols())
		}
		if g.Rows() != 3 {
			t.Errorf("expected rows=3, got %d", g.Rows())
		}

		// Verify all original data preserved
		if !g.B.Test(g.Index(0, 0)) {
			t.Error("expected bit at (0,0)")
		}
		if !g.B.Test(g.Index(1, 1)) {
			t.Error("expected bit at (1,1)")
		}
		if !g.B.Test(g.Index(0, 2)) {
			t.Error("expected bit at (0,2)")
		}
		if !g.B.Test(g.Index(1, 2)) {
			t.Error("expected bit at (1,2)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})

	t.Run("panics on negative cols", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative cols")
			}
		}()
		g := btmp.NewGridWithSize(5, 5)
		g.EnsureCols(-1)
	})

	t.Run("returns grid for chaining", func(t *testing.T) {
		g := btmp.NewGridWithSize(3, 2)
		result := g.EnsureCols(5)

		if result != g {
			t.Error("expected same grid instance")
		}

		// Verify chaining works
		g2 := btmp.NewGridWithSize(2, 2).
			EnsureCols(4).
			EnsureCols(6).
			EnsureCols(3) // no-op

		if g2.Cols() != 6 {
			t.Errorf("expected cols=6, got %d", g2.Cols())
		}
	})

	t.Run("handles large growth", func(t *testing.T) {
		g := btmp.NewGridWithSize(2, 100)
		// Set first and last row
		g.B.SetBit(g.Index(0, 0))
		g.B.SetBit(g.Index(1, 0))
		g.B.SetBit(g.Index(0, 99))
		g.B.SetBit(g.Index(1, 99))

		g.EnsureCols(100)

		if g.Cols() != 100 {
			t.Errorf("expected cols=100, got %d", g.Cols())
		}
		if g.Rows() != 100 {
			t.Errorf("expected rows=100, got %d", g.Rows())
		}

		// Verify corner bits preserved
		if !g.B.Test(g.Index(0, 0)) {
			t.Error("expected bit at (0,0)")
		}
		if !g.B.Test(g.Index(1, 0)) {
			t.Error("expected bit at (1,0)")
		}
		if !g.B.Test(g.Index(0, 99)) {
			t.Error("expected bit at (0,99)")
		}
		if !g.B.Test(g.Index(1, 99)) {
			t.Error("expected bit at (1,99)")
		}

		if g.B.Count() != 4 {
			t.Errorf("expected count=4, got %d", g.B.Count())
		}
	})
}
