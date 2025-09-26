package btmp

// Grid is a zero-copy row-major view over a Bitmap.
// Cols is the fixed number of columns per row. Grid mutators keep
// Len() == Rows()*Cols after each operation.
type Grid struct {
	B    *Bitmap
	cols int
}

// NewGrid returns a Grid with a new underlying Bitmap and the given column count.
// Accepts cols == 0.
func NewGrid(cols int) *Grid {
	validateCols(cols)
	return &Grid{
		B:    New(0),
		cols: cols,
	}
}

// NewGridWithSize returns a Grid sized to rows*cols bits.
// Accepts cols == 0 or rows == 0.
func NewGridWithSize(cols, rows int) *Grid {
	validateCols(cols)
	validateRows(rows)
	validateGridSize(cols, rows)

	size := cols * rows
	return &Grid{
		B:    New(uint(size)),
		cols: cols,
	}
}

// NewGridFrom wraps an existing Bitmap. Panics if b is nil.
func NewGridFrom(b *Bitmap, cols int) *Grid {
	validateBitmap(b)
	validateCols(cols)
	return &Grid{
		B:    b,
		cols: cols,
	}
}

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
	validateCoordinate(x, y)
	return y*g.cols + x
}

// EnsureCols grows Cols to at least cols, repositioning like GrowCols when needed.
// No-op if cols <= Cols. Returns g.
func (g *Grid) EnsureCols(cols int) *Grid {
	validateCols(cols)
	g.ensureCols(cols)
	return g
}

// EnsureRows ensures at least rows rows exist. No repositioning. Returns g.
func (g *Grid) EnsureRows(rows int) *Grid {
	validateRows(rows)
	g.ensureRows(rows)
	return g
}

// GrowCols increases Cols by delta (>0) and repositions existing rows so each
// cell (x,y) remains at the same coordinates under the new Cols.
// Newly created columns are zero. Returns g.
func (g *Grid) GrowCols(delta int) *Grid {
	validateDelta(delta)
	g.growCols(delta)
	return g
}

// GrowRows appends delta (>0) empty rows below current content. Returns g.
func (g *Grid) GrowRows(delta int) *Grid {
	validateDelta(delta)
	g.growRows(delta)
	return g
}

// SetRect sets to 1 a rectangle of size w×h at origin (x,y).
// Auto-grows rows to fit y+h. Panics if x<0, y<0, w<0, h<0, or x+w > Cols.
// Returns g.
func (g *Grid) SetRect(x, y, w, h int) *Grid {
	g.validateSetRect(x, y, w, h)
	g.setRect(x, y, w, h)
	return g
}

// ClearRect clears to 0 a rectangle of size w×h at origin (x,y).
// Panics if rectangle exceeds current Rows() or Cols(). Returns g.
func (g *Grid) ClearRect(x, y, w, h int) *Grid {
	g.validateClearRect(x, y, w, h)
	g.clearRect(x, y, w, h)
	return g
}
