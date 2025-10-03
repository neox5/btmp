// Package btmp provides a compact, growable bitmap optimized for fast range
// updates and overlap-safe copies, plus a zero-copy row-major 2D Grid view.
//
// Conventions:
//   - Length is in bits (Len).
//   - Storage is []uint64 words, exposed via Words() for read-only inspection.
//   - Ranges use (start, count).
//   - All operations are in-bounds only - no auto-growth.
//   - All mutating methods return *Bitmap for chaining.
//
// Invariant:
//   - After any public mutator returns, all bits at indexes >= Len() are zero,
//     even when count == 0.
package btmp

import "math/bits"

const (
	WordBits         = 64
	WordShift        = 6            // log2(64), divide by 64 via >> 6
	IndexMask        = WordBits - 1 // for i & IndexMask
	WordMask  uint64 = ^uint64(0)   // 0xFFFFFFFFFFFFFFFF
)

// Bitmap is a growable bitset backed by 64-bit words.
type Bitmap struct {
	words       []uint64
	lenBits     int
	lastWordIdx int    // index of last logical word; -1 if Len()==0
	tailMask    uint64 // mask for last logical word; 0 if Len()==0; WordMask if Len()%64==0
}

// ========================================
// Constructor Functions
// ========================================

// New returns an empty bitmap sized for n bits (Len==n).
func New(n uint) *Bitmap {
	b := &Bitmap{
		words:   make([]uint64, (n+IndexMask)>>WordShift),
		lenBits: int(n),
	}
	b.computeCache()
	return b
}

// ========================================
// Accessors
// ========================================

// Len returns the logical length in bits.
func (b *Bitmap) Len() int { return b.lenBits }

// Words exposes the underlying words slice (length may exceed the logical need).
func (b *Bitmap) Words() []uint64 { return b.words }

// ========================================
// Growth Operations
// ========================================

// EnsureBits grows the logical length to at least n bits. No-op if n <= Len().
// Returns *Bitmap for chaining. Panics if n < 0.
func (b *Bitmap) EnsureBits(n int) *Bitmap {
	validateNonNegative(n, "n")

	if n > b.lenBits {
		b.ensureBits(n)
		b.computeCache()
	}
	return b
}

// AddBits grows the logical length by n bits.
// Returns *Bitmap for chaining. Panics if n < 0.
func (b *Bitmap) AddBits(n int) *Bitmap {
	validateNonNegative(n, "n")

	if n > 0 {
		b.addBits(n)
		b.computeCache()
	}
	return b
}

// ========================================
// Query Operations
// ========================================

// Test reports whether bit pos is set. Panics if pos is out of [0, Len()).
func (b *Bitmap) Test(pos int) bool {
	validateNonNegative(pos, "pos")
	b.validateInBounds(pos)

	w, off := wordIndex(pos)
	return (b.words[w]>>off)&1 == 1
}

// GetBits extracts n bits starting from pos, returned right-aligned.
// The result contains the extracted bits in the least significant positions.
// Panics if pos < 0, n <= 0, n > 64, or pos+n > Len().
//
// Example: bitmap 11010110, GetBits(2, 3) returns 101 (bits at positions 2,3,4).
func (b *Bitmap) GetBits(pos, n int) uint64 {
	validateNonNegative(pos, "pos")
	validateWordBits(n)
	b.validateRange(pos, n)

	return b.getBits(pos, n)
}

// Any reports whether any bit in [0, Len()) is set.
func (b *Bitmap) Any() bool {
	if b.lenBits == 0 {
		return false
	}
	// full words except the last
	for i := range b.lastWordIdx {
		if b.words[i] != 0 {
			return true
		}
	}
	// masked last word
	return (b.words[b.lastWordIdx] & b.tailMask) != 0
}

// Count returns the number of set bits in [0, Len()).
func (b *Bitmap) Count() int {
	if b.lenBits == 0 {
		return 0
	}
	sum := 0
	for i := range b.lastWordIdx {
		sum += bits.OnesCount64(b.words[i])
	}
	return sum + bits.OnesCount64(b.words[b.lastWordIdx]&b.tailMask)
}

