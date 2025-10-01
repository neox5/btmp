package btmp

// ensureCols grows Cols to at least cols without validation.
// Internal implementation - no bounds checking.
func (g *Grid) ensureCols(cols int) {
	if cols <= g.cols {
		return
	}
	g.growCols(cols - g.cols)
}

// ensureRows ensures at least rows rows exist without validation.
// Internal implementation - no bounds checking.
func (g *Grid) ensureRows(rows int) {
	if g.cols == 0 {
		// Cannot have rows without columns
		return
	}
	if rows <= g.Rows() {
		return
	}
	g.B.EnsureBits(rows * g.cols)
}

// growCols increases Cols by delta and repositions existing rows.
// Internal implementation - no validation.
//
// Algorithm:
//  1. Resize backing store to accommodate new layout (rows * newCols)
//  2. Move each row from old stride to new stride using MoveRange (bottom-to-top)
//  3. MoveRange handles both copying and clearing source, maintaining bitmap invariants
//
// Example: 3x2 grid growing to 5x2:
//  Before: [0,1,2][3,4,5] (cols=3, len=6)
//  After:  [0,1,2,_,_][3,4,5,_,_] (cols=5, len=10)
//
// Row repositioning (backward iteration):
//  Row 1: MoveRange([3,4,5] → [5,6,7]) - moves data and clears non-overlapping source
//  Row 0: MoveRange([0,1,2] → [0,1,2]) - no-op, srcStart == dstStart
//
// MoveRange automatically clears the gap [3,5) after moving Row 1,
// removing old inter-row bits that are now invalid in the new layout.
func (g *Grid) growCols(delta int) {
	oldCols := g.cols
	newCols := oldCols + delta
	rows := g.Rows()

	if rows == 0 {
		// No data to reposition
		g.cols = newCols
		return
	}

	// Resize backing store to new size
	g.B.EnsureBits(rows * newCols)

	// Move each row to its new stride from bottom to top to avoid overlap issues
	if oldCols > 0 {
		for r := rows - 1; r >= 0; r-- {
			srcStart := r * oldCols
			dstStart := r * newCols

			// MoveRange handles copy + clear of non-overlapping source region
			g.B.MoveRange(srcStart, dstStart, oldCols)
		}
	}

	g.cols = newCols
}

// growRows appends delta empty rows without validation.
// Internal implementation - no validation.
func (g *Grid) growRows(delta int) {
	if g.cols == 0 {
		// Cannot add rows without columns
		return
	}
	newRows := g.Rows() + delta
	g.B.EnsureBits(newRows * g.cols)
	// New bits are already zero from EnsureBits
}
