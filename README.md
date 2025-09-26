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

g := btmp.NewGrid(16).
    SetRect(3, 2, 5, 4).
    GrowCols(8).
    GrowRows(10)

idx := g.Index(7, 3)
bits := b.GetBits(100, 16)
_ = idx
_ = bits
```

## API surface (V1)

**Bitmap**

Constructor:
* `New(n uint) *Bitmap`

Accessors:
* `(*Bitmap) Len() int`
* `(*Bitmap) Words() []uint64`

Growth mutators:
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

Bulk mutators:
* `(*Bitmap) SetAll() *Bitmap`
* `(*Bitmap) ClearAll() *Bitmap`

Logical operations:
* `(*Bitmap) And(other *Bitmap) *Bitmap`
* `(*Bitmap) Or(other *Bitmap) *Bitmap`
* `(*Bitmap) Xor(other *Bitmap) *Bitmap`
* `(*Bitmap) Not() *Bitmap`

**Grid**

* `NewGrid(cols int) *Grid`
* `NewGridWithSize(cols, rows int) *Grid`
* `NewGridFrom(b *Bitmap, cols int) *Grid`
* `(*Grid) Index(x, y int) int`
* `(*Grid) Cols() int`
* `(*Grid) Rows() int`
* `(*Grid) SetRect(x, y, w, h int) *Grid`
* `(*Grid) ClearRect(x, y, w, h int) *Grid`
* `(*Grid) GrowCols(delta int) *Grid`
* `(*Grid) EnsureCols(cols int) *Grid`
* `(*Grid) GrowRows(delta int) *Grid`
* `(*Grid) EnsureRows(rows int) *Grid`

## License

MIT. See `LICENSE`.