// ========================================
// Single-Bit Mutators
// ========================================

// SetBit sets bit pos to 1. Panics if pos < 0 or pos >= Len().
func (b *Bitmap) SetBit(pos int) *Bitmap {
	validateNonNegative(pos, "pos")
	b.validateInBounds(pos)

	b.setBit(pos)
	return b
}

// ClearBit sets bit pos to 0. Panics if pos < 0 or pos >= Len().
func (b *Bitmap) ClearBit(pos int) *Bitmap {
	validateNonNegative(pos, "pos")
	b.validateInBounds(pos)

	b.clearBit(pos)
	return b
}

// FlipBit toggles bit pos. Panics if pos < 0 or pos >= Len().
func (b *Bitmap) FlipBit(pos int) *Bitmap {
	validateNonNegative(pos, "pos")
	b.validateInBounds(pos)

	b.flipBit(pos)
	return b
}

// ========================================
// Multi-Bit Mutators
// ========================================

// SetBits inserts the low n bits of val into the bitmap starting at pos.
// Only the least significant n bits of val are used; higher bits are ignored.
// Preserves surrounding bits unchanged. Panics if pos < 0, n <= 0, n > 64, or pos+n > Len().
// Returns *Bitmap for chaining.
//
// Example: SetBits(2, 3, 0b101) sets 3 bits starting at position 2 to the pattern 101.
func (b *Bitmap) SetBits(pos, n int, val uint64) *Bitmap {
	validateNonNegative(pos, "pos")
	validateWordBits(n)
	b.validateRange(pos, n)

	b.setBits(pos, n, val)
	return b
}

// ========================================
// Range Mutators
// ========================================

// SetRange sets bits in [start, start+count) to 1. In-bounds only.
// Returns *Bitmap for chaining. Panics on negative inputs, overflow, or out-of-bounds.
func (b *Bitmap) SetRange(start, count int) *Bitmap {
	b.validateRange(start, count)

	b.setRange(start, count)
	return b
}

// ClearRange clears bits in [start, start+count) to 0. In-bounds only.
// Returns *Bitmap for chaining. Panics on negative inputs, overflow, or out-of-bounds.
func (b *Bitmap) ClearRange(start, count int) *Bitmap {
	b.validateRange(start, count)

	b.clearRange(start, count)
	return b
}

// CopyRange copies count bits from src[srcStart:] to dst[dstStart:].
// In-bounds only for both src and dst. Overlap-safe with memmove semantics.
// Returns *Bitmap for chaining. Panics on negative inputs, nil src, or out-of-bounds.
func (b *Bitmap) CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap {
	validateNotNil(src, "src")

	src.validateRange(srcStart, count)
	b.validateRange(dstStart, count)

	b.copyRange(src, srcStart, dstStart, count)
	return b
}

// MoveRange moves count bits from [srcStart, srcStart+count) to [dstStart, dstStart+count).
// The source range is cleared after copying. Overlap-safe with memmove semantics.
// In-bounds only for both source and destination ranges.
// Returns *Bitmap for chaining. Panics on negative inputs, overflow, or out-of-bounds.
func (b *Bitmap) MoveRange(srcStart, dstStart, count int) *Bitmap {
	b.validateRange(srcStart, count)
	b.validateRange(dstStart, count)

	b.moveRange(srcStart, dstStart, count)
	return b
}

// ========================================
// Bulk Mutators
// ========================================

// SetAll sets all bits in [0, Len()) to 1.
// Equivalent to SetRange(0, Len()) but optimized for full bitmap operations.
// Returns *Bitmap for chaining.
func (b *Bitmap) SetAll() *Bitmap {
	b.setAll()
	return b
}

// ClearAll clears all bits in [0, Len()) to 0.
// Equivalent to ClearRange(0, Len()) but optimized for full bitmap operations.
// Returns *Bitmap for chaining.
func (b *Bitmap) ClearAll() *Bitmap {
	b.clearAll()
	return b
}

