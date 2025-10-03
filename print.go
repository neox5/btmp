package btmp

import (
	"fmt"
	"strings"
)

// formatBits formats a bit sequence into a string representation.
//
// Parameters:
//   - bits: source bits, right-aligned (low bits used if bitCount < 64)
//   - bitCount: number of valid bits to format (1-64)
//   - base: output base (2 for binary, 16 for hexadecimal)
//   - grouped: if true, insert separators between groups
//   - groupSize: units per group - for base 2: bits, for base 16: hex digits
//   - sep: separator string inserted between groups
//
// For base 16:
//   - Groups 4 bits per hex digit, left-to-right
//   - Right-pads incomplete final group with zeros
//   - Example: 6 bits "101100" → "B0" (treated as "10110000")
//
// For base 2:
//   - Outputs '0' and '1' characters in index order (left-to-right)
//   - No padding
//
// Grouping:
//   - Inserts sep every groupSize output units
//   - For base 2: groupSize is bit count
//   - For base 16: groupSize is hex digit count
//   - Last group may be shorter than groupSize
//   - Example base 2: bits=0xFF, bitCount=8, groupSize=4 → "1111_1111"
//   - Example base 16: bits=0xABCD, bitCount=16, groupSize=2 → "AB CD"
//
// Panics if bitCount <= 0, bitCount > 64, base not in {2,16},
// or grouped && groupSize <= 0.
func formatBits(bits uint64, bitCount int, base int, grouped bool, groupSize int, sep string) string {
	// Validation
	if bitCount <= 0 || bitCount > WordBits {
		panic("bitCount must be > 0 and <= 64")
	}
	if base != 2 && base != 16 {
		panic("base must be 2 or 16")
	}
	if grouped && groupSize <= 0 {
		panic("groupSize must be positive when grouped")
	}

	var s string
	if base == 2 {
		s = formatBinary(bits, bitCount)
	} else { // base == 16
		s = formatHex(bits, bitCount)
	}

	if grouped {
		s = applyGrouping(s, groupSize, sep)
	}

	return s
}

// formatBinary formats bits as binary string with exact bitCount digits.
// Pads left with zeros if needed. Takes rightmost bitCount bits.
// Internal helper - no validation, no grouping.
func formatBinary(bits uint64, bitCount int) string {
	s := fmt.Sprintf("%b", bits)

	// Pad left if needed
	if len(s) < bitCount {
		s = strings.Repeat("0", bitCount-len(s)) + s
	}

	// Take rightmost bitCount characters
	return s[len(s)-bitCount:]
}

// formatHex formats bits as hexadecimal string (uppercase).
// Right-pads to complete hex digit if bitCount not divisible by 4.
// Internal helper - no validation, no grouping.
func formatHex(bits uint64, bitCount int) string {
	// Calculate number of hex digits needed (ceiling division)
	hexDigits := (bitCount + 3) / 4

	// Create format string with zero-padding
	format := fmt.Sprintf("%%0%dX", hexDigits)

	return fmt.Sprintf(format, bits)
}

// applyGrouping inserts separators every groupSize characters from left to right.
// Last group may be shorter than groupSize.
// Internal helper - no validation.
func applyGrouping(s string, groupSize int, sep string) string {
	if groupSize <= 0 || groupSize >= len(s) {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(s) + (len(s)/groupSize)*len(sep))

	for i, c := range s {
		if i > 0 && i%groupSize == 0 {
			builder.WriteString(sep)
		}
		builder.WriteRune(c)
	}

	return builder.String()
}
