package btmp

import (
	"fmt"
	"strings"
)

// print formats the grid as a coordinate-labeled visualization.
// Internal implementation - no validation.
func (g *Grid) print() string {
	rows := g.Rows()
	cols := g.cols

	// Empty grid
	if rows == 0 || cols == 0 {
		return ""
	}

	// Calculate widths for alignment
	rowWidth := len(fmt.Sprintf("%d", rows-1))
	colWidth := len(fmt.Sprintf("%d", cols-1))

	// Each column needs colWidth + 1 space separator (except last)
	cellWidth := colWidth + 1

	// Estimate capacity
	headerLen := rowWidth + 1 + (cols * cellWidth)
	lineLen := rowWidth + 1 + (cols * cellWidth) + 1 // +1 for newline
	capacity := headerLen + (rows * lineLen)
	var builder strings.Builder
	builder.Grow(capacity)

	// Build column header
	for range rowWidth {
		builder.WriteByte(' ')
	}
	builder.WriteByte(' ')
	for col := range cols {
		// Right-align column number within colWidth
		colStr := fmt.Sprintf("%*d", colWidth, col)
		builder.WriteString(colStr)
		// Space separator after each column except last
		if col < cols-1 {
			builder.WriteByte(' ')
		}
	}
	builder.WriteByte('\n')

	// Build grid by iterating through bitmap words
	bitLen := g.B.Len()
	words := g.B.Words()
	bitIdx := 0
	col := 0
	row := 0

	for _, word := range words {
		for bitOff := range WordBits {
			if bitIdx >= bitLen {
				break
			}

			// Row index at start of each row
			if col == 0 {
				rowStr := fmt.Sprintf("%*d", rowWidth, row)
				builder.WriteString(rowStr)
				builder.WriteByte(' ')
			}

			// Right-align cell value within colWidth
			if (word>>bitOff)&1 == 1 {
				builder.WriteString(fmt.Sprintf("%*c", colWidth, '#'))
			} else {
				builder.WriteString(fmt.Sprintf("%*c", colWidth, '.'))
			}

			// Space separator after each column except last
			if col < cols-1 {
				builder.WriteByte(' ')
			}

			bitIdx++
			col++

			// End of row
			if col == cols {
				col = 0
				row++
				// Newline (except after last row)
				if bitIdx < bitLen {
					builder.WriteByte('\n')
				}
			}
		}
	}

	return builder.String()
}
