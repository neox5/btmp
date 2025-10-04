package btmp

import "fmt"

// ValidationError represents a validation failure with context about what failed.
type ValidationError struct {
	Field   string // Name of the parameter that failed validation
	Value   any    // The actual value that failed (for debugging)
	Message string // Description of the validation failure
	Context string // Optional context (e.g., "Grid.SetRect", "Bitmap.CopyRange")
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Context != "" {
		return fmt.Sprintf("%s: %s: %s (got %v)", e.Context, e.Field, e.Message, e.Value)
	}
	return fmt.Sprintf("%s: %s (got %v)", e.Field, e.Message, e.Value)
}

// WithContext adds context to the validation error.
func (e *ValidationError) WithContext(ctx string) *ValidationError {
	e.Context = ctx
	return e
}

// validateNonNegative validates that value is non-negative.
// Returns ValidationError if value < 0.
func validateNonNegative(value int, name string) error {
	if value < 0 {
		return &ValidationError{
			Field:   name,
			Value:   value,
			Message: "must be non-negative",
		}
	}
	return nil
}

// validatePositive validates that value is positive (> 0).
// Returns ValidationError if value <= 0.
func validatePositive(value int, name string) error {
	if value <= 0 {
		return &ValidationError{
			Field:   name,
			Value:   value,
			Message: "must be positive",
		}
	}
	return nil
}

// validateNotNil validates that pointer is not nil.
// Returns ValidationError if ptr is nil.
func validateNotNil(ptr any, name string) error {
	if ptr == nil {
		return &ValidationError{
			Field:   name,
			Value:   nil,
			Message: "must not be nil",
		}
	}
	return nil
}

// validateRangeOverflow validates that start + count doesn't overflow.
// Returns ValidationError if start + count < start.
func validateRangeOverflow(start, count int) error {
	if start+count < start {
		return &ValidationError{
			Field:   "range",
			Value:   fmt.Sprintf("start=%d, count=%d", start, count),
			Message: "overflow",
		}
	}
	return nil
}

// validateGridSizeMax validates that rows * cols doesn't overflow.
// Returns ValidationError if rows * cols < 0.
func validateGridSizeMax(rows, cols int) error {
	size := rows * cols
	if size < 0 {
		return &ValidationError{
			Field:   "size",
			Value:   fmt.Sprintf("rows=%d, cols=%d", rows, cols),
			Message: "overflow",
		}
	}
	return nil
}

// validateWordBits validates that n is within word bit range for internal operations.
// Returns ValidationError if n <= 0 or n > WordBits (64).
func validateWordBits(n int) error {
	if n <= 0 || n > WordBits {
		return &ValidationError{
			Field:   "n",
			Value:   n,
			Message: fmt.Sprintf("must be > 0 and <= %d", WordBits),
		}
	}
	return nil
}

// validateSameLength validates that two bitmaps have identical length.
// Returns ValidationError if lengths differ.
func validateSameLength(a, b *Bitmap) error {
	if a.Len() != b.Len() {
		return &ValidationError{
			Field:   "length",
			Value:   fmt.Sprintf("a=%d, b=%d", a.Len(), b.Len()),
			Message: "bitmaps must have same length",
		}
	}
	return nil
}
