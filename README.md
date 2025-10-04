<br/>
<br/>

<div align="center">
  <img src="logo.png" alt="btmp" width="500"/>
</div>

<br/>
<br/>

# btmp

btmp ("bitmap") is a pure Go bitmap library designed as a building block for your data structures. It provides tested, validated operations for manipulating dense boolean data without implementing bit math yourself.

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
count := b.Count()               // Number of set bits
bits := b.GetBits(100, 16)       // Extract 16 bits starting at 100

// Grid - zero-copy 2D view
g := btmp.NewGridWithSize(16, 10)
g.SetRect(3, 2, 5, 4)            // Set 5×4 rectangle at (3,2)
if g.IsFree(8, 5, 3, 3) {        // Check if 3×3 region is available
    g.SetRect(8, 5, 3, 3)
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

**Bitmap** - 34 methods
- Construction (1): `New`
- Access (2): `Len`, `Words`
- Growth (2): `EnsureBits`, `AddBits`
- Query (4): `Test`, `GetBits`, `Any`, `Count`
- Validation (2): `ValidateInBounds`, `ValidateRange`
- Single-bit (3): `SetBit`, `ClearBit`, `FlipBit`
- Multi-bit (1): `SetBits`
- Range (4): `SetRange`, `ClearRange`, `CopyRange`, `MoveRange`
- Bulk (2): `SetAll`, `ClearAll`
- Logic (4): `And`, `Or`, `Xor`, `Not`
- Print (4): `Print`, `PrintRange`, `PrintFormat`, `PrintRangeFormat`

**Grid** - 20 methods
- Construction (2): `NewGrid`, `NewGridWithSize`
- Access (3): `Cols`, `Rows`, `Index`
- Growth (4): `EnsureCols`, `EnsureRows`, `GrowCols`, `GrowRows`
- Query (5): `IsFree`, `CanShiftRight`, `CanShiftLeft`, `CanShiftUp`, `CanShiftDown`
- Validation (2): `ValidateCoordinate`, `ValidateRect`
- Rectangle (6): `SetRect`, `ClearRect`, `ShiftRect{Right,Left,Up,Down}`
- Print (1): `Print`

## License

MIT. See `LICENSE`.
