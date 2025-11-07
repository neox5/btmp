package btmp

// Grid is a zero-copy row-major view over a Bitmap.
// Cols is the fixed number of columns per row. Grid mutators keep
// Len() == Rows()*Cols after each operation.
type Grid struct {
	B    *Bitmap
	cols int
	rows int
}

// ========================================
// Constructor Functions
// ========================================

// NewGrid returns a Grid with zero rows and columns.
func NewGrid() *Grid {
	return NewGridWithSize(0, 0)
}

// NewGridWithSize returns a Grid sized to rows*cols bits.
// Accepts rows == 0 or cols == 0. Panics if rows < 0, cols < 0, or size overflows.
func NewGridWithSize(rows, cols int) *Grid {
	if err := validateNonNegative(rows, "rows"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NewGridWithSize"))
	}
	if err := validateNonNegative(cols, "cols"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NewGridWithSize"))
	}
	if err := validateGridSizeMax(rows, cols); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NewGridWithSize"))
	}

	size := rows * cols
	return &Grid{
		B:    New(uint(size)),
		cols: cols,
		rows: rows,
	}
}

// ========================================
// Accessors
// ========================================

// Cols returns the number of columns.
func (g *Grid) Cols() int {
	return g.cols
}

// Rows returns the number of rows.
func (g *Grid) Rows() int {
	return g.rows
}

// Index returns r*Cols + c. Panics on negative r or c.
func (g *Grid) Index(r, c int) int {
	if err := validateNonNegative(r, "r"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.Index"))
	}
	if err := validateNonNegative(c, "c"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.Index"))
	}
	return r*g.cols + c
}

// ========================================
// Growth Operations
// ========================================

// EnsureCols grows Cols to at least cols, repositioning like GrowCols when needed.
// No-op if cols <= Cols. Returns g. Panics if cols < 0.
func (g *Grid) EnsureCols(cols int) *Grid {
	if err := validateNonNegative(cols, "cols"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.EnsureCols"))
	}
	g.ensureCols(cols)
	return g
}

// EnsureRows ensures at least rows rows exist. No repositioning. Returns g.
// Panics if rows < 0.
func (g *Grid) EnsureRows(rows int) *Grid {
	if err := validateNonNegative(rows, "rows"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.EnsureRows"))
	}
	g.ensureRows(rows)
	return g
}

// GrowCols increases Cols by delta (>0) and repositions existing rows so each
// cell (r,c) remains at the same coordinates under the new Cols.
// Newly created columns are zero. Returns g. Panics if delta < 0.
func (g *Grid) GrowCols(delta int) *Grid {
	if err := validateNonNegative(delta, "delta"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.GrowCols"))
	}
	if delta > 0 {
		g.growCols(delta)
	}
	return g
}

// GrowRows appends delta (>0) empty rows below current content. Returns g.
// Panics if delta < 0.
func (g *Grid) GrowRows(delta int) *Grid {
	if err := validateNonNegative(delta, "delta"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.GrowRows"))
	}
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
func (g *Grid) IsFree(r, c, h, w int) bool {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.IsFree"))
	}
	return g.isFree(r, c, h, w)
}

// NextZeroInRow returns the column index of the next zero bit in row r,
// starting search from column c.
// Search is constrained to row r only - does not continue to next row.
// Returns -1 if no zero bit exists in [c, Cols()).
// Panics if r < 0, c < 0, r >= Rows(), or c >= Cols().
func (g *Grid) NextZeroInRow(r, c int) int {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NextZeroInRow"))
	}
	return g.nextZeroInRow(r, c)
}

// NextOneInRow returns the column index of the next set bit in row r,
// starting search from column c.
// Search is constrained to row r only - does not continue to next row.
// Returns -1 if no set bit exists in [c, Cols()).
// Panics if r < 0, c < 0, r >= Rows(), or c >= Cols().
func (g *Grid) NextOneInRow(r, c int) int {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NextOneInRow"))
	}
	return g.nextOneInRow(r, c)
}

// NextZeroInRowRange returns the column index of the next zero bit in row r,
// searching within [c, c+count).
// Search is constrained to specified range only.
// Returns -1 if no zero bit exists in range.
// Panics if r < 0, c < 0, count <= 0, r >= Rows(), or c >= Cols().
func (g *Grid) NextZeroInRowRange(r, c, count int) int {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NextZeroInRowRange"))
	}
	if err := validatePositive(count, "count"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NextZeroInRowRange"))
	}
	return g.nextZeroInRowRange(r, c, count)
}

// NextOneInRowRange returns the column index of the next set bit in row r,
// searching within [c, c+count).
// Search is constrained to specified range only.
// Returns -1 if no set bit exists in range.
// Panics if r < 0, c < 0, count <= 0, r >= Rows(), or c >= Cols().
func (g *Grid) NextOneInRowRange(r, c, count int) int {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NextOneInRowRange"))
	}
	if err := validatePositive(count, "count"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.NextOneInRowRange"))
	}
	return g.nextOneInRowRange(r, c, count)
}

