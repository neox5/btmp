# btmp

Compact, growable bitmap for Go. Fast range updates, overlap-safe copies, and a zero-copy row-major grid view.

**Module:** `github.com/neox5/btmp`

## Install
```bash
go get github.com/neox5/btmp
```

## Quick start

```go
import "github.com/neox5/btmp"

b := btmp.New(0).
    EnsureBits(8192).
    SetRange(100, 32).
    ClearRange(110, 4).
    CopyRange(b, 0, 256, 64).
    SetBits(500, 8, 0xFF).
    SetAll()

b2 := btmp.New(8192).SetRange(200, 50)
b.And(b2).Or(b2).Xor(b2).Not()

g := btmp.NewGridWithSize(16, 10).
    SetRect(3, 2, 5, 4).
    GrowCols(8).
    GrowRows(10)

idx := g.Index(7, 3)
bits := b.GetBits(100, 16)
_ = idx
_ = bits

// Print bits for inspection
fmt.Println(b.Print())                              // Binary string
fmt.Println(b.PrintFormat(16, false, 0, ""))        // Hexadecimal
fmt.Println(b.PrintFormat(2, true, 8, "_"))         // Binary grouped by 8
fmt.Println(b.PrintRangeFormat(0, 64, 16, true, 4, " "))  // Hex grouped
fmt.Println(g.Print())                              // Grid visualization
```

## API surface (V1)

**Bitmap**

Constructor:
* `New(n uint) *Bitmap`

Accessors:
* `(*Bitmap) Len() int`
* `(*Bitmap) Words() []uint64`

Growth operations:
* `(*Bitmap) EnsureBits(n int) *Bitmap`
* `(*Bitmap) AddBits(n int) *Bitmap`

Query operations:
* `(*Bitmap) Test(pos int) bool`
* `(*Bitmap) GetBits(pos, n int) uint64`
* `(*Bitmap) Any() bool`
* `(*Bitmap) Count() int`

Single-bit mutators:
* `(*Bitmap) SetBit(pos int) *Bitmap`
* `(*Bitmap) ClearBit(pos int) *Bitmap`
* `(*Bitmap) FlipBit(pos int) *Bitmap`

Multi-bit mutators:
* `(*Bitmap) SetBits(pos, n int, val uint64) *Bitmap`

Range mutators:
* `(*Bitmap) SetRange(start, count int) *Bitmap`
* `(*Bitmap) ClearRange(start, count int) *Bitmap`
* `(*Bitmap) CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap`
* `(*Bitmap) MoveRange(srcStart, dstStart, count int) *Bitmap`

Bulk mutators:
* `(*Bitmap) SetAll() *Bitmap`
* `(*Bitmap) ClearAll() *Bitmap`

Logical operations:
* `(*Bitmap) And(other *Bitmap) *Bitmap`
* `(*Bitmap) Or(other *Bitmap) *Bitmap`
* `(*Bitmap) Xor(other *Bitmap) *Bitmap`
* `(*Bitmap) Not() *Bitmap`

Print operations:
* `(*Bitmap) Print() string`
* `(*Bitmap) PrintRange(start, count int) string`
* `(*Bitmap) PrintFormat(base int, grouped bool, groupSize int, sep string) string`
* `(*Bitmap) PrintRangeFormat(start, count int, base int, grouped bool, groupSize int, sep string) string`

**Grid**

Constructors:
* `NewGrid() *Grid`
* `NewGridWithSize(cols, rows int) *Grid`

Accessors:
* `(*Grid) Cols() int`
* `(*Grid) Rows() int`
* `(*Grid) Index(x, y int) int`

Growth operations:
* `(*Grid) EnsureCols(cols int) *Grid`
* `(*Grid) EnsureRows(rows int) *Grid`
* `(*Grid) GrowCols(delta int) *Grid`
* `(*Grid) GrowRows(delta int) *Grid`

Query operations:
* `(*Grid) IsFree(x, y, w, h int) bool`
* `(*Grid) CanShiftRight(x, y, w, h int) bool`
* `(*Grid) CanShiftLeft(x, y, w, h int) bool`
* `(*Grid) CanShiftUp(x, y, w, h int) bool`
* `(*Grid) CanShiftDown(x, y, w, h int) bool`

Rectangle mutators:
* `(*Grid) SetRect(x, y, w, h int) *Grid`
* `(*Grid) ClearRect(x, y, w, h int) *Grid`
* `(*Grid) ShiftRectRight(x, y, w, h int) *Grid`
* `(*Grid) ShiftRectLeft(x, y, w, h int) *Grid`
* `(*Grid) ShiftRectUp(x, y, w, h int) *Grid`
* `(*Grid) ShiftRectDown(x, y, w, h int) *Grid`

Print operations:
* `(*Grid) Print() string`

## License

MIT. See `LICENSE`.
