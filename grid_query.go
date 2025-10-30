package btmp

// ========================================
// Internal Helpers
// ========================================

// rowStart returns the bitmap index for the first bit of row r.
// Internal helper - no validation.
func (g *Grid) rowStart(r int) int {
	return r * g.cols
}

// ========================================
// Query Implementations
// ========================================

// isFree reports whether the specified rectangle contains only zeros.
// Internal implementation - no validation, assumes valid bounds.
func (g *Grid) isFree(r, c, h, w int) bool {
	// Check each row of the rectangle
	for row := range h {
		start := (r+row)*g.cols + c
		// Check if any bit is set in this row segment
		for i := range w {
			if g.B.Test(start + i) {
				return false
			}
		}
	}
	return true
}

// canShiftRight reports whether rectangle can shift right.
// Checks if column c+w is free for rows [r, r+h).
// Internal implementation - no validation, assumes valid bounds.
func (g *Grid) canShiftRight(r, c, h, w int) bool {
	targetCol := c + w
	if targetCol >= g.cols {
		return false
	}
	return g.isFree(r, targetCol, h, 1)
}

// canShiftLeft reports whether rectangle can shift left.
// Checks if column c-1 is free for rows [r, r+h).
// Internal implementation - no validation, assumes valid bounds.
func (g *Grid) canShiftLeft(r, c, h, w int) bool {
	if c == 0 {
		return false
	}
	targetCol := c - 1
	return g.isFree(r, targetCol, h, 1)
}

// canShiftUp reports whether rectangle can shift up.
// Checks if row r-1 is free for columns [c, c+w).
// Internal implementation - no validation, assumes valid bounds.
func (g *Grid) canShiftUp(r, c, h, w int) bool {
	if r == 0 {
		return false
	}
	targetRow := r - 1
	return g.isFree(targetRow, c, 1, w)
}

// canShiftDown reports whether rectangle can shift down.
// Checks if row r+h is free for columns [c, c+w).
// Internal implementation - no validation, assumes valid bounds.
func (g *Grid) canShiftDown(r, c, h, w int) bool {
	targetRow := r + h
	if targetRow >= g.Rows() {
		return false
	}
	return g.isFree(targetRow, c, 1, w)
}

// nextFreeCol returns the column index of the next unoccupied cell in row r,
// starting search from column c.
// Returns -1 if no free column exists in [c, Cols()).
// Internal implementation - no validation.
func (g *Grid) nextFreeCol(r, c int) int {
	start := g.rowStart(r) + c
	remaining := g.cols - c

	pos := g.B.nextZeroInRange(start, remaining)
	if pos == -1 {
		return -1
	}

	// Convert bitmap position back to column index
	return pos - g.rowStart(r)
}

// nextFreeColInRange returns the column index of the next unoccupied cell in row r,
// searching within [c, c+count).
// Returns -1 if no free column exists in range.
// Internal implementation - no validation.
func (g *Grid) nextFreeColInRange(r, c, count int) int {
	start := g.rowStart(r) + c

	// Limit count to available columns
	available := g.cols - c
	searchCount := min(count, available)

	if searchCount <= 0 {
		return -1
	}

	pos := g.B.nextZeroInRange(start, searchCount)
	if pos == -1 {
		return -1
	}

	// Convert bitmap position back to column index
	return pos - g.rowStart(r)
}

// freeColsFrom returns the count of consecutive unoccupied columns in row r
// starting at column c.
// Returns 0 if cell at (r,c) is occupied.
// Stops at first occupied cell or end of row.
// Internal implementation - no validation.
func (g *Grid) freeColsFrom(r, c int) int {
	start := g.rowStart(r) + c
	remaining := g.cols - c

	return g.B.countZerosFromInRange(start, remaining)
}

// canFitWidth reports whether columns [c, c+w) in row r contain only unoccupied cells.
// Returns false if any cell in range is occupied or if c+w exceeds Cols().
// Internal implementation - no validation.
func (g *Grid) canFitWidth(r, c, w int) bool {
	// Check if width fits within row bounds
	if c+w > g.cols {
		return false
	}

	start := g.rowStart(r) + c

	// Check if all bits in range are zero (unoccupied)
	return !g.B.anyRange(start, w)
}

// canFit reports whether a rectangle of size h×w fits at position (r, c).
// Checks only boundary constraints, not cell occupancy.
// Returns two booleans:
//   - fitRow: true if r+h <= rows (height fits)
//   - fitCol: true if c+w <= cols (width fits)
//
// Internal implementation - no validation.
func (g *Grid) canFit(r, c, h, w int) (fitRow, fitCol bool) {
	return r+h <= g.rows, c+w <= g.cols
}

// allGrid returns true if all bits in the grid are set.
// Returns false for empty grid.
// Internal implementation - no validation.
func (g *Grid) allGrid() bool {
	// Empty grid has no bits to check
	if g.rows == 0 || g.cols == 0 {
		return false
	}

	// Delegate to bitmap's All() method
	return g.B.All()
}

// allRow returns true if all bits in row r are set.
// Returns false for empty row.
// Internal implementation - no validation.
func (g *Grid) allRow(r int) bool {
	// Empty row has no bits to check
	if g.cols == 0 {
		return false
	}

	start := g.rowStart(r)
	return g.B.AllRange(start, g.cols)
}
