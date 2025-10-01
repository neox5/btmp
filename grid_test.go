package btmp_test

import (
	"testing"

	"github.com/neox5/btmp"
)

// TestNewGridWithSize validates NewGridWithSize constructor behavior.
// Tests:
//   - Creates grid with exact rows×cols bits
//   - Zero rows/cols accepted
//   - Negative values panic
//   - Overflow detection (rows * cols < 0)
//   - Bitmap length matches rows×cols
func TestNewGridWithSize(t *testing.T) {
	t.Run("creates grid with specified dimensions", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 5)
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

	t.Run("accepts zero columns", func(t *testing.T) {
		g := btmp.NewGridWithSize(0, 5)
		if g.Cols() != 0 {
			t.Errorf("expected cols=0, got %d", g.Cols())
		}
		if g.B.Len() != 0 {
			t.Errorf("expected bitmap len=0, got %d", g.B.Len())
		}
	})

	t.Run("accepts zero rows", func(t *testing.T) {
		g := btmp.NewGridWithSize(10, 0)
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

	t.Run("accepts both dimensions zero", func(t *testing.T) {
		g := btmp.NewGridWithSize(0, 0)
		if g.Cols() != 0 {
			t.Errorf("expected cols=0, got %d", g.Cols())
		}
		if g.B.Len() != 0 {
			t.Errorf("expected bitmap len=0, got %d", g.B.Len())
		}
	})

	t.Run("panics on negative columns", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative columns")
			}
		}()
		btmp.NewGridWithSize(-1, 5)
	})

	t.Run("panics on negative rows", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for negative rows")
			}
		}()
		btmp.NewGridWithSize(10, -1)
	})

	t.Run("panics on overflow", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic for overflow")
			}
		}()
		// MaxInt rows × 2 cols should overflow
		btmp.NewGridWithSize(2, int(^uint(0)>>1))
	})

	t.Run("bitmap initialized to zeros", func(t *testing.T) {
		g := btmp.NewGridWithSize(8, 4)
		if g.B.Any() {
			t.Error("expected all bits to be zero")
		}
		if g.B.Count() != 0 {
			t.Errorf("expected count=0, got %d", g.B.Count())
		}
	})
}