// CountZerosFromInRow returns the count of consecutive zero bits in row r
// starting at column c.
// Count is constrained to row r only - stops at Cols() boundary.
// Returns 0 if bit at (r,c) is set.
// Stops at first set bit or end of row.
// Panics if r < 0, c < 0, r >= Rows(), or c >= Cols().
func (g *Grid) CountZerosFromInRow(r, c int) int {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CountZerosFromInRow"))
	}
	return g.countZerosFromInRow(r, c)
}

// CountOnesFromInRow returns the count of consecutive set bits in row r
// starting at column c.
// Count is constrained to row r only - stops at Cols() boundary.
// Returns 0 if bit at (r,c) is zero.
// Stops at first zero bit or end of row.
// Panics if r < 0, c < 0, r >= Rows(), or c >= Cols().
func (g *Grid) CountOnesFromInRow(r, c int) int {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CountOnesFromInRow"))
	}
	return g.countOnesFromInRow(r, c)
}

// CountZerosFromInRowRange returns the count of consecutive zero bits in row r
// starting at column c, within [c, c+count).
// Count is constrained to specified range only.
// Returns 0 if bit at (r,c) is set.
// Stops at first set bit or end of range.
// Panics if r < 0, c < 0, count <= 0, r >= Rows(), or c >= Cols().
func (g *Grid) CountZerosFromInRowRange(r, c, count int) int {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CountZerosFromInRowRange"))
	}
	if err := validatePositive(count, "count"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CountZerosFromInRowRange"))
	}
	return g.countZerosFromInRowRange(r, c, count)
}

// CountOnesFromInRowRange returns the count of consecutive set bits in row r
// starting at column c, within [c, c+count).
// Count is constrained to specified range only.
// Returns 0 if bit at (r,c) is zero.
// Stops at first zero bit or end of range.
// Panics if r < 0, c < 0, count <= 0, r >= Rows(), or c >= Cols().
func (g *Grid) CountOnesFromInRowRange(r, c, count int) int {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CountOnesFromInRowRange"))
	}
	if err := validatePositive(count, "count"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CountOnesFromInRowRange"))
	}
	return g.countOnesFromInRowRange(r, c, count)
}

// CanShiftRight reports whether the rectangle can shift one column right.
// Checks if column c+w exists and is free (all zeros) for rows [r, r+h).
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) CanShiftRight(r, c, h, w int) bool {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanShiftRight"))
	}
	return g.canShiftRight(r, c, h, w)
}

// CanShiftLeft reports whether the rectangle can shift one column left.
// Checks if column c-1 exists and is free (all zeros) for rows [r, r+h).
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) CanShiftLeft(r, c, h, w int) bool {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanShiftLeft"))
	}
	return g.canShiftLeft(r, c, h, w)
}

// CanShiftUp reports whether the rectangle can shift one row up.
// Checks if row r-1 exists and is free (all zeros) for columns [c, c+w).
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) CanShiftUp(r, c, h, w int) bool {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanShiftUp"))
	}
	return g.canShiftUp(r, c, h, w)
}

// CanShiftDown reports whether the rectangle can shift one row down.
// Checks if row r+h exists and is free (all zeros) for columns [c, c+w).
// Panics if rectangle is invalid or out of bounds.
func (g *Grid) CanShiftDown(r, c, h, w int) bool {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanShiftDown"))
	}
	return g.canShiftDown(r, c, h, w)
}

// CanFitWidth reports whether a cell with width w can fit in row r starting at
// column c (i.e., whether columns [c, c+w) contain only unoccupied cells).
// Returns false if any cell in the range is occupied or if c+w exceeds Cols().
// Panics if r < 0, c < 0, w <= 0, r >= Rows(), or c >= Cols().
func (g *Grid) CanFitWidth(r, c, w int) bool {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanFitWidth"))
	}
	if err := validatePositive(w, "w"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanFitWidth"))
	}
	return g.canFitWidth(r, c, w)
}

// CanFit reports whether a rectangle of size h×w can fit at position (r, c) in the grid.
// This checks only boundary constraints, not whether the cells are occupied.
// Returns two booleans:
//   - fitRow: true if r+h <= Rows() (height fits within grid)
//   - fitCol: true if c+w <= Cols() (width fits within grid)
//
// Panics if r < 0, c < 0, h < 0, w < 0, r >= Rows(), or c >= Cols().
func (g *Grid) CanFit(r, c, h, w int) (fitRow, fitCol bool) {
	if err := g.validateCoordinate(r, c); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanFit"))
	}
	if err := validateNonNegative(h, "h"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanFit"))
	}
	if err := validateNonNegative(w, "w"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.CanFit"))
	}
	return g.canFit(r, c, h, w)
}

// AllGrid returns true if all bits in the grid are set.
// Returns false for empty grid (0 rows or 0 columns).
func (g *Grid) AllGrid() bool {
	return g.allGrid()
}

