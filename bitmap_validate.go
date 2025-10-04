package btmp

import "fmt"

// validateInBounds validates that position is within bitmap bounds.
// Returns ValidationError if pos >= bitmap length.
func (b *Bitmap) validateInBounds(pos int) error {
	if pos >= b.lenBits {
		return &ValidationError{
			Field:   "pos",
			Value:   fmt.Sprintf("pos=%d, len=%d", pos, b.lenBits),
			Message: "position out of bounds",
		}
	}
	return nil
}

// validateRange validates a complete range operation against bitmap bounds.
// Validates start >= 0, count >= 0, no overflow, and range within bounds.
// Returns ValidationError on any validation failure.
func (b *Bitmap) validateRange(start, count int) error {
	if err := validateNonNegative(start, "start"); err != nil {
		return err
	}
	if err := validateNonNegative(count, "count"); err != nil {
		return err
	}
	if err := validateRangeOverflow(start, count); err != nil {
		return err
	}
	if start+count > b.lenBits {
		return &ValidationError{
			Field:   "range",
			Value:   fmt.Sprintf("start=%d, count=%d, len=%d", start, count, b.lenBits),
			Message: "exceeds bitmap bounds",
		}
	}
	return nil
}
