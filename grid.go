package btmp

// Grid is a zero-copy row-major view over a Bitmap.
// Cols is the fixed number of columns per row. Grid mutators keep
// Len() == Rows()*Cols after each operation.
type Grid struct {
	B    *Bitmap
	cols int
}

// ========================================
// Constructor Functions
// ========================================

// NewGrid returns a Grid with a new underlying Bitmap and the given column count.
// Accepts cols == 0. Panics if cols < 0.
func NewGrid(cols int) *Grid {
	validateNonNegative(cols, "cols")
	return &Grid{
		B:    New(0),
		cols: cols,
	}
}

// NewGridWithSize returns a Grid sized to rows*cols bits.
// Accepts cols == 0 or rows == 0. Panics if cols < 0, rows < 0, or size overflows.
func NewGridWithSize(cols, rows int) *Grid {
	validateNonNegative(cols, "cols")
	validateNonNegative(rows, "rows")
	validateGridSizeMax(rows, cols)

	size := cols * rows
	return &Grid{
		B:    New(uint(size)),
		cols: cols,
	}
}

// NewGridFrom wraps an existing Bitmap. Panics if b is nil or cols < 0.
func NewGridFrom(b *Bitmap, cols int) *Grid {
	validateNotNil(b, "b")
	validateNonNegative(cols, "cols")
	return &Grid{
		B:    b,
		cols: cols,
	}
}

// ========================================
// Accessors
// ========================================

// Cols returns the number of columns.
func (g *Grid) Cols() int {
	return g.cols
}

// Rows reports Len()/Cols. If Cols==0 or Len()==0, Rows==0.
func (g *Grid) Rows() int {
	if g.cols == 0 || g.B.Len() == 0 {
		return 0
	}
	return g.B.Len() / g.cols
}

// Index returns y*Cols + x. Panics on negative x or y.
func (g *Grid) Index(x, y int) int {
	validateNonNegative(x, "x")
	validateNonNegative(y, "y")
	return y*g.cols + x
}

// ========================================
// Growth Operations
// ========================================

// EnsureCols grows Cols to at least cols, repositioning like GrowCols when needed.
// No-op if cols <= Cols. Returns g. Panics if cols < 0.
func (g *Grid) EnsureCols(cols int) *Grid {
	validateNonNegative(cols, "cols")
	g.ensureCols(cols)
	return g
}

// EnsureRows ensures at least rows rows exist. No repositioning. Returns g.
// Panics if rows < 0.
func (g *Grid) EnsureRows(rows int) *Grid {
	validateNonNegative(rows, "rows")
	g.ensureRows(rows)
	return g
}

// GrowCols increases Cols by delta (>0) and repositions existing rows so each
// cell (x,y) remains at the same coordinates under the new Cols.
// Newly created columns are zero. Returns g. Panics if delta < 0.
func (g *Grid) GrowCols(delta int) *Grid {
	validateNonNegative(delta, "delta")
	if delta > 0 {
		g.growCols(delta)
	}
	return g
}

// GrowRows appends delta (>0) empty rows below current content. Returns g.
// Panics if delta < 0.
func (g *Grid) GrowRows(delta int) *Grid {
	validateNonNegative(delta, "delta")
	if delta > 0 {
		g.growRows(delta)
	}
	return g
}

// ========================================
// Query Operations
// ========================================

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
	targetCol := x + w
	if targetCol >= g.cols {
		panic("target column out of bounds")
	}

	return g.IsFree(targetCol, y, 1, h)
}

// CanShiftLeft reports whether the rectangle can shift one column left.
// Checks if column x-1 exists and is free (all zeros) for rows [y, y+h).
// Panics if rectangle is invalid, out of bounds, or target column doesn't exist.
func (g *Grid) CanShiftLeft(x, y, w, h int) bool {
	if x == 0 {
		panic("target column out of bounds")
	}

	targetCol := x - 1
	return g.IsFree(targetCol, y, 1, h)
}

// CanShiftUp reports whether the rectangle can shift one row up.
// Checks if row y-1 exists and is free (all zeros) for columns [x, x+w).
// Panics if rectangle is invalid, out of bounds, or target row doesn't exist.
func (g *Grid) CanShiftUp(x, y, w, h int) bool {
	if y == 0 {
		panic("target row out of bounds")
	}

	targetRow := y - 1
	return g.IsFree(x, targetRow, w, 1)
}

// CanShiftDown reports whether the rectangle can shift one row down.
// Checks if row y+h exists and is free (all zeros) for columns [x, x+w).
// Panics if rectangle is invalid, out of bounds, or target row doesn't exist.
func (g *Grid) CanShiftDown(x, y, w, h int) bool {
	targetRow := y + h
	if targetRow >= g.Rows() {
		panic("target row out of bounds")
	}

	return g.IsFree(x, targetRow, w, 1)
}

// ========================================
// Rectangle Mutators
// ========================================

// SetRect sets to 1 a rectangle of size w×h at origin (x,y).
// All coordinates must be in bounds. Panics if x<0, y<0, w<0, h<0,
// x+w > Cols, or y+h > Rows.
// Returns *Grid for chaining.
func (g *Grid) SetRect(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	g.setRect(x, y, w, h)
	return g
}

// ClearRect clears to 0 a rectangle of size w×h at origin (x,y).
// Panics if rectangle exceeds current Rows() or Cols(). Returns g.
func (g *Grid) ClearRect(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	g.clearRect(x, y, w, h)
	return g
}

// ShiftRectRight shifts a rectangle one column to the right.
// Moves bits from [x,y,w,h) to [x+1,y,w,h) and clears the leftmost column.
// Target column (x+w) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// target column doesn't exist, or target column is not free.
func (g *Grid) ShiftRectRight(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	if !g.CanShiftRight(x, y, w, h) {
		panic("target column not free or out of bounds")
	}
	g.shiftRectRight(x, y, w, h)
	return g
}

// ShiftRectLeft shifts a rectangle one column to the left.
// Moves bits from [x,y,w,h) to [x-1,y,w,h) and clears the rightmost column.
// Target column (x-1) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// target column doesn't exist, or target column is not free.
func (g *Grid) ShiftRectLeft(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	if !g.CanShiftLeft(x, y, w, h) {
		panic("target column not free or out of bounds")
	}
	g.shiftRectLeft(x, y, w, h)
	return g
}

// ShiftRectUp shifts a rectangle one row up.
// Moves bits from [x,y,w,h) to [x,y-1,w,h) and clears the bottom row.
// Target row (y-1) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// target row doesn't exist, or target row is not free.
func (g *Grid) ShiftRectUp(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	if !g.CanShiftUp(x, y, w, h) {
		panic("target row not free or out of bounds")
	}
	g.shiftRectUp(x, y, w, h)
	return g
}

// ShiftRectDown shifts a rectangle one row down.
// Moves bits from [x,y,w,h) to [x,y+1,w,h) and clears the top row.
// Target row (y+h) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// target row doesn't exist, or target row is not free.
func (g *Grid) ShiftRectDown(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	if !g.CanShiftDown(x, y, w, h) {
		panic("target row not free or out of bounds")
	}
	g.shiftRectDown(x, y, w, h)
	return g
}
