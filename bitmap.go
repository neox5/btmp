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
	if err := validateNonNegative(n, "n"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.EnsureBits"))
	}

	if n > b.lenBits {
		b.ensureBits(n)
		b.computeCache()
	}
	return b
}

// AddBits grows the logical length by n bits.
// Returns *Bitmap for chaining. Panics if n < 0.
func (b *Bitmap) AddBits(n int) *Bitmap {
	if err := validateNonNegative(n, "n"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.AddBits"))
	}

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
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.Test"))
	}
	if err := b.validateInBounds(pos); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.Test"))
	}

	return b.test(pos)
}

// Any reports whether any bit in [0, Len()) is set.
func (b *Bitmap) Any() bool {
	return b.any()
}

// All reports whether all bits in [0, Len()) are set.
// Returns true for empty bitmaps (vacuously true).
func (b *Bitmap) All() bool {
	return b.all()
}

// Count returns the number of set bits in [0, Len()).
func (b *Bitmap) Count() int {
	return b.count()
}

// AnyRange reports whether any bit in [start, start+count) is set.
// Returns false for empty ranges (count == 0).
// Panics if start < 0, count < 0, or start+count > Len().
func (b *Bitmap) AnyRange(start, count int) bool {
	if err := b.validateRange(start, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.AnyRange"))
	}

	return b.anyRange(start, count)
}

// AllRange reports whether all bits in [start, start+count) are set.
// Returns true for empty ranges (vacuously true).
// Panics if start < 0, count < 0, or start+count > Len().
func (b *Bitmap) AllRange(start, count int) bool {
	if err := b.validateRange(start, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.AllRange"))
	}

	return b.allRange(start, count)
}

// CountRange returns the number of set bits in [start, start+count).
// Returns 0 for empty ranges (count == 0).
// Panics if start < 0, count < 0, or start+count > Len().
func (b *Bitmap) CountRange(start, count int) int {
	if err := b.validateRange(start, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CountRange"))
	}

	return b.countRange(start, count)
}

// NextZero returns the position of the next zero bit at or after pos.
// Returns -1 if no zero bit exists in [pos, Len()).
// Panics if pos < 0 or pos >= Len().
func (b *Bitmap) NextZero(pos int) int {
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.NextZero"))
	}
	if err := b.validateInBounds(pos); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.NextZero"))
	}

	return b.nextZero(pos)
}

// NextOne returns the position of the next set bit at or after pos.
// Returns -1 if no set bit exists in [pos, Len()).
// Panics if pos < 0 or pos >= Len().
func (b *Bitmap) NextOne(pos int) int {
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.NextOne"))
	}
	if err := b.validateInBounds(pos); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.NextOne"))
	}

	return b.nextOne(pos)
}

// NextZeroInRange returns the position of the next zero bit in [pos, pos+count).
// Returns -1 if no zero bit exists in range.
// Panics if pos < 0, count <= 0, or pos+count > Len().
func (b *Bitmap) NextZeroInRange(pos, count int) int {
	if err := b.validateRange(pos, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.NextZeroInRange"))
	}

	return b.nextZeroInRange(pos, count)
}

// NextOneInRange returns the position of the next set bit in [pos, pos+count).
// Returns -1 if no set bit exists in range.
// Panics if pos < 0, count <= 0, or pos+count > Len().
func (b *Bitmap) NextOneInRange(pos, count int) int {
	if err := b.validateRange(pos, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.NextOneInRange"))
	}

	return b.nextOneInRange(pos, count)
}

// CountZerosFrom counts consecutive zero bits starting at pos.
// Returns 0 if bit at pos is set.
// Stops at first set bit or end of bitmap.
// Panics if pos < 0 or pos >= Len().
func (b *Bitmap) CountZerosFrom(pos int) int {
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CountZerosFrom"))
	}
	if err := b.validateInBounds(pos); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CountZerosFrom"))
	}

	return b.countZerosFrom(pos)
}

// CountOnesFrom counts consecutive set bits starting at pos.
// Returns 0 if bit at pos is clear.
// Stops at first zero bit or end of bitmap.
// Panics if pos < 0 or pos >= Len().
func (b *Bitmap) CountOnesFrom(pos int) int {
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CountOnesFrom"))
	}
	if err := b.validateInBounds(pos); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CountOnesFrom"))
	}

	return b.countOnesFrom(pos)
}

// CountZerosFromInRange counts consecutive zero bits starting at pos within [pos, pos+count).
// Returns 0 if bit at pos is set.
// Stops at first set bit or end of range.
// Panics if pos < 0, count <= 0, or pos+count > Len().
func (b *Bitmap) CountZerosFromInRange(pos, count int) int {
	if err := b.validateRange(pos, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CountZerosFromInRange"))
	}

	return b.countZerosFromInRange(pos, count)
}

// CountOnesFromInRange counts consecutive set bits starting at pos within [pos, pos+count).
// Returns 0 if bit at pos is clear.
// Stops at first zero bit or end of range.
// Panics if pos < 0, count <= 0, or pos+count > Len().
func (b *Bitmap) CountOnesFromInRange(pos, count int) int {
	if err := b.validateRange(pos, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CountOnesFromInRange"))
	}

	return b.countOnesFromInRange(pos, count)
}

// ========================================
// Validation Operations
// ========================================

