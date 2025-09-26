package btmp

// validateCols validates that cols is non-negative.
// Panics if cols < 0.
func validateCols(cols int) {
	if cols < 0 {
		panic("negative cols")
	}
}

// validateRows validates that rows is non-negative.
// Panics if rows < 0.
func validateRows(rows int) {
	if rows < 0 {
		panic("negative rows")
	}
}

// validateDelta validates that delta is positive.
// Panics if delta <= 0.
func validateDelta(delta int) {
	if delta <= 0 {
		panic("delta must be > 0")
	}
}

// validateCoordinate validates that x and y are non-negative.
// Panics if x < 0 or y < 0.
func validateCoordinate(x, y int) {
	if x < 0 || y < 0 {
		panic("negative coordinate")
	}
}

// validateRectParams validates rectangle parameters are non-negative.
// Panics if x < 0, y < 0, w < 0, or h < 0.
func validateRectParams(x, y, w, h int) {
	if x < 0 || y < 0 || w < 0 || h < 0 {
		panic("negative rectangle parameter")
	}
}

// validateGridSize validates that cols*rows doesn't overflow.
// Panics if cols*rows < 0.
func validateGridSize(cols, rows int) {
	size := cols * rows
	if size < 0 {
		panic("grid size overflow")
	}
}

// validateBitmap validates that bitmap is not nil.
// Panics if b is nil.
func validateBitmap(b *Bitmap) {
	if b == nil {
		panic("nil bitmap")
	}
}

// validateColsZero validates operation when cols == 0.
// Panics if cols == 0 and w > 0.
func validateColsZero(cols, w int) {
	if cols == 0 && w > 0 {
		panic("Cols == 0")
	}
}

// validateRectInCols validates rectangle fits within columns.
// Panics if x+w > cols.
func (g *Grid) validateRectInCols(x, w int) {
	if x+w > g.cols {
		panic("rectangle exceeds columns")
	}
}

// validateRectInRows validates rectangle fits within rows.
// Panics if y+h > rows.
func (g *Grid) validateRectInRows(y, h int) {
	if y+h > g.Rows() {
		panic("rectangle exceeds rows")
	}
}

// validateSetRect validates complete SetRect operation.
// Validates params are non-negative, cols != 0 if w > 0, and rect fits in columns.
// Does NOT validate row bounds since SetRect auto-grows.
func (g *Grid) validateSetRect(x, y, w, h int) {
	validateRectParams(x, y, w, h)
	validateColsZero(g.cols, w)
	g.validateRectInCols(x, w)
}

// validateClearRect validates complete ClearRect operation.
// Validates params are non-negative and rect fits in both dimensions.
func (g *Grid) validateClearRect(x, y, w, h int) {
	validateRectParams(x, y, w, h)
	g.validateRectInCols(x, w)
	g.validateRectInRows(y, h)
}
