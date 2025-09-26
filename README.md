# btmp

Compact, growable bitmap for Go. Fast range updates, overlap-safe copies, and a zero-copy row-major grid view.

**Module:** `github.com/neox5/btmp`

## Install
```bash
go get github.com/neox5/btmp
````

## Quick start

```go
import "github.com/neox5/btmp"

b := btmp.New().
    EnsureBits(8192).
    SetRange(100, 32).
    ClearRange(110, 4).
    CopyRange(b, 0, 256, 64)

g := btmp.NewGrid(16).
    SetRect(3, 2, 5, 4).
    GrowCols(8).
    GrowRows(10)

idx := g.Index(7, 3)
_ = idx
```

## API surface (V1)

**Bitmap**

* `New() *Bitmap`
* `NewWithCap(capBits int) *Bitmap`
* `(*Bitmap) Len() int`
* `(*Bitmap) Words() []uint64`
* `(*Bitmap) Test(i int) bool`
* `(*Bitmap) Any() bool`
* `(*Bitmap) Count() int`
* `(*Bitmap) SetRange(start, count int) *Bitmap`
* `(*Bitmap) ClearRange(start, count int) *Bitmap`
* `(*Bitmap) CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap`
* `(*Bitmap) EnsureBits(n int) *Bitmap`
* `(*Bitmap) AddBits(n int) *Bitmap`

**Grid**

* `NewGrid(cols int) *Grid`
* `NewGridWithCap(cols, rowsCap int) *Grid`
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

## Semantics

* Length in bits. Storage `[]uint64`.
* Ranges use `(start, count)`.
* `SetRange`, `CopyRange` auto-grow. `ClearRange` in-bounds.
* All mutators return `*Bitmap` for chaining; pointer identity is stable.
* Grid maintains `Len() == Rows()*Cols`.

## License

MIT. See `LICENSE`.
