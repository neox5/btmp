package btmp

import "fmt"

// validateCoordinate validates that c and r are non-negative and within grid bounds.
// Returns ValidationError if c < 0, r < 0, c >= g.Cols(), or r >= g.Rows().
func (g *Grid) validateCoordinate(c, r int) error {
	if err := validateNonNegative(c, "c"); err != nil {
		return err
	}
	if err := validateNonNegative(r, "r"); err != nil {
		return err
	}
	if c >= g.cols {
		return &ValidationError{
			Field:   "c",
			Value:   fmt.Sprintf("c=%d, cols=%d", c, g.cols),
			Message: "out of bounds",
		}
	}
	if r >= g.Rows() {
		return &ValidationError{
			Field:   "r",
			Value:   fmt.Sprintf("r=%d, rows=%d", r, g.Rows()),
			Message: "out of bounds",
		}
	}
	return nil
}

// validateRect validates that rectangle parameters are non-negative
// and rectangle is fully contained within grid bounds.
// Returns ValidationError if c < 0, r < 0, w < 0, h < 0, c+w > g.Cols(), or r+h > g.Rows().
func (g *Grid) validateRect(c, r, w, h int) error {
	if err := g.validateCoordinate(c, r); err != nil {
		return err
	}
	if err := validatePositive(w, "w"); err != nil {
		return err
	}
	if err := validatePositive(h, "h"); err != nil {
		return err
	}
	if c+w > g.cols {
		return &ValidationError{
			Field:   "rectangle",
			Value:   fmt.Sprintf("c=%d, w=%d, cols=%d", c, w, g.cols),
			Message: "exceeds columns",
		}
	}
	if r+h > g.Rows() {
		return &ValidationError{
			Field:   "rectangle",
			Value:   fmt.Sprintf("r=%d, h=%d, rows=%d", r, h, g.Rows()),
			Message: "exceeds rows",
		}
	}
	return nil
}
