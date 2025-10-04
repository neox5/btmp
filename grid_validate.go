package btmp

import "fmt"

// validateCoordinate validates that x and y are non-negative and within grid bounds.
// Returns ValidationError if x < 0, y < 0, x >= g.Cols(), or y >= g.Rows().
func (g *Grid) validateCoordinate(x, y int) error {
	if err := validateNonNegative(x, "x"); err != nil {
		return err
	}
	if err := validateNonNegative(y, "y"); err != nil {
		return err
	}
	if x >= g.cols {
		return &ValidationError{
			Field:   "x",
			Value:   fmt.Sprintf("x=%d, cols=%d", x, g.cols),
			Message: "out of bounds",
		}
	}
	if y >= g.Rows() {
		return &ValidationError{
			Field:   "y",
			Value:   fmt.Sprintf("y=%d, rows=%d", y, g.Rows()),
			Message: "out of bounds",
		}
	}
	return nil
}

// validateRect validates that rectangle parameters are non-negative
// and rectangle is fully contained within grid bounds.
// Returns ValidationError if x < 0, y < 0, w < 0, h < 0, x+w > g.Cols(), or y+h > g.Rows().
func (g *Grid) validateRect(x, y, w, h int) error {
	if err := g.validateCoordinate(x, y); err != nil {
		return err
	}
	if err := validatePositive(w, "w"); err != nil {
		return err
	}
	if err := validatePositive(h, "h"); err != nil {
		return err
	}
	if x+w > g.cols {
		return &ValidationError{
			Field:   "rectangle",
			Value:   fmt.Sprintf("x=%d, w=%d, cols=%d", x, w, g.cols),
			Message: "exceeds columns",
		}
	}
	if y+h > g.Rows() {
		return &ValidationError{
			Field:   "rectangle",
			Value:   fmt.Sprintf("y=%d, h=%d, rows=%d", y, h, g.Rows()),
			Message: "exceeds rows",
		}
	}
	return nil
}
