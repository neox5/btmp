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

// rectZero reports whether the specified rectangle contains only zeros.
// Internal implementation - no validation, assumes valid bounds.
func (g *Grid) rectZero(r, c, h, w int) bool {
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

// rectOne reports whether the specified rectangle contains only ones.
// Internal implementation - no validation, assumes valid bounds.
func (g *Grid) rectOne(r, c, h, w int) bool {
	// Check each row of the rectangle
	for row := range h {
		start := (r+row)*g.cols + c
		// Check if any bit is zero in this row segment
		for i := range w {
			if !g.B.Test(start + i) {
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