// AllRow returns true if all bits in row r are set.
// Returns false for empty row (Cols() == 0).
// Panics if r < 0 or r >= Rows().
func (g *Grid) AllRow(r int) bool {
	if err := validateNonNegative(r, "r"); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.AllRow"))
	}
	if r >= g.rows {
		panic(&ValidationError{
			Field:   "r",
			Value:   r,
			Message: "out of bounds",
			Context: "Grid.AllRow",
		})
	}
	return g.allRow(r)
}

// ========================================
// Validation Operations
// ========================================

// ValidateCoordinate validates that r and c are non-negative and within grid bounds.
// Returns ValidationError if r < 0, c < 0, r >= g.Rows(), or c >= g.Cols().
func (g *Grid) ValidateCoordinate(r, c int) error {
	return g.validateCoordinate(r, c)
}

// ValidateRect validates that rectangle parameters are non-negative
// and rectangle is fully contained within grid bounds.
// Returns ValidationError if r < 0, c < 0, h < 0, w < 0, r+h > g.Rows(), or c+w > g.Cols().
func (g *Grid) ValidateRect(r, c, h, w int) error {
	return g.validateRect(r, c, h, w)
}

// ========================================
// Rectangle Mutators
// ========================================

// SetRect sets to 1 a rectangle of size h×w at origin (r,c).
// All coordinates must be in bounds. Panics if r<0, c<0, h<0, w<0,
// r+h > Rows, or c+w > Cols.
// Returns *Grid for chaining.
func (g *Grid) SetRect(r, c, h, w int) *Grid {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.SetRect"))
	}
	g.setRect(r, c, h, w)
	return g
}

// ClearRect clears to 0 a rectangle of size h×w at origin (r,c).
// Panics if rectangle exceeds current Rows() or Cols(). Returns g.
func (g *Grid) ClearRect(r, c, h, w int) *Grid {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.ClearRect"))
	}
	g.clearRect(r, c, h, w)
	return g
}

// ShiftRectRight shifts a rectangle one column to the right.
// Moves bits from [r,c,h,w) to [r,c+1,h,w) and clears the leftmost column.
// Target column (c+w) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// or target column is not free.
func (g *Grid) ShiftRectRight(r, c, h, w int) *Grid {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.ShiftRectRight"))
	}
	if !g.canShiftRight(r, c, h, w) {
		panic(&ValidationError{
			Field:   "shift",
			Value:   "right",
			Message: "target column not free",
			Context: "Grid.ShiftRectRight",
		})
	}
	g.shiftRectRight(r, c, h, w)
	return g
}

// ShiftRectLeft shifts a rectangle one column to the left.
// Moves bits from [r,c,h,w) to [r,c-1,h,w) and clears the rightmost column.
// Target column (c-1) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// or target column is not free.
func (g *Grid) ShiftRectLeft(r, c, h, w int) *Grid {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.ShiftRectLeft"))
	}
	if !g.canShiftLeft(r, c, h, w) {
		panic(&ValidationError{
			Field:   "shift",
			Value:   "left",
			Message: "target column not free",
			Context: "Grid.ShiftRectLeft",
		})
	}
	g.shiftRectLeft(r, c, h, w)
	return g
}

// ShiftRectUp shifts a rectangle one row up.
// Moves bits from [r,c,h,w) to [r-1,c,h,w) and clears the bottom row.
// Target row (r-1) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// or target row is not free.
func (g *Grid) ShiftRectUp(r, c, h, w int) *Grid {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.ShiftRectUp"))
	}
	if !g.canShiftUp(r, c, h, w) {
		panic(&ValidationError{
			Field:   "shift",
			Value:   "up",
			Message: "target row not free",
			Context: "Grid.ShiftRectUp",
		})
	}
	g.shiftRectUp(r, c, h, w)
	return g
}

// ShiftRectDown shifts a rectangle one row down.
// Moves bits from [r,c,h,w) to [r+1,c,h,w) and clears the top row.
// Target row (r+h) must exist and be free (all zeros).
// Returns *Grid for chaining. Panics if rectangle is invalid, out of bounds,
// or target row is not free.
func (g *Grid) ShiftRectDown(r, c, h, w int) *Grid {
	if err := g.validateRect(r, c, h, w); err != nil {
		panic(err.(*ValidationError).WithContext("Grid.ShiftRectDown"))
	}
	if !g.canShiftDown(r, c, h, w) {
		panic(&ValidationError{
			Field:   "shift",
			Value:   "down",
			Message: "target row not free",
			Context: "Grid.ShiftRectDown",
		})
	}
	g.shiftRectDown(r, c, h, w)
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
// Example output for a 3x5 grid (3 rows, 5 cols) with bits set at (0,1) and (1,3):
//
//	  0 1 2 3 4
//	0 . # . . .
//	1 . . . # .
//	2 . . . . .
func (g *Grid) Print() string {
	return g.print()
}
