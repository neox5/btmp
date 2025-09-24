package btmp_test

import (
	"testing"

	btmp "github.com/neox5/btmp"
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
	if g.Cols()!= 8 {
		t.Fatalf("Cols=%d want=8", g.Cols())
	}
	if g.Rows() != 0 || g.B.Len() != 0 {
		t.Fatalf("cap ctor must not change length")
	}

	g = btmp.NewGridWithSize(5, 7)
	if g.Cols() != 5 || g.Rows() != 7 || g.B.Len() != 35 {
		t.Fatalf("size ctor mismatch: Cols=%d Rows=%d Len=%d", g.Cols(), g.Rows(), g.B.Len())
	}

	base := btmp.New().EnsureBits(24)
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

func TestGrid_Panics(t *testing.T) {
	t.Parallel()

	// NewGridFrom nil
	mustPanic(t, func() { _ = btmp.NewGridFrom(nil, 8) })

	g := btmp.NewGrid(8)

	// Negative args
	mustPanic(t, func() { _ = g.SetRect(-1, 0, 1, 1) })
	mustPanic(t, func() { _ = g.SetRect(0, -1, 1, 1) })
	mustPanic(t, func() { _ = g.ClearRect(-1, 0, 1, 1) })

	// Exceed cols
	mustPanic(t, func() { _ = g.SetRect(7, 0, 2, 1) })
	mustPanic(t, func() { _ = g.ClearRect(7, 0, 2, 1) })

	// Clear beyond rows
	g = g.GrowRows(2)
	mustPanic(t, func() { _ = g.ClearRect(0, 1, 8, 2) })

	// GrowCols invalid
	mustPanic(t, func() { _ = g.GrowCols(0) })

	// EnsureCols negative
	mustPanic(t, func() { _ = g.EnsureCols(-1) })

	// GrowRows invalid
	mustPanic(t, func() { _ = g.GrowRows(0) })

	// EnsureRows negative
	mustPanic(t, func() { _ = g.EnsureRows(-1) })
}
