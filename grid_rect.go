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
