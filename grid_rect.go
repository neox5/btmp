package btmp

// setRect sets rectangle to 1 without validation.
// Internal implementation - no auto-growth, requires in-bounds.
func (g *Grid) setRect(x, y, w, h int) {
	if w == 0 || h == 0 {
		// Empty rectangle, nothing to do
		return
	}

	// Set each row of the rectangle
	for row := range h {
		start := (y+row)*g.cols + x
		g.B.setRange(start, w)
	}
}

// clearRect clears rectangle to 0 without validation.
// Internal implementation - no auto-growth.
func (g *Grid) clearRect(x, y, w, h int) {
	if w == 0 || h == 0 {
		// Empty rectangle, nothing to do
		return
	}

	// Clear each row of the rectangle
	for row := range h {
		start := (y+row)*g.cols + x
		g.B.clearRange(start, w)
	}
}

// shiftRectRight shifts a rectangle one column to the right.
// Moves bits from [x,y,w,h) to [x+1,y,w,h).
// The leftmost column (x) is cleared.
// Internal implementation - no validation, requires in-bounds and target column free.
func (g *Grid) shiftRectRight(x, y, w, h int) {
	if w == 0 || h == 0 {
		return
	}

	for row := range h {
		srcStart := (y+row)*g.cols + x
		dstStart := srcStart + 1
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectLeft shifts a rectangle one column to the left.
// Moves bits from [x,y,w,h) to [x-1,y,w,h).
// The rightmost column (x+w-1) is cleared.
// Internal implementation - no validation, requires in-bounds and target column free.
func (g *Grid) shiftRectLeft(x, y, w, h int) {
	if w == 0 || h == 0 {
		return
	}

	for row := range h {
		srcStart := (y+row)*g.cols + x
		dstStart := srcStart - 1
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectUp shifts a rectangle one row up.
// Moves bits from [x,y,w,h) to [x,y-1,w,h).
// The bottom row (y+h-1) is cleared.
// Internal implementation - no validation, requires in-bounds and target row free.
func (g *Grid) shiftRectUp(x, y, w, h int) {
	if w == 0 || h == 0 {
		return
	}

	for row := range h {
		srcStart := (y+row)*g.cols + x
		dstStart := ((y-1)+row)*g.cols + x
		g.B.MoveRange(srcStart, dstStart, w)
	}
}

// shiftRectDown shifts a rectangle one row down.
// Moves bits from [x,y,w,h) to [x,y+1,w,h).
// The top row (y) is cleared.
// Internal implementation - no validation, requires in-bounds and target row free.
func (g *Grid) shiftRectDown(x, y, w, h int) {
	if w == 0 || h == 0 {
		return
	}

	// Process rows in reverse to avoid overlap issues
	for row := h - 1; row >= 0; row-- {
		srcStart := (y+row)*g.cols + x
		dstStart := ((y+1)+row)*g.cols + x
		g.B.MoveRange(srcStart, dstStart, w)
	}
}