// ValidateInBounds validates that position is within bitmap bounds.
// Returns ValidationError if pos >= bitmap length.
func (b *Bitmap) ValidateInBounds(pos int) error {
	return b.validateInBounds(pos)
}

// ValidateRange validates a complete range operation against bitmap bounds.
// Validates start >= 0, count >= 0, no overflow, and range within bounds.
// Returns ValidationError on any validation failure.
func (b *Bitmap) ValidateRange(start, count int) error {
	return b.validateRange(start, count)
}

// ========================================
// Single-Bit Mutators
// ========================================

// SetBit sets bit pos to 1. Panics if pos < 0 or pos >= Len().
func (b *Bitmap) SetBit(pos int) *Bitmap {
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.SetBit"))
	}
	if err := b.validateInBounds(pos); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.SetBit"))
	}

	b.setBit(pos)
	return b
}

// ClearBit sets bit pos to 0. Panics if pos < 0 or pos >= Len().
func (b *Bitmap) ClearBit(pos int) *Bitmap {
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.ClearBit"))
	}
	if err := b.validateInBounds(pos); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.ClearBit"))
	}

	b.clearBit(pos)
	return b
}

// FlipBit toggles bit pos. Panics if pos < 0 or pos >= Len().
func (b *Bitmap) FlipBit(pos int) *Bitmap {
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.FlipBit"))
	}
	if err := b.validateInBounds(pos); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.FlipBit"))
	}

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
// This method is primarily useful for initializing bitmaps from constants:
//
//	b.SetBits(0, 16, 0xABCD)  // Set hex pattern
//	b.SetBits(8, 4, 0b1010)   // Set binary pattern
func (b *Bitmap) SetBits(pos, n int, val uint64) *Bitmap {
	if err := validateNonNegative(pos, "pos"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.SetBits"))
	}
	if err := validateWordBits(n); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.SetBits"))
	}
	if err := b.validateRange(pos, n); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.SetBits"))
	}

	b.setBits(pos, n, val)
	return b
}

// ========================================
// Range Mutators
// ========================================

// SetRange sets bits in [start, start+count) to 1. In-bounds only.
// Returns *Bitmap for chaining. Panics on negative inputs, overflow, or out-of-bounds.
func (b *Bitmap) SetRange(start, count int) *Bitmap {
	if err := b.validateRange(start, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.SetRange"))
	}

	b.setRange(start, count)
	return b
}

// ClearRange clears bits in [start, start+count) to 0. In-bounds only.
// Returns *Bitmap for chaining. Panics on negative inputs, overflow, or out-of-bounds.
func (b *Bitmap) ClearRange(start, count int) *Bitmap {
	if err := b.validateRange(start, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.ClearRange"))
	}

	b.clearRange(start, count)
	return b
}

// CopyRange copies count bits from src[srcStart:] to dst[dstStart:].
// In-bounds only for both src and dst. Overlap-safe with memmove semantics.
// Returns *Bitmap for chaining. Panics on negative inputs, nil src, or out-of-bounds.
func (b *Bitmap) CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap {
	if err := validateNotNil(src, "src"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CopyRange"))
	}
	if err := src.validateRange(srcStart, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CopyRange"))
	}
	if err := b.validateRange(dstStart, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.CopyRange"))
	}

	b.copyRange(src, srcStart, dstStart, count)
	return b
}

// MoveRange moves count bits from [srcStart, srcStart+count) to [dstStart, dstStart+count).
// The source range is cleared after copying. Overlap-safe with memmove semantics.
// In-bounds only for both source and destination ranges.
// Returns *Bitmap for chaining. Panics on negative inputs, overflow, or out-of-bounds.
func (b *Bitmap) MoveRange(srcStart, dstStart, count int) *Bitmap {
	if err := b.validateRange(srcStart, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.MoveRange"))
	}
	if err := b.validateRange(dstStart, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.MoveRange"))
	}

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
	if err := validateNotNil(other, "other"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.And"))
	}
	if err := validateSameLength(b, other); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.And"))
	}

	b.and(other)
	return b
}

// Or performs bitwise OR with other bitmap. Both bitmaps must have the same length.
// Returns *Bitmap for chaining. Panics if other is nil or lengths differ.
func (b *Bitmap) Or(other *Bitmap) *Bitmap {
	if err := validateNotNil(other, "other"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.Or"))
	}
	if err := validateSameLength(b, other); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.Or"))
	}

	b.or(other)
	return b
}

// Xor performs bitwise XOR with other bitmap. Both bitmaps must have the same length.
// Returns *Bitmap for chaining. Panics if other is nil or lengths differ.
func (b *Bitmap) Xor(other *Bitmap) *Bitmap {
	if err := validateNotNil(other, "other"); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.Xor"))
	}
	if err := validateSameLength(b, other); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.Xor"))
	}

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
	if err := b.validateRange(start, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.PrintRange"))
	}
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
	if err := b.validateRange(start, count); err != nil {
		panic(err.(*ValidationError).WithContext("Bitmap.PrintRangeFormat"))
	}

	if base != 2 && base != 16 {
		panic(&ValidationError{
			Field:   "base",
			Value:   base,
			Message: "must be 2 or 16",
			Context: "Bitmap.PrintRangeFormat",
		})
	}
	if grouped && groupSize <= 0 {
		panic(&ValidationError{
			Field:   "groupSize",
			Value:   groupSize,
			Message: "must be positive when grouped",
			Context: "Bitmap.PrintRangeFormat",
		})
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
