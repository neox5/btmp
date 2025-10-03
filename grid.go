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

// NewGrid returns a Grid with a new empty underlying Bitmap and zero columns.
// The grid must be configured with GrowCols or EnsureCols before use.
func NewGrid() *Grid {
	return &Grid{
		B:    New(0),
		cols: 0,
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
	return g.isFree(x, y, w, h)
}

// CanShiftRight reports whether the rectangle can shift one column right.
// Checks if column x+w exists and is free (all zeros) for rows [y, y+h).
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) CanShiftRight(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)
	return g.canShiftRight(x, y, w, h)
}

// CanShiftLeft reports whether the rectangle can shift one column left.
// Checks if column x-1 exists and is free (all zeros) for rows [y, y+h).
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) CanShiftLeft(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)
	return g.canShiftLeft(x, y, w, h)
}

// CanShiftUp reports whether the rectangle can shift one row up.
// Checks if row y-1 exists and is free (all zeros) for columns [x, x+w).
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) CanShiftUp(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)
	return g.canShiftUp(x, y, w, h)
}

// CanShiftDown reports whether the rectangle can shift one row down.
// Checks if row y+h exists and is free (all zeros) for columns [x, x+w).
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) CanShiftDown(x, y, w, h int) bool {
	g.validateRect(x, y, w, h)
	return g.canShiftDown(x, y, w, h)
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
// or target column is not free.
func (g *Grid) ShiftRectRight(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	if !g.canShiftRight(x, y, w, h) {
		panic("cannot shift right")
	}
	g.shiftRectRight(x, y, w, h)
	return g
}

// ShiftRectLeft shifts a rectangle one column to the left.
// Moves bits from [x,y,w,h) to [x-1,y,w,h) and clears the rightmost column.
// Target column (x-1) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// or target column is not free.
func (g *Grid) ShiftRectLeft(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	if !g.canShiftLeft(x, y, w, h) {
		panic("cannot shift left")
	}
	g.shiftRectLeft(x, y, w, h)
	return g
}

// ShiftRectUp shifts a rectangle one row up.
// Moves bits from [x,y,w,h) to [x,y-1,w,h) and clears the bottom row.
// Target row (y-1) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// or target row is not free.
func (g *Grid) ShiftRectUp(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	if !g.canShiftUp(x, y, w, h) {
		panic("cannot shift up")
	}
	g.shiftRectUp(x, y, w, h)
	return g
}

// ShiftRectDown shifts a rectangle one row down.
// Moves bits from [x,y,w,h) to [x,y+1,w,h) and clears the top row.
// Target row (y+h) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// or target row is not free.
func (g *Grid) ShiftRectDown(x, y, w, h int) *Grid {
	g.validateRect(x, y, w, h)
	if !g.canShiftDown(x, y, w, h) {
		panic("cannot shift down")
	}
	g.shiftRectDown(x, y, w, h)
	return g
}

// ========================================
// Print Operations
// ========================================

// Print formats the grid as a coordinate-labeled visualization.
// Each row is prefixed with its row number, and column indices are shown at the top.
// Uses '.' for zero bits and '#' for set bits.
// Returns empty string if grid has no rows or columns.
//
// Example output for a 5x3 grid with bits set at (1,0) and (3,1):
//   0 1 2 3 4
// 0 . # . . .
// 1 . . . # .
// 2 . . . . .
func (g *Grid) Print() string {
	return g.print()
}
