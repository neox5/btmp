package btmp

// Grid is a zero-copy row-major view over a Bitmap.
// Cols is the fixed number of columns per row. Grid mutators keep
// Len() == Rows()*Cols after each operation.
type Grid struct {
	B    *Bitmap
	cols int
}

// NewGrid returns a Grid with a new underlying Bitmap and the given column count.
// Accepts cols==0.
func NewGrid(cols int) *Grid {
	return &Grid{B: New(), cols: cols}
}

// NewGridWithCap returns a Grid with a new Bitmap and reserves capacity
// for rowsCap rows (Len() unchanged). Accepts cols==0 or rowsCap==0.
func NewGridWithCap(cols, rowsCap int) *Grid {
	return &Grid{B: NewWithCap(cols * rowsCap), cols: cols}
}

// NewGridWithSize returns a Grid with a new Bitmap and sets Len() to rows*cols bits
// (growing zeroed as needed). Accepts cols==0 or rows==0.
func NewGridWithSize(cols, rows int) *Grid {
	if cols < 0 || rows < 0 {
		panic("btmp: negative cols or rows")
	}
	g := &Grid{B: New(), cols: cols}
	g.B = g.B.EnsureBits(cols * rows)
	return g
}

// NewGridFrom wraps an existing Bitmap without allocation. B must not be nil.
func NewGridFrom(b *Bitmap, cols int) *Grid {
	if b == nil {
		panic("btmp: NewGridFrom nil Bitmap")
	}
	return &Grid{B: b, cols: cols}
}

// Index returns the linear bit index for (x,y): y*Cols + x.
// Panics on negative x or y. Does not check against Len().
func (g *Grid) Index(x, y int) int {
	if x < 0 || y < 0 {
		panic("btmp: Grid.Index negative coordinate")
	}
	return y*g.cols + x
}

// Cols returns the number of columns of the grid.
func (g *Grid) Cols() int {
	return g.cols
}

// Rows reports Len()/Cols. If Cols==0 or Len()==0, Rows==0.
func (g *Grid) Rows() int {
	if g.cols <= 0 || g.B.Len() == 0 {
		return 0
	}
	return g.B.Len() / g.cols
}

// SetRect sets to 1 a rectangle of size w×h at origin (x,y).
// Auto-grows rows as needed to fit y+h, but does not change Cols.
// Panics if x<0, y<0, w<0, h<0, or if x+w > Cols.
// Returns g for chaining.
func (g *Grid) SetRect(x, y, w, h int) *Grid {
	if x < 0 || y < 0 || w < 0 || h < 0 {
		panic("btmp: SetRect negative args")
	}
	if w == 0 || h == 0 {
		return g
	}
	if g.cols <= 0 {
		panic("btmp: SetRect requires Cols > 0")
	}
	if x+w > g.cols {
		panic("btmp: SetRect exceeds Cols")
	}
	needRows := y + h
	g.B = g.B.EnsureBits(needRows * g.cols)
	start := g.Index(x, y)
	for r := range h {
		g.B = g.B.SetRange(start+r*g.cols, w)
	}
	return g
}

// ClearRect clears to 0 a rectangle of size w×h at origin (x,y).
// Panics if any part of the rectangle exceeds current Rows() or Cols().
// Returns g for chaining.
func (g *Grid) ClearRect(x, y, w, h int) *Grid {
	if x < 0 || y < 0 || w < 0 || h < 0 {
		panic("btmp: ClearRect negative args")
	}
	if w == 0 || h == 0 {
		return g
	}
	if g.cols <= 0 {
		panic("btmp: ClearRect requires Cols > 0")
	}
	if x+w > g.cols {
		panic("btmp: ClearRect exceeds Cols")
	}
	if y+h > g.Rows() {
		panic("btmp: ClearRect exceeds Rows")
	}
	start := g.Index(x, y)
	for r := range h {
		g.B = g.B.ClearRange(start+r*g.cols, w)
	}
	return g
}

// GrowCols increases Cols by delta (>0) and repositions existing rows so that
// each cell (x,y) remains at the same coordinates under the new Cols.
// After return, newly created columns are zero and Len()==Rows()*Cols.
// Returns g for chaining.
func (g *Grid) GrowCols(delta int) *Grid {
	if delta <= 0 {
		panic("btmp: GrowCols delta must be > 0")
	}
	oldC := g.cols
	newC := oldC + delta
	if oldC == 0 {
		g.cols = newC
		return g
	}
	rows := g.Rows()
	g.B = g.B.EnsureBits(rows * newC)
	// Bottom-up to avoid overlap.
	for y := rows - 1; y >= 0; y-- {
		src := y * oldC
		dst := y * newC
		g.B = g.B.CopyRange(g.B, src, dst, oldC)
	}
	g.cols = newC
	return g
}

// EnsureCols grows Cols to at least cols and performs the same repositioning
// as GrowCols(cols - Cols) when cols > Cols. No-op if cols <= Cols.
// Returns g for chaining.
func (g *Grid) EnsureCols(cols int) *Grid {
	if cols < 0 {
		panic("btmp: EnsureCols negative")
	}
	if cols <= g.cols {
		return g
	}
	return g.GrowCols(cols - g.cols)
}

// GrowRows appends delta (>0) empty rows below current content. No repositioning.
// After return, Len()==Rows()*Cols and new rows are zero. Returns g for chaining.
func (g *Grid) GrowRows(delta int) *Grid {
	if delta <= 0 {
		panic("btmp: GrowRows delta must be > 0")
	}
	if g.cols < 0 {
		panic("btmp: GrowRows Cols negative")
	}
	if g.cols == 0 {
		return g
	}
	newRows := g.Rows() + delta
	g.B = g.B.EnsureBits(newRows * g.cols)
	return g
}

// EnsureRows ensures at least rows rows exist. No repositioning.
// After return, Len()==Rows()*Cols. No-op if rows <= current Rows().
// Returns g for chaining.
func (g *Grid) EnsureRows(rows int) *Grid {
	if rows < 0 {
		panic("btmp: EnsureRows negative")
	}
	if g.cols == 0 {
		return g
	}
	cur := g.Rows()
	if rows <= cur {
		return g
	}
	g.B = g.B.EnsureBits(rows * g.cols)
	return g
}
