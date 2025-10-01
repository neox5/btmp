package btmp

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
	validateNonNegative(start, "start")
	validateNonNegative(count, "count")
	validateRangeOverflow(start, count)
	if start+count > b.lenBits {
		panic("range exceeds bitmap bounds")
	}
}
