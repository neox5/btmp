package btmp_test

import (
	"math"
	"testing"

	"github.com/neox5/btmp"
)

// bitAt reads bit i via public API.
func bitAt(b *btmp.Bitmap, i int) bool { return b.Test(i) }

func TestGrid_ConstructorsAndRows(t *testing.T) {
	t.Parallel()

	g := btmp.NewGrid(16)
	if g.Cols() != 16 {
		t.Fatalf("Cols=%d want=16", g.Cols())
	}
	if g.Rows() != 0 || g.B.Len() != 0 {
		t.Fatalf("new grid not empty")
	}

	g = btmp.NewGridWithCap(8, 10)
	if g.Cols() != 8 {
		t.Fatalf("Cols=%d want=8", g.Cols())
	}
	if g.Rows() != 0 || g.B.Len() != 0 {
		t.Fatalf("cap ctor must not change length")
	}

	g = btmp.NewGridWithSize(5, 7)
	if g.Cols() != 5 || g.Rows() != 7 || g.B.Len() != 35 {
		t.Fatalf("size ctor mismatch: Cols=%d Rows=%d Len=%d", g.Cols(), g.Rows(), g.B.Len())
	}

	base := btmp.New(24)
	g = btmp.NewGridFrom(base, 6)
	if g.Cols() != 6 || g.Rows() != 4 {
		t.Fatalf("from mismatch: Cols=%d Rows=%d", g.Cols(), g.Rows())
	}
}

func TestGrid_IndexAndRect(t *testing.T) {
	t.Parallel()

	g := btmp.NewGrid(16)
	// Fill 5x4 at (3,2) → rows grow to at least y+h = 6
	g = g.SetRect(3, 2, 5, 4)

	if g.Rows() != 6 {
		t.Fatalf("Rows=%d want=6", g.Rows())
	}
	// Check rectangle bits set, others zero within covered rows.
	for y := range g.Rows() {
		for x := range g.Cols() {
			i := g.Index(x, y)
			inRect := (x >= 3 && x < 3+5) && (y >= 2 && y < 2+4)
			if bitAt(g.B, i) != inRect {
				t.Fatalf("bit(%d,%d)=%v want %v", x, y, bitAt(g.B, i), inRect)
			}
		}
	}

	// Clear a 2x2 sub-rect at (4,3)
	g = g.ClearRect(4, 3, 2, 2)
	for y := 3; y < 5; y++ {
		for x := 4; x < 6; x++ {
			if bitAt(g.B, g.Index(x, y)) {
				t.Fatalf("expected cleared at (%d,%d)", x, y)
			}
		}
	}
}

func TestGrid_GrowCols_PreservesCoordinates(t *testing.T) {
	t.Parallel()

	g := btmp.NewGrid(8)
	// Mark a few points
	points := [][2]int{{1, 1}, {3, 2}, {7, 0}, {0, 3}}
	for _, p := range points {
		g = g.SetRect(p[0], p[1], 1, 1)
	}
	rowsBefore := g.Rows()

	// Grow columns by 5 → Cols=13
	g = g.GrowCols(5)
	if g.Cols() != 13 {
		t.Fatalf("Cols=%d want=13", g.Cols())
	}
	if g.Rows() != rowsBefore {
		t.Fatalf("Rows changed on GrowCols: %d -> %d", rowsBefore, g.Rows())
	}

	// All previous (x,y) must still be set at same coordinates.
	for _, p := range points {
		if !bitAt(g.B, g.Index(p[0], p[1])) {
			t.Fatalf("lost bit at (%d,%d) after GrowCols", p[0], p[1])
		}
	}
	// Newly added columns must be zero.
	for y := range g.Rows() {
		for x := 8; x < 13; x++ {
			if bitAt(g.B, g.Index(x, y)) {
				t.Fatalf("new col not zero at (%d,%d)", x, y)
			}
		}
	}
}

func TestGrid_GrowEnsureRows(t *testing.T) {
	t.Parallel()

	g := btmp.NewGrid(10)
	g = g.GrowRows(3)
	if g.Rows() != 3 {
		t.Fatalf("Rows=%d want=3", g.Rows())
	}
	// Write on last row
	g = g.SetRect(0, 2, 10, 1)

	// EnsureRows no-op then growth
	g = g.EnsureRows(2)
	if g.Rows() != 3 {
		t.Fatalf("EnsureRows shrink should be no-op")
	}
	g = g.EnsureRows(8)
	if g.Rows() != 8 {
		t.Fatalf("Rows=%d want=8", g.Rows())
	}
	// Prior data intact
	for x := range 10 {
		if !bitAt(g.B, g.Index(x, 2)) {
			t.Fatalf("lost bit at row 2 col %d after EnsureRows", x)
		}
	}
	// New rows zero
	for y := 3; y < 8; y++ {
		for x := range 10 {
			if bitAt(g.B, g.Index(x, y)) {
				t.Fatalf("expected zero in new rows at (%d,%d)", x, y)
			}
		}
	}
}

