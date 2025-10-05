package btmp

import "fmt"

// validateCoordinate validates that r and c are non-negative and within grid bounds.
// Returns ValidationError if r < 0, c < 0, r >= g.Rows(), or c >= g.Cols().
func (g *Grid) validateCoordinate(r, c int) error {
	if err := validateNonNegative(r, "r"); err != nil {
		return err
	}
	if err := validateNonNegative(c, "c"); err != nil {
		return err
	}
	if r >= g.rows {
		return &ValidationError{
			Field:   "r",
			Value:   fmt.Sprintf("r=%d, rows=%d", r, g.rows),
			Message: "out of bounds",
		}
	}
	if c >= g.cols {
		return &ValidationError{
			Field:   "c",
			Value:   fmt.Sprintf("c=%d, cols=%d", c, g.cols),
			Message: "out of bounds",
		}
	}
	return nil
}

// validateRect validates that rectangle parameters are non-negative
// and rectangle is fully contained within grid bounds.
// Returns ValidationError if r < 0, c < 0, h < 0, w < 0, r+h > g.Rows(), or c+w > g.Cols().
func (g *Grid) validateRect(r, c, h, w int) error {
	if err := g.validateCoordinate(r, c); err != nil {
		return err
	}
	if err := validatePositive(h, "h"); err != nil {
		return err
	}
	if err := validatePositive(w, "w"); err != nil {
		return err
	}
	if r+h > g.rows {
		return &ValidationError{
			Field:   "rectangle",
			Value:   fmt.Sprintf("r=%d, h=%d, rows=%d", r, h, g.rows),
			Message: "exceeds rows",
		}
	}
	if c+w > g.cols {
		return &ValidationError{
			Field:   "rectangle",
			Value:   fmt.Sprintf("c=%d, w=%d, cols=%d", c, w, g.cols),
			Message: "exceeds columns",
		}
	}
	return nil
}