// ========================================
// Logical Operations
// ========================================

// And performs bitwise AND with other bitmap. Both bitmaps must have the same length.
// Returns *Bitmap for chaining. Panics if other is nil or lengths differ.
func (b *Bitmap) And(other *Bitmap) *Bitmap {
	validateNotNil(other, "other")
	validateSameLength(b, other)

	b.and(other)
	return b
}

// Or performs bitwise OR with other bitmap. Both bitmaps must have the same length.
// Returns *Bitmap for chaining. Panics if other is nil or lengths differ.
func (b *Bitmap) Or(other *Bitmap) *Bitmap {
	validateNotNil(other, "other")
	validateSameLength(b, other)

	b.or(other)
	return b
}

// Xor performs bitwise XOR with other bitmap. Both bitmaps must have the same length.
// Returns *Bitmap for chaining. Panics if other is nil or lengths differ.
func (b *Bitmap) Xor(other *Bitmap) *Bitmap {
	validateNotNil(other, "other")
	validateSameLength(b, other)

	b.xor(other)
	return b
}

// Not performs bitwise NOT, flipping all bits in [0, Len()).
// Returns *Bitmap for chaining.
func (b *Bitmap) Not() *Bitmap {
	b.not()
	return b
}

// ========================================
// Print Operations
// ========================================

// Print formats all bits in [0, Len()) as binary string.
// Returns empty string if Len() == 0.
func (b *Bitmap) Print() string {
	return b.PrintRange(0, b.lenBits)
}

// PrintRange formats bits in [start, start+count) as binary string.
// Returns empty string if count == 0.
// Panics if start < 0, count < 0, or start+count > Len().
func (b *Bitmap) PrintRange(start, count int) string {
	b.validateRange(start, count)
	return b.printRangeFormat(start, count, 2, false, 0, "")
}

// PrintFormat formats all bits according to format parameters.
// base: 2 (binary) or 16 (hexadecimal)
// grouped: insert separators between bit groups
// groupSize: units per group (bits for base 2, hex digits for base 16)
// sep: separator string
// Panics if base not in {2,16} or grouped && groupSize <= 0.
func (b *Bitmap) PrintFormat(base int, grouped bool, groupSize int, sep string) string {
	return b.PrintRangeFormat(0, b.lenBits, base, grouped, groupSize, sep)
}

// PrintRangeFormat formats bits in [start, start+count) with format parameters.
// base: 2 (binary) or 16 (hexadecimal)
// grouped: insert separators between bit groups
// groupSize: units per group (bits for base 2, hex digits for base 16)
// sep: separator string
// Panics if start < 0, count < 0, start+count > Len(), base not in {2,16},
// or grouped && groupSize <= 0.
func (b *Bitmap) PrintRangeFormat(start, count int, base int, grouped bool, groupSize int, sep string) string {
	b.validateRange(start, count)

	if base != 2 && base != 16 {
		panic("base must be 2 or 16")
	}
	if grouped && groupSize <= 0 {
		panic("groupSize must be positive when grouped")
	}

	return b.printRangeFormat(start, count, base, grouped, groupSize, sep)
}

// ========================================
// Internal Helpers
// ========================================

// computeCache recomputes cache fields from lenBits only.
func (b *Bitmap) computeCache() {
	if b.lenBits == 0 {
		b.lastWordIdx = -1
		b.tailMask = 0
		return
	}
	// ceil(lenBits/64) - 1
	b.lastWordIdx = int((b.lenBits+IndexMask)>>WordShift) - 1

	r := uint(b.lenBits) & IndexMask // bits used in last word, 0..63
	if r == 0 {
		b.tailMask = WordMask
		return
	}
	b.tailMask = MaskUpto(r)
}

// wordIndex converts a bit index to (wordIdx, bitOffset).
func wordIndex(i int) (w int, off int) {
	return i >> WordShift, i & IndexMask
}
