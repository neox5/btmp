package btmp

// setRect sets rectangle to 1 without validation.
// Internal implementation - no auto-growth, requires in-bounds.
func (g *Grid) setRect(r, c, h, w int) {
	if h == 0 || w == 0 {
		// Empty rectangle, nothing to do
		return
	}

	// Set each row of the rectangle
	for row := range h {
		start := (r+row)*g.cols + c
		g.B.setRange(start, w)
	}
}

// clearRect clears rectangle to 0 without validation.
// Internal implementation - no auto-growth.
func (g *Grid) clearRect(r, c, h, w int) {
	if h == 0 || w == 0 {
		// Empty rectangle, nothing to do
		return
	}

	// Clear each row of the rectangle
	for row := range h {
		start := (r+row)*g.cols + c
		g.B.clearRange(start, w)
	}
}

// isFree reports whether the specified rectangle contains only zeros.
// Internal helper - no validation, assumes valid bounds.
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
// Internal helper - no validation, assumes valid bounds.
func (g *Grid) canShiftRight(r, c, h, w int) bool {
	targetCol := c + w
	if targetCol >= g.cols {
		return false
	}
	return g.isFree(r, targetCol, h, 1)
}

// canShiftLeft reports whether rectangle can shift left.
// Checks if column c-1 is free for rows [r, r+h).
// Internal helper - no validation, assumes valid bounds.
func (g *Grid) canShiftLeft(r, c, h, w int) bool {
	if c == 0 {
		return false
	}
	targetCol := c - 1
	return g.isFree(r, targetCol, h, 1)
}

// canShiftUp reports whether rectangle can shift up.
// Checks if row r-1 is free for columns [c, c+w).
// Internal helper - no validation, assumes valid bounds.
func (g *Grid) canShiftUp(r, c, h, w int) bool {
	if r == 0 {
		return false
	}
	targetRow := r - 1
	return g.isFree(targetRow, c, 1, w)
}

// canShiftDown reports whether rectangle can shift down.
// Checks if row r+h is free for columns [c, c+w).
// Internal helper - no validation, assumes valid bounds.
func (g *Grid) canShiftDown(r, c, h, w int) bool {
	targetRow := r + h
	if targetRow >= g.Rows() {
		return false
	}
	return g.isFree(targetRow, c, 1, w)
}

// shiftRectRight shifts a rectangle one column to the right.
// Moves bits from [r,c,h,w) to [r,c+1,h,w).
// The leftmost column (c) is cleared.
// Internal implementation - no validation, requires in-bounds and target column free.
func (g *Grid) shiftRectRight(r, c, h, w int) {
	if h == 0 || w == 0 {
		return
	}

	for row := range h {
		srcStart := (r+row)*g.cols + c
		dstStart := srcStart + 1
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectLeft shifts a rectangle one column to the left.
// Moves bits from [r,c,h,w) to [r,c-1,h,w).
// The rightmost column (c+w-1) is cleared.
// Internal implementation - no validation, requires in-bounds and target column free.
func (g *Grid) shiftRectLeft(r, c, h, w int) {
	if h == 0 || w == 0 {
		return
	}

	for row := range h {
		srcStart := (r+row)*g.cols + c
		dstStart := srcStart - 1
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectUp shifts a rectangle one row up.
// Moves bits from [r,c,h,w) to [r-1,c,h,w).
// The bottom row (r+h-1) is cleared.
// Internal implementation - no validation, requires in-bounds and target row free.
func (g *Grid) shiftRectUp(r, c, h, w int) {
	if h == 0 || w == 0 {
		return
	}

	for row := range h {
		srcStart := (r+row)*g.cols + c
		dstStart := ((r-1)+row)*g.cols + c
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectDown shifts a rectangle one row down.
// Moves bits from [r,c,h,w) to [r+1,c,h,w).
// The top row (r) is cleared.
// Internal implementation - no validation, requires in-bounds and target row free.
func (g *Grid) shiftRectDown(r, c, h, w int) {
	if h == 0 || w == 0 {
		return
	}

	// Process rows in reverse to avoid overlap issues
	for row := h - 1; row >= 0; row-- {
		srcStart := (r+row)*g.cols + c
		dstStart := ((r+1)+row)*g.cols + c
		g.B.MoveRange(srcStart, dstStart, w)
	}
}
