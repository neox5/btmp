// bitmap_validation.go
package btmp

// validateBitCount validates that n is within valid bit count range.
// Panics if n <= 0 or n > WordBits (64).
func validateBitCount(n int) {
	if n <= 0 || n > WordBits {
		panic("bit count must be > 0 and <= 64")
	}
}

// validatePosition validates that pos is non-negative.
// Panics if pos < 0.
func validatePosition(pos int) {
	if pos < 0 {
		panic("position must be non-negative")
	}
}

// validateBitOp performs complete validation for bit operations.
// Validates position, bit count, integer overflow, and bitmap bounds.
// Panics on any validation failure.
func (b *Bitmap) validateBitOp(pos, n int) {
	validatePosition(pos)
	validateBitCount(n)

	end := pos + n
	if end < pos { // integer overflow check
		panic("position + bit count overflow")
	}

	// Bitmap bounds check - operations must be within current length
	if end > b.lenBits {
		panic("operation exceeds bitmap bounds")
	}
}
