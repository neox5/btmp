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
	g.B.ensureBits(rows * g.cols)
}

// growCols increases Cols by delta and repositions existing rows.
// Internal implementation - no validation.
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
	g.B.ensureBits(rows * newCols)

	// Move each row to its new stride from bottom to top to avoid overlap issues
	if oldCols > 0 {
		for r := rows - 1; r >= 0; r-- {
			srcStart := r * oldCols
			dstStart := r * newCols

			// Copy old row data to new position
			g.B.copyRange(g.B, srcStart, dstStart, oldCols)

			// Clear the gap between old and new position if moving forward
			if srcStart < dstStart {
				g.B.clearRange(srcStart, dstStart-srcStart)
			}
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
	g.B.ensureBits(newRows * g.cols)
	// New bits are already zero from ensureBits
}
