package btmp

// validateCoordinate validates that x and y are non-negative and within grid bounds.
// Panics if x < 0, y < 0, x >= g.Cols(), or y >= g.Rows().
func (g *Grid) validateCoordinate(x, y int) {
	validateNonNegative(x, "x")
	validateNonNegative(y, "y")
	if x >= g.cols {
		panic("x out of bounds")
	}
	if y >= g.Rows() {
		panic("y out of bounds")
	}
}

// validateRect validates that rectangle parameters are non-negative
// and rectangle is fully contained within grid bounds.
// Panics if x < 0, y < 0, w < 0, h < 0, x+w > g.Cols(), or y+h > g.Rows().
func (g *Grid) validateRect(x, y, w, h int) {
	g.validateCoordinate(x, y)
	validateNonNegative(w, "w")
	validateNonNegative(h, "h")
	if x+w > g.cols {
		panic("rectangle exceeds columns")
	}
	if y+h > g.Rows() {
		panic("rectangle exceeds rows")
	}
}
