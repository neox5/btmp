package btmp

// validateNonNegative validates that value is non-negative.
// Panics if value < 0.
func validateNonNegative(value int, name string) {
	if value < 0 {
		panic(name + " must be non-negative")
	}
}

// validatePositive validates that value is positive (> 0).
// Panics if value <= 0.
func validatePositive(value int, name string) {
	if value <= 0 {
		panic(name + " must be positive")
	}
}

// validateNotNil validates that pointer is not nil.
// Panics if ptr is nil.
func validateNotNil(ptr any, name string) {
	if ptr == nil {
		panic(name + " must not be nil")
	}
}

// validateRangeOverflow validates that start + count doesn't overflow.
// Panics if start + count < start.
func validateRangeOverflow(start, count int) {
	if start+count < start {
		panic("start + count overflow")
	}
}

// validateGridSizeMax validates that rows * cols doesn't overflow.
// Panics if rows * cols < 0.
func validateGridSizeMax(rows, cols int) {
	size := rows * cols
	if size < 0 {
		panic("grid size overflow")
	}
}

// validateWordBits validates that n is within word bit range for internal operations.
// Panics if n <= 0 or n > WordBits (64).
func validateWordBits(n int) {
	if n <= 0 || n > WordBits {
		panic("bit count must be > 0 and <= 64")
	}
}

// validateSameLength validates that two bitmaps have identical length.
// Panics if lengths differ.
func validateSameLength(a, b *Bitmap) {
	if a.Len() != b.Len() {
		panic("bitmaps must have same length")
	}
}
