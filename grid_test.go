package btmp_test

import (
	"testing"

	"github.com/neox5/btmp"
)

// TestNewGridWithSize validates NewGridWithSize constructor behavior.
func TestNewGridWithSize(t *testing.T) {
	t.Run("creates grid with specified dimensions", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		if g.Cols() != 10 {
			t.Errorf("expected cols=10, got %d", g.Cols())
		}
		if g.Rows() != 5 {
			t.Errorf("expected rows=5, got %d", g.Rows())
		}
		if g.B.Len() != 50 {
			t.Errorf("expected bitmap len=50, got %d", g.B.Len())
		}
	})

	t.Run("accepts zero rows", func(t *testing.T) {
		g := btmp.NewGridWithSize(0, 10)
		if g.Cols() != 10 {
			t.Errorf("expected cols=10, got %d", g.Cols())
		}
		if g.Rows() != 0 {
			t.Errorf("expected rows=0, got %d", g.Rows())
		}
		if g.B.Len() != 0 {
			t.Errorf("expected bitmap len=0, got %d", g.B.Len())
		}
	})

	t.Run("accepts zero columns", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 0)
		if g.Cols() != 0 {
			t.Errorf("expected cols=0, got %d", g.Cols())
		}
		if g.B.Len() != 0 {
			t.Errorf("expected bitmap len=0, got %d", g.B.Len())
		}
	})

	t.Run("accepts both dimensions zero", func(t *testing.T) {
		g := btmp.NewGridWithSize(0, 0)
		if g.Cols() != 0 {
			t.Errorf("expected cols=0, got %d", g.Cols())
		}
		if g.B.Len() != 0 {
			t.Errorf("expected bitmap len=0, got %d", g.B.Len())
		}
	})

	t.Run("panics on negative rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative rows")
			}
		}()
		btmp.NewGridWithSize(-1, 10)
	})

	t.Run("panics on negative columns", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative columns")
			}
		}()
		btmp.NewGridWithSize(5, -1)
	})

	t.Run("panics on overflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for overflow")
			}
		}()
		// MaxInt rows × 2 cols should overflow
		btmp.NewGridWithSize(int(^uint(0)>>1), 2)
	})

	t.Run("bitmap initialized to zeros", func(t *testing.T) {
		g := btmp.NewGridWithSize(4, 8)
		if g.B.Any() {
			t.Error("expected all bits to be zero")
		}
		if g.B.Count() != 0 {
			t.Errorf("expected count=0, got %d", g.B.Count())
		}
	})
}

// TestGridCols validates Grid.Cols() accessor behavior.
func TestGridCols(t *testing.T) {
	t.Run("returns correct column count", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		if g.Cols() != 10 {
			t.Errorf("expected cols=10, got %d", g.Cols())
		}
	})

	t.Run("returns 0 for empty grid", func(t *testing.T) {
		g := btmp.NewGrid()
		if g.Cols() != 0 {
			t.Errorf("expected cols=0, got %d", g.Cols())
		}
	})
}

// TestGridRows validates Grid.Rows() accessor behavior.
func TestGridRows(t *testing.T) {
	t.Run("returns correct row count", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)
		if g.Rows() != 5 {
			t.Errorf("expected rows=5, got %d", g.Rows())
		}
	})

	t.Run("returns 0 for empty grid", func(t *testing.T) {
		g := btmp.NewGrid()
		if g.Rows() != 0 {
			t.Errorf("expected rows=0, got %d", g.Rows())
		}
	})

	t.Run("returns 0 when cols is 0", func(t *testing.T) {
		g := btmp.NewGridWithSize(0, 0)
		if g.Rows() != 0 {
			t.Errorf("expected rows=0, got %d", g.Rows())
		}
	})
}

// TestGridIndex validates Grid.Index() coordinate conversion.
func TestGridIndex(t *testing.T) {
	t.Run("calculates correct row-major index", func(t *testing.T) {
		g := btmp.NewGridWithSize(5, 10)

		// Test various coordinates
		tests := []struct {
			r, c int
			want int
		}{
			{0, 0, 0},
			{0, 5, 5},
			{1, 0, 10},
			{2, 5, 25},
			{4, 9, 49},
		}

		for _, tt := range tests {
			got := g.Index(tt.r, tt.c)
			if got != tt.want {
				t.Errorf("Index(%d, %d) = %d, want %d", tt.r, tt.c, got, tt.want)
			}
		}
	})

	t.Run("panics on negative r", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative r")
			}
		}()
		g := btmp.NewGridWithSize(5, 10)
		g.Index(-1, 0)
	})

	t.Run("panics on negative c", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative c")
			}
		}()
		g := btmp.NewGridWithSize(5, 10)
		g.Index(0, -1)
	})
}
