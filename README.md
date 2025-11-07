<br/>
<br/>

<div align="center">
  <img src="logo.png" alt="btmp" width="500"/>
</div>

<br/>
<br/>

# btmp

btmp ("bitmap") is a pure Go bitmap library designed as a building block for your data structures. It provides tested, validated operations for manipulating dense boolean data without implementing bit math yourself.

Grid serves as a full-featured reference implementation, demonstrating how to build zero-copy 2D data structures on the bitmap foundation.

### When NOT to use

- **Sparse data** (< 1% density) → use `map[int]struct{}` or RoaringBitmap
- **Unknown/unpredictable growth** → use bits-and-blooms/bitset with auto-grow
- **Need compression** → use RoaringBitmap

### When to use

- **Building data structures** on top of bitmap operations
- **Size is known** or grows predictably
- **Explicit control** over bounds and validation
- **Range operations** - bulk sets, clears, copies

### Implementation

btmp abstracts away 64-bit word boundaries so you can work with bit positions directly. The library uses panics for validation failures - incorrect usage fails immediately at the source rather than propagating errors through your code.

## Install

```bash
go get github.com/neox5/btmp
```

## Quick start

```go
import "github.com/neox5/btmp"

// Create and manipulate bitmap
b := btmp.New(1024)

// Range operations (the key feature)
b.SetRange(100, 50)              // Set bits [100, 150)
b.ClearRange(110, 10)            // Clear bits [110, 120)
b.CopyRange(b, 0, 500, 100)      // Copy bits [0, 100) to [500, 600)

// Single-bit operations
b.SetBit(42).ClearBit(43).FlipBit(44)

// Multi-bit operations
b.SetBits(200, 8, 0xFF)          // Insert 8 bits at position 200

// Boolean operations
b2 := btmp.New(1024).SetRange(50, 100)
b.And(b2).Or(b2).Xor(b2).Not()

// Query operations
if b.Test(42) { /* bit is set */ }
if b.Any() { /* has any set bits */ }
if b.All() { /* all bits set */ }
count := b.Count()                    // Number of set bits
count = b.CountRange(100, 50)         // Count in range [100, 150)
if b.AnyRange(200, 10) { /* ... */ }  // Any bits set in range

// Grid - zero-copy 2D view (row-major)
g := btmp.NewGridWithSize(10, 16) // 10 rows, 16 columns
g.SetRect(2, 3, 4, 5)             // Set 4×5 rectangle at row 2, col 3
if g.IsFree(5, 8, 3, 3) {         // Check if 3×3 region is available at row 5, col 8
    g.SetRect(5, 8, 3, 3)
}
```

## Examples

The `examples/` directory contains complete working examples:

- **[bitmap_print](examples/bitmap_print/)** - Bitmap formatting and visualization (binary, hexadecimal, grouped output)
- **[grid_print](examples/grid_print/)** - Grid visualization and pattern creation

To run an example:

```bash
go run examples/bitmap_print/main.go
```

## API

### Bitmap (40 methods)

