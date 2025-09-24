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
	if cols < 0 {
		panic("NewGrid: negative cols")
	}
	return &Grid{B: New(), cols: cols}
}

// NewGridWithCap returns a Grid with capacity for rowsCap rows.
// Accepts cols == 0 or rowsCap == 0.
func NewGridWithCap(cols, rowsCap int) *Grid {
	if cols < 0 || rowsCap < 0 {
		panic("NewGridWithCap: negative input")
	}
	return &Grid{B: NewWithCap(rowsCap * cols), cols: cols}
}

// NewGridWithSize returns a Grid sized to rows*cols bits.
// Accepts cols == 0 or rows == 0.
func NewGridWithSize(cols, rows int) *Grid {
	if cols < 0 || rows < 0 {
		panic("NewGridWithSize: negative input")
	}
	g := &Grid{B: New(), cols: cols}
	g.B.EnsureBits(rows * cols)
	return g
}

// NewGridFrom wraps an existing Bitmap. Panics if b is nil.
func NewGridFrom(b *Bitmap, cols int) *Grid {
	if b == nil {
		panic("NewGridFrom: nil bitmap")
	}
	if cols < 0 {
		panic("NewGridFrom: negative cols")
	}
	return &Grid{B: b, cols: cols}
}

// Index returns y*Cols + x. Panics on negative x or y.
func (g *Grid) Index(x, y int) int {
	if x < 0 || y < 0 {
		panic("Index: negative")
	}
	return y*g.cols + x
}

// Cols returns the number of columns.
func (g *Grid) Cols() int { return g.cols }

// Rows reports Len()/Cols. If Cols==0 or Len()==0, Rows==0.
func (g *Grid) Rows() int {
	if g.cols <= 0 || g.B.Len() == 0 {
		return 0
	}
	return g.B.Len() / g.cols
}

// SetRect sets to 1 a rectangle of size w×h at origin (x,y).
// Auto-grows rows to fit y+h. Panics if x<0, y<0, w<0, h<0, or x+w > Cols.
// Returns g.
//
// Invariant: after return, g.B.Len() == g.Rows()*g.Cols() and tail bits are masked.
func (g *Grid) SetRect(x, y, w, h int) *Grid {
	if x < 0 || y < 0 || w < 0 || h < 0 {
		panic("SetRect: negative input")
	}
	if g.cols == 0 && w > 0 {
		panic("SetRect: Cols == 0")
	}
	if x+w > g.cols {
		panic("SetRect: rectangle exceeds columns")
	}
	needRows := y + h
	if needRows > g.Rows() {
		g.B.EnsureBits(needRows * g.cols)
	}
	for row := range h {
		start := (y+row)*g.cols + x
		g.B.SetRange(start, w)
	}
	return g
}

// ClearRect clears to 0 a rectangle of size w×h at origin (x,y).
// Panics if rectangle exceeds current Rows() or Cols(). Returns g.
//
// Invariant: after return, g.B.Len() == g.Rows()*g.Cols() and tail bits are masked.
func (g *Grid) ClearRect(x, y, w, h int) *Grid {
	if x < 0 || y < 0 || w < 0 || h < 0 {
		panic("ClearRect: negative input")
	}
	if x+w > g.cols {
		panic("ClearRect: rectangle exceeds columns")
	}
	if y+h > g.Rows() {
		panic("ClearRect: rectangle exceeds rows")
	}
	for row := range h {
		start := (y+row)*g.cols + x
		g.B.ClearRange(start, w)
	}
	return g
}

// GrowCols increases Cols by delta (>0) and repositions existing rows so each
// cell (x,y) remains at the same coordinates under the new Cols.
// Newly created columns are zero. Returns g.
//
// Invariant: after return, g.B.Len() == g.Rows()*g.Cols() and tail bits are masked.
func (g *Grid) GrowCols(delta int) *Grid {
	if delta <= 0 {
		panic("GrowCols: delta must be > 0")
	}
	oldCols := g.cols
	newCols := oldCols + delta
	rows := g.Rows()
	if rows == 0 {
		g.cols = newCols
		return g
	}
	// Resize backing store to new size.
	g.B.EnsureBits(rows * newCols)

	// Move each row to its new stride from bottom to top to avoid overlap issues.
	if oldCols > 0 {
		for r := rows - 1; r >= 0; r-- {
			srcStart := r * oldCols
			dstStart := r * newCols
			g.B.CopyRange(g.B, srcStart, dstStart, oldCols)
			// New columns are already zero because EnsureBits zero-initializes newly added space.
		}
	}
	g.cols = newCols
	return g
}

// EnsureCols grows Cols to at least cols, repositioning like GrowCols when needed.
// No-op if cols <= Cols. Returns g.
//
// Invariant: after return, g.B.Len() == g.Rows()*g.Cols() and tail bits are masked.
func (g *Grid) EnsureCols(cols int) *Grid {
	if cols < 0 {
		panic("EnsureCols: negative cols")
	}
	if cols <= g.cols {
		return g
	}
	return g.GrowCols(cols - g.cols)
}

// GrowRows appends delta (>0) empty rows below current content. Returns g.
//
// Invariant: after return, g.B.Len() == g.Rows()*g.Cols() and tail bits are masked.
func (g *Grid) GrowRows(delta int) *Grid {
	if delta <= 0 {
		panic("GrowRows: delta must be > 0")
	}
	if g.cols == 0 {
		// No columns, rows are meaningless. Only update length 0.
		return g
	}
	newRows := g.Rows() + delta
	g.B.EnsureBits(newRows * g.cols)
	return g
}

// EnsureRows ensures at least rows rows exist. No repositioning. Returns g.
//
// Invariant: after return, g.B.Len() == g.Rows()*g.Cols() and tail bits are masked.
func (g *Grid) EnsureRows(rows int) *Grid {
	if rows < 0 {
		panic("EnsureRows: negative rows")
	}
	if g.cols == 0 {
		return g
	}
	if rows <= g.Rows() {
		return g
	}
	g.B.EnsureBits(rows * g.cols)
	return g
}
