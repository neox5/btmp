package btmp

import "strings"

// printRangeFormat formats bits in [start, start+count) with format parameters.
// Internal implementation - no validation.
// base: 2 (binary) or 16 (hexadecimal)
// grouped: insert separators between bit groups
// groupSize: units per group (bits for base 2, hex digits for base 16)
// sep: separator string
func (b *Bitmap) printRangeFormat(start, count int, base int, grouped bool, groupSize int, sep string) string {
	if count == 0 {
		return ""
	}

	// For ranges <= 64 bits, single format call
	if count <= WordBits {
		bits := b.getBits(start, count)
		return formatBits(bits, count, base, grouped, groupSize, sep)
	}

	// For ranges > 64 bits:
	// 1. Build ungrouped string from chunks
	// 2. Apply grouping to complete string

	var builder strings.Builder
	estimatedSize := count
	if base == 16 {
		estimatedSize = (count + 3) / 4
	}
	builder.Grow(estimatedSize)

	remaining := count
	pos := start

	for remaining > 0 {
		chunkSize := min(remaining, WordBits)
		bits := b.getBits(pos, chunkSize)
		// Format without grouping
		builder.WriteString(formatBits(bits, chunkSize, base, false, 0, ""))

		remaining -= chunkSize
		pos += chunkSize
	}

	ungrouped := builder.String()

	if grouped {
		return applyGrouping(ungrouped, groupSize, sep)
	}

	return ungrouped
}
