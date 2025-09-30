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