func TestGrid_ZeroColsEdgeCases(t *testing.T) {
	t.Parallel()

	// Zero cols with GrowRows
	g := btmp.NewGrid(0)
	g = g.GrowRows(5) // Should be no-op
	if g.Rows() != 0 || g.B.Len() != 0 {
		t.Fatalf("GrowRows with zero cols should be no-op")
	}

	// Zero cols with EnsureRows
	g = g.EnsureRows(10) // Should be no-op
	if g.Rows() != 0 || g.B.Len() != 0 {
		t.Fatalf("EnsureRows with zero cols should be no-op")
	}
}

func TestSetRect_AutoGrow(t *testing.T) {
	t.Parallel()

	// Test auto-grow when rectangle exceeds current rows
	g := btmp.NewGrid(10)
	// Initially 0 rows
	if g.Rows() != 0 {
		t.Fatalf("new grid should have 0 rows, got %d", g.Rows())
	}

	// SetRect that requires growing to y+h rows
	g = g.SetRect(2, 5, 3, 4) // Rectangle at (2,5) with size 3x4
	expectedRows := 5 + 4     // y + h = 9
	if g.Rows() != expectedRows {
		t.Fatalf("SetRect should auto-grow to %d rows, got %d", expectedRows, g.Rows())
	}

	// Verify the rectangle was set correctly
	for y := 5; y < 9; y++ {
		for x := 2; x < 5; x++ {
			if !bitAt(g.B, g.Index(x, y)) {
				t.Fatalf("bit should be set at (%d,%d)", x, y)
			}
		}
	}

	// Test growing from existing rows
	g = g.SetRect(0, 10, 5, 2) // Extend beyond current 9 rows
	expectedRows = 10 + 2      // y + h = 12
	if g.Rows() != expectedRows {
		t.Fatalf("SetRect should auto-grow to %d rows, got %d", expectedRows, g.Rows())
	}
}

func TestGrowRows_ZeroCols(t *testing.T) {
	t.Parallel()

	// Test GrowRows with zero columns
	g := btmp.NewGrid(0) // Zero columns
	if g.Cols() != 0 {
		t.Fatalf("expected 0 cols, got %d", g.Cols())
	}
	if g.Rows() != 0 {
		t.Fatalf("expected 0 rows, got %d", g.Rows())
	}

	// GrowRows should be no-op when cols == 0
	g = g.GrowRows(5)
	if g.Rows() != 0 {
		t.Fatalf("GrowRows with zero cols should be no-op, got %d rows", g.Rows())
	}
	if g.B.Len() != 0 {
		t.Fatalf("bitmap length should remain 0, got %d", g.B.Len())
	}

	// Multiple GrowRows calls
	g = g.GrowRows(10)
	if g.Rows() != 0 {
		t.Fatalf("multiple GrowRows with zero cols should be no-op, got %d rows", g.Rows())
	}
}

func TestEnsureRows_ZeroCols(t *testing.T) {
	t.Parallel()

	// Test EnsureRows with zero columns
	g := btmp.NewGrid(0) // Zero columns
	if g.Cols() != 0 {
		t.Fatalf("expected 0 cols, got %d", g.Cols())
	}

	// EnsureRows should be no-op when cols == 0
	g = g.EnsureRows(5)
	if g.Rows() != 0 {
		t.Fatalf("EnsureRows with zero cols should be no-op, got %d rows", g.Rows())
	}
	if g.B.Len() != 0 {
		t.Fatalf("bitmap length should remain 0, got %d", g.B.Len())
	}

	// EnsureRows with larger value
	g = g.EnsureRows(100)
	if g.Rows() != 0 {
		t.Fatalf("EnsureRows with zero cols should be no-op, got %d rows", g.Rows())
	}
	if g.B.Len() != 0 {
		t.Fatalf("bitmap length should remain 0, got %d", g.B.Len())
	}

	// EnsureRows no-op (requesting fewer rows than current)
	g = g.EnsureRows(0)
	if g.Rows() != 0 {
		t.Fatalf("EnsureRows(0) should be no-op, got %d rows", g.Rows())
	}
}

func TestGrid_ZeroColsComprehensive(t *testing.T) {
	t.Parallel()

	// Comprehensive test of zero-column grid behavior
	g := btmp.NewGrid(0)

	// All row operations should be no-ops
	g = g.GrowRows(1)
	g = g.EnsureRows(10)
	g = g.GrowRows(5)

	if g.Rows() != 0 {
		t.Fatalf("all row operations on zero-col grid should be no-ops, got %d rows", g.Rows())
	}
	if g.Cols() != 0 {
		t.Fatalf("columns should remain 0, got %d", g.Cols())
	}
	if g.B.Len() != 0 {
		t.Fatalf("bitmap length should remain 0, got %d", g.B.Len())
	}

	// Verify grid remains in consistent state
	if g.Rows() != 0 || g.Cols() != 0 {
		t.Fatal("zero-col grid should maintain zero rows and cols")
	}
}

