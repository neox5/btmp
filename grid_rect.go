package btmp

// setRect sets rectangle to 1 without validation.
// Internal implementation - no auto-growth, requires in-bounds.
func (g *Grid) setRect(c, r, w, h int) {
	if w == 0 || h == 0 {
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
func (g *Grid) clearRect(c, r, w, h int) {
	if w == 0 || h == 0 {
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
func (g *Grid) isFree(c, r, w, h int) bool {
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
func (g *Grid) canShiftRight(c, r, w, h int) bool {
	targetCol := c + w
	if targetCol >= g.cols {
		return false
	}
	return g.isFree(targetCol, r, 1, h)
}

// canShiftLeft reports whether rectangle can shift left.
// Checks if column c-1 is free for rows [r, r+h).
// Internal helper - no validation, assumes valid bounds.
func (g *Grid) canShiftLeft(c, r, w, h int) bool {
	if c == 0 {
		return false
	}
	targetCol := c - 1
	return g.isFree(targetCol, r, 1, h)
}

// canShiftUp reports whether rectangle can shift up.
// Checks if row r-1 is free for columns [c, c+w).
// Internal helper - no validation, assumes valid bounds.
func (g *Grid) canShiftUp(c, r, w, h int) bool {
	if r == 0 {
		return false
	}
	targetRow := r - 1
	return g.isFree(c, targetRow, w, 1)
}

// canShiftDown reports whether rectangle can shift down.
// Checks if row r+h is free for columns [c, c+w).
// Internal helper - no validation, assumes valid bounds.
func (g *Grid) canShiftDown(c, r, w, h int) bool {
	targetRow := r + h
	if targetRow >= g.Rows() {
		return false
	}
	return g.isFree(c, targetRow, w, 1)
}

// shiftRectRight shifts a rectangle one column to the right.
// Moves bits from [c,r,w,h) to [c+1,r,w,h).
// The leftmost column (c) is cleared.
// Internal implementation - no validation, requires in-bounds and target column free.
func (g *Grid) shiftRectRight(c, r, w, h int) {
	if w == 0 || h == 0 {
		return
	}

	for row := range h {
		srcStart := (r+row)*g.cols + c
		dstStart := srcStart + 1
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectLeft shifts a rectangle one column to the left.
// Moves bits from [c,r,w,h) to [c-1,r,w,h).
// The rightmost column (c+w-1) is cleared.
// Internal implementation - no validation, requires in-bounds and target column free.
func (g *Grid) shiftRectLeft(c, r, w, h int) {
	if w == 0 || h == 0 {
		return
	}

	for row := range h {
		srcStart := (r+row)*g.cols + c
		dstStart := srcStart - 1
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectUp shifts a rectangle one row up.
// Moves bits from [c,r,w,h) to [c,r-1,w,h).
// The bottom row (r+h-1) is cleared.
// Internal implementation - no validation, requires in-bounds and target row free.
func (g *Grid) shiftRectUp(c, r, w, h int) {
	if w == 0 || h == 0 {
		return
	}

	for row := range h {
		srcStart := (r+row)*g.cols + c
		dstStart := ((r-1)+row)*g.cols + c
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectDown shifts a rectangle one row down.
// Moves bits from [c,r,w,h) to [c,r+1,w,h).
// The top row (r) is cleared.
// Internal implementation - no validation, requires in-bounds and target row free.
func (g *Grid) shiftRectDown(c, r, w, h int) {
	if w == 0 || h == 0 {
		return
	}

	// Process rows in reverse to avoid overlap issues
	for row := h - 1; row >= 0; row-- {
		srcStart := (r+row)*g.cols + c
		dstStart := ((r+1)+row)*g.cols + c
		g.B.MoveRange(srcStart, dstStart, w)
	}
}
