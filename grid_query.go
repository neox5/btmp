package btmp

// IsFree reports whether the specified rectangle contains only zeros.
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) IsFree(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)

	if w == 0 || h == 0 {
		return true
	}

	// Check each row of the rectangle
	for row := range h {
		start := (y+row)*g.cols + x
		// Check if any bit is set in this row segment
		for i := range w {
			if g.B.Test(start + i) {
				return false
			}
		}
	}

	return true
}

// CanShiftRight reports whether the rectangle can shift one column right.
// Checks if column x+w exists and is free (all zeros) for rows [y, y+h).
// Panics if rectangle is invalid, out of bounds, or target column doesn't exist.
func (g *Grid) CanShiftRight(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)

	targetCol := x + w
	if targetCol >= g.cols {
		panic("target column out of bounds")
	}

	if h == 0 {
		return true
	}

	// Check if target column is free for all rows in range
	for row := range h {
		pos := (y+row)*g.cols + targetCol
		if g.B.Test(pos) {
			return false
		}
	}

	return true
}

// CanShiftLeft reports whether the rectangle can shift one column left.
// Checks if column x-1 exists and is free (all zeros) for rows [y, y+h).
// Panics if rectangle is invalid, out of bounds, or target column doesn't exist.
func (g *Grid) CanShiftLeft(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)

	if x == 0 {
		panic("target column out of bounds")
	}

	targetCol := x - 1

	if h == 0 {
		return true
	}

	// Check if target column is free for all rows in range
	for row := range h {
		pos := (y+row)*g.cols + targetCol
		if g.B.Test(pos) {
			return false
		}
	}

	return true
}

// CanShiftUp reports whether the rectangle can shift one row up.
// Checks if row y-1 exists and is free (all zeros) for columns [x, x+w).
// Panics if rectangle is invalid, out of bounds, or target row doesn't exist.
func (g *Grid) CanShiftUp(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)

	if y == 0 {
		panic("target row out of bounds")
	}

	targetRow := y - 1

	if w == 0 {
		return true
	}

	// Check if target row is free for all columns in range
	start := targetRow*g.cols + x
	for col := range w {
		if g.B.Test(start + col) {
			return false
		}
	}

	return true
}

// CanShiftDown reports whether the rectangle can shift one row down.
// Checks if row y+h exists and is free (all zeros) for columns [x, x+w).
// Panics if rectangle is invalid, out of bounds, or target row doesn't exist.
func (g *Grid) CanShiftDown(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)

	targetRow := y + h
	if targetRow >= g.Rows() {
		panic("target row out of bounds")
	}

	if w == 0 {
		return true
	}

	// Check if target row is free for all columns in range
	start := targetRow*g.cols + x
	for col := range w {
		if g.B.Test(start + col) {
			return false
		}
	}

	return true
}