| Category             | Method                                                                                         |
| -------------------- | ---------------------------------------------------------------------------------------------- |
| **Construction** (1) | `New(n uint) *Bitmap`                                                                          |
| **Access** (2)       | `Len() int`                                                                                    |
|                      | `Words() []uint64`                                                                             |
| **Growth** (2)       | `EnsureBits(n int) *Bitmap`                                                                    |
|                      | `AddBits(n int) *Bitmap`                                                                       |
| **Query** (15)       | `Test(pos int) bool`                                                                           |
|                      | `Any() bool`                                                                                   |
|                      | `All() bool`                                                                                   |
|                      | `Count() int`                                                                                  |
|                      | `AnyRange(start, count int) bool`                                                              |
|                      | `AllRange(start, count int) bool`                                                              |
|                      | `CountRange(start, count int) int`                                                             |
|                      | `NextZero(pos int) int`                                                                        |
|                      | `NextOne(pos int) int`                                                                         |
|                      | `NextZeroInRange(pos, count int) int`                                                          |
|                      | `NextOneInRange(pos, count int) int`                                                           |
|                      | `CountZerosFrom(pos int) int`                                                                  |
|                      | `CountOnesFrom(pos int) int`                                                                   |
|                      | `CountZerosFromInRange(pos, count int) int`                                                    |
|                      | `CountOnesFromInRange(pos, count int) int`                                                     |
| **Validation** (2)   | `ValidateInBounds(pos int) error`                                                              |
|                      | `ValidateRange(start, count int) error`                                                        |
| **Single-bit** (3)   | `SetBit(pos int) *Bitmap`                                                                      |
|                      | `ClearBit(pos int) *Bitmap`                                                                    |
|                      | `FlipBit(pos int) *Bitmap`                                                                     |
| **Multi-bit** (1)    | `SetBits(pos, n int, val uint64) *Bitmap`                                                      |
| **Range** (4)        | `SetRange(start, count int) *Bitmap`                                                           |
|                      | `ClearRange(start, count int) *Bitmap`                                                         |
|                      | `CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap`                                |
|                      | `MoveRange(srcStart, dstStart, count int) *Bitmap`                                             |
| **Bulk** (2)         | `SetAll() *Bitmap`                                                                             |
|                      | `ClearAll() *Bitmap`                                                                           |
| **Logic** (4)        | `And(other *Bitmap) *Bitmap`                                                                   |
|                      | `Or(other *Bitmap) *Bitmap`                                                                    |
|                      | `Xor(other *Bitmap) *Bitmap`                                                                   |
|                      | `Not() *Bitmap`                                                                                |
| **Print** (4)        | `Print() string`                                                                               |
|                      | `PrintRange(start, count int) string`                                                          |
|                      | `PrintFormat(base int, grouped bool, groupSize int, sep string) string`                        |
|                      | `PrintRangeFormat(start, count int, base int, grouped bool, groupSize int, sep string) string` |

### Grid (29 methods)

| Category                   | Method                                          |
| -------------------------- | ----------------------------------------------- |
| **Construction** (2)       | `NewGrid() *Grid`                               |
|                            | `NewGridWithSize(rows, cols int) *Grid`         |
| **Access** (3)             | `Rows() int`                                    |
|                            | `Cols() int`                                    |
|                            | `Index(r, c int) int`                           |
| **Growth** (4)             | `EnsureRows(rows int) *Grid`                    |
|                            | `GrowRows(delta int) *Grid`                     |
|                            | `EnsureCols(cols int) *Grid`                    |
|                            | `GrowCols(delta int) *Grid`                     |
| **Query** (11)             | `RectZero(r, c, h, w int) bool`                 |
|                            | `RectOne(r, c, h, w int) bool`                  |
|                            | `NextZeroInRow(r, c int) int`                   |
|                            | `NextOneInRow(r, c int) int`                    |
|                            | `NextZeroInRowRange(r, c, count int) int`       |
|                            | `NextOneInRowRange(r, c, count int) int`        |
|                            | `CountZerosFromInRow(r, c int) int`             |
|                            | `CountOnesFromInRow(r, c int) int`              |
|                            | `CountZerosFromInRowRange(r, c, count int) int` |
|                            | `CountOnesFromInRowRange(r, c, count int) int`  |
|                            | `AllRow(r int) bool`                            |
| **Validation** (2)         | `ValidateCoordinate(r, c int) error`            |
|                            | `ValidateRect(r, c, h, w int) error`            |
| **Rectangle Mutators** (6) | `SetRect(r, c, h, w int) *Grid`                 |
|                            | `ClearRect(r, c, h, w int) *Grid`               |
|                            | `ShiftRectRight(r, c, h, w int) *Grid`          |
|                            | `ShiftRectLeft(r, c, h, w int) *Grid`           |
|                            | `ShiftRectUp(r, c, h, w int) *Grid`             |
|                            | `ShiftRectDown(r, c, h, w int) *Grid`           |
| **Print** (1)              | `Print() string`                                |

## License

MIT. See `LICENSE`.
