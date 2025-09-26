package btmp

// validatePosition validates that pos is non-negative.
// Panics if pos < 0.
func validatePosition(pos int) {
	if pos < 0 {
		panic("position must be non-negative")
	}
}

// validateCount validates that count is non-negative.
// Panics if count < 0.
func validateCount(count int) {
	if count < 0 {
		panic("count must be non-negative")
	}
}

// validateWordBits validates that n is within word bit range for internal operations.
// Panics if n <= 0 or n > WordBits (64).
func validateWordBits(n int) {
	if n <= 0 || n > WordBits {
		panic("bit count must be > 0 and <= 64")
	}
}

// validateNoOverflow validates that start + count doesn't overflow.
// Panics if start + count < start.
func validateNoOverflow(start, count int) {
	if start+count < start {
		panic("start + count overflow")
	}
}

// validateSameLength validates that two bitmaps have identical length.
// Panics if lengths differ.
func validateSameLength(a, b *Bitmap) {
	if a.Len() != b.Len() {
		panic("bitmaps must have same length")
	}
}

// validateInBounds validates that position is within bitmap bounds.
// Panics if pos >= bitmap length.
func (b *Bitmap) validateInBounds(pos int) {
	if pos >= b.lenBits {
		panic("position out of bounds")
	}
}

// validateRange validates a complete range operation against bitmap bounds.
// Validates start >= 0, count >= 0, no overflow, and range within bounds.
// Panics on any validation failure.
func (b *Bitmap) validateRange(start, count int) {
	validatePosition(start)
	validateCount(count)
	validateNoOverflow(start, count)
	if start+count > b.lenBits {
		panic("range exceeds bitmap bounds")
	}
}