func TestGrowCols_EmptyGrid(t *testing.T) {
	t.Parallel()

	// Test GrowCols when rows == 0 (covers line 140)
	g := btmp.NewGrid(5)
	// No rows initially
	if g.Rows() != 0 {
		t.Fatalf("expected 0 rows, got %d", g.Rows())
	}

	// GrowCols should just update cols when rows == 0
	g = g.GrowCols(3)
	if g.Cols() != 8 {
		t.Fatalf("expected 8 cols after grow, got %d", g.Cols())
	}
	if g.Rows() != 0 {
		t.Fatalf("expected 0 rows still, got %d", g.Rows())
	}
}

func TestEnsureCols(t *testing.T) {
	t.Parallel()

	// Test panic on negative cols (covers line 169)
	g := btmp.NewGrid(10)
	mustPanic(t, func() { _ = g.EnsureCols(-1) })

	// Test no-op when cols <= current (covers lines 171-173)
	g = g.SetRect(0, 0, 10, 5) // Create 5 rows
	originalBits := g.B.Count()

	// EnsureCols with same amount - should be no-op
	g = g.EnsureCols(10)
	if g.Cols() != 10 {
		t.Fatalf("EnsureCols(10) should keep 10 cols, got %d", g.Cols())
	}
	if g.B.Count() != originalBits {
		t.Fatal("EnsureCols no-op should not change bits")
	}

	// EnsureCols with less - should be no-op
	g = g.EnsureCols(5)
	if g.Cols() != 10 {
		t.Fatalf("EnsureCols(5) should keep 10 cols, got %d", g.Cols())
	}

	// EnsureCols with more - should call GrowCols (covers line 174)
	g = g.EnsureCols(15)
	if g.Cols() != 15 {
		t.Fatalf("EnsureCols(15) should result in 15 cols, got %d", g.Cols())
	}

	// Verify data preserved after growth
	for y := range 5 {
		for x := range 10 {
			if !g.B.Test(g.Index(x, y)) {
				t.Fatalf("lost bit at (%d,%d) after EnsureCols", x, y)
			}
		}
	}
}

func TestGrid_Panics(t *testing.T) {
	t.Parallel()

	// Constructor panics
	mustPanic(t, func() { _ = btmp.NewGrid(-1) })
	mustPanic(t, func() { _ = btmp.NewGridWithCap(-1, 0) })
	mustPanic(t, func() { _ = btmp.NewGridWithCap(0, -1) })
	mustPanic(t, func() { _ = btmp.NewGridWithCap(math.MaxInt/2, math.MaxInt/2+1) }) // overflow
	mustPanic(t, func() { _ = btmp.NewGridWithSize(-1, 0) })
	mustPanic(t, func() { _ = btmp.NewGridWithSize(0, -1) })
	mustPanic(t, func() { _ = btmp.NewGridWithSize(math.MaxInt/2, math.MaxInt/2+1) }) // overflow
	mustPanic(t, func() { _ = btmp.NewGridFrom(nil, 8) })
	mustPanic(t, func() { _ = btmp.NewGridFrom(btmp.New(0), -1) })

	g := btmp.NewGrid(8)

	// Index panics
	mustPanic(t, func() { _ = g.Index(-1, 0) })
	mustPanic(t, func() { _ = g.Index(0, -1) })

	// SetRect panics
	mustPanic(t, func() { _ = g.SetRect(-1, 0, 1, 1) })
	mustPanic(t, func() { _ = g.SetRect(0, -1, 1, 1) })
	mustPanic(t, func() { _ = g.SetRect(0, 0, -1, 1) })
	mustPanic(t, func() { _ = g.SetRect(0, 0, 1, -1) })
	mustPanic(t, func() { _ = g.SetRect(7, 0, 2, 1) }) // x+w > cols

	// SetRect with zero cols
	gZero := btmp.NewGrid(0)
	mustPanic(t, func() { _ = gZero.SetRect(0, 0, 1, 1) }) // cols == 0 and w > 0

	// ClearRect panics
	mustPanic(t, func() { _ = g.ClearRect(-1, 0, 1, 1) })
	mustPanic(t, func() { _ = g.ClearRect(0, -1, 1, 1) })
	mustPanic(t, func() { _ = g.ClearRect(0, 0, -1, 1) })
	mustPanic(t, func() { _ = g.ClearRect(0, 0, 1, -1) })
	mustPanic(t, func() { _ = g.ClearRect(7, 0, 2, 1) }) // x+w > cols

	// ClearRect beyond current rows
	g = g.GrowRows(2)
	mustPanic(t, func() { _ = g.ClearRect(0, 1, 8, 2) }) // y+h > rows

	// GrowCols/EnsureCols panics
	mustPanic(t, func() { _ = g.GrowCols(0) })
	mustPanic(t, func() { _ = g.GrowCols(-1) })
	mustPanic(t, func() { _ = g.EnsureCols(-1) })

	// GrowRows/EnsureRows panics
	mustPanic(t, func() { _ = g.GrowRows(0) })
	mustPanic(t, func() { _ = g.GrowRows(-1) })
	mustPanic(t, func() { _ = g.EnsureRows(-1) })
}
