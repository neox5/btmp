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

// nextZeroInRow returns the column index of the next zero bit in row r,
// starting search from column c.
// Returns -1 if no zero bit exists in [c, Cols()).
// Internal implementation - no validation.
func (g *Grid) nextZeroInRow(r, c int) int {
	start := g.rowStart(r) + c
	remaining := g.cols - c

	pos := g.B.NextZeroInRange(start, remaining)
	if pos == -1 {
		return -1
	}

	// Convert bitmap position back to column index
	return pos - g.rowStart(r)
}

// nextOneInRow returns the column index of the next set bit in row r,
// starting search from column c.
// Returns -1 if no set bit exists in [c, Cols()).
// Internal implementation - no validation.
func (g *Grid) nextOneInRow(r, c int) int {
	start := g.rowStart(r) + c
	remaining := g.cols - c

	pos := g.B.NextOneInRange(start, remaining)
	if pos == -1 {
		return -1
	}

	// Convert bitmap position back to column index
	return pos - g.rowStart(r)
}

// nextZeroInRowRange returns the column index of the next zero bit in row r,
// searching within [c, c+count).
// Returns -1 if no zero bit exists in range.
// Internal implementation - no validation.
func (g *Grid) nextZeroInRowRange(r, c, count int) int {
	start := g.rowStart(r) + c

	// Limit count to available columns
	available := g.cols - c
	searchCount := min(count, available)

	if searchCount <= 0 {
		return -1
	}

	pos := g.B.NextZeroInRange(start, searchCount)
	if pos == -1 {
		return -1
	}

	// Convert bitmap position back to column index
	return pos - g.rowStart(r)
}

// nextOneInRowRange returns the column index of the next set bit in row r,
// searching within [c, c+count).
// Returns -1 if no set bit exists in range.
// Internal implementation - no validation.
func (g *Grid) nextOneInRowRange(r, c, count int) int {
	start := g.rowStart(r) + c

	// Limit count to available columns
	available := g.cols - c
	searchCount := min(count, available)

	if searchCount <= 0 {
		return -1
	}

	pos := g.B.NextOneInRange(start, searchCount)
	if pos == -1 {
		return -1
	}

	// Convert bitmap position back to column index
	return pos - g.rowStart(r)
}

// countZerosFromInRow returns the count of consecutive zero bits in row r
// starting at column c.
// Returns 0 if bit at (r,c) is set.
// Stops at first set bit or end of row.
// Internal implementation - no validation.
func (g *Grid) countZerosFromInRow(r, c int) int {
	start := g.rowStart(r) + c
	remaining := g.cols - c

	return g.B.CountZerosFromInRange(start, remaining)
}

// countOnesFromInRow returns the count of consecutive set bits in row r
// starting at column c.
// Returns 0 if bit at (r,c) is zero.
// Stops at first zero bit or end of row.
// Internal implementation - no validation.
func (g *Grid) countOnesFromInRow(r, c int) int {
	start := g.rowStart(r) + c
	remaining := g.cols - c

	return g.B.CountOnesFromInRange(start, remaining)
}

// countZerosFromInRowRange returns the count of consecutive zero bits in row r
// starting at column c, within [c, c+count).
// Returns 0 if bit at (r,c) is set.
// Stops at first set bit or end of range.
// Internal implementation - no validation.
func (g *Grid) countZerosFromInRowRange(r, c, count int) int {
	start := g.rowStart(r) + c

	// Limit count to available columns
	available := g.cols - c
	searchCount := min(count, available)

	if searchCount <= 0 {
		return 0
	}

	return g.B.CountZerosFromInRange(start, searchCount)
}

// countOnesFromInRowRange returns the count of consecutive set bits in row r
// starting at column c, within [c, c+count).
// Returns 0 if bit at (r,c) is zero.
// Stops at first zero bit or end of range.
// Internal implementation - no validation.
func (g *Grid) countOnesFromInRowRange(r, c, count int) int {
	start := g.rowStart(r) + c

	// Limit count to available columns
	available := g.cols - c
	searchCount := min(count, available)

	if searchCount <= 0 {
		return 0
	}

	return g.B.CountOnesFromInRange(start, searchCount)
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
