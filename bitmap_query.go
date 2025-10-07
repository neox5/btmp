package btmp

import "math/bits"

// test reports whether bit pos is set.
// Internal implementation - no validation.
func (b *Bitmap) test(pos int) bool {
	w := wordIdx(pos)
	off := bitOffset(pos)
	return (b.words[w]>>off)&1 == 1
}

// any reports whether any bit in [0, Len()) is set.
// Internal implementation - no validation.
func (b *Bitmap) any() bool {
	if b.lenBits == 0 {
		return false
	}
	// full words except the last
	for i := range b.lastWordIdx {
		if b.words[i] != 0 {
			return true
		}
	}
	// masked last word
	return (b.words[b.lastWordIdx] & b.tailMask) != 0
}

// all reports whether all bits in [0, Len()) are set.
// Internal implementation - no validation.
func (b *Bitmap) all() bool {
	if b.lenBits == 0 {
		return true // vacuously true for empty bitmap
	}
	// full words except the last
	for i := range b.lastWordIdx {
		if b.words[i] != WordMask {
			return false
		}
	}
	// masked last word
	return (b.words[b.lastWordIdx] & b.tailMask) == b.tailMask
}

// count returns the number of set bits in [0, Len()).
// Internal implementation - no validation.
func (b *Bitmap) count() int {
	if b.lenBits == 0 {
		return 0
	}
	sum := 0
	for i := range b.lastWordIdx {
		sum += bits.OnesCount64(b.words[i])
	}
	return sum + bits.OnesCount64(b.words[b.lastWordIdx]&b.tailMask)
}

// nextZero returns the position of the next zero bit at or after pos.
// Returns -1 if no zero bit exists in [pos, Len()).
// Internal implementation - no validation.
func (b *Bitmap) nextZero(pos int) int {
	return b.nextBitInRange(pos, b.lenBits-pos, false)
}

// nextOne returns the position of the next set bit at or after pos.
// Returns -1 if no set bit exists in [pos, Len()).
// Internal implementation - no validation.
func (b *Bitmap) nextOne(pos int) int {
	return b.nextBitInRange(pos, b.lenBits-pos, true)
}

// nextZeroInRange returns the position of the next zero bit in [pos, pos+count).
// Returns -1 if no zero bit exists in range.
// Internal implementation - no validation.
func (b *Bitmap) nextZeroInRange(pos, count int) int {
	return b.nextBitInRange(pos, count, false)
}

// nextOneInRange returns the position of the next set bit in [pos, pos+count).
// Returns -1 if no set bit exists in range.
// Internal implementation - no validation.
func (b *Bitmap) nextOneInRange(pos, count int) int {
	return b.nextBitInRange(pos, count, true)
}

// nextBitInRange returns the position of the next bit matching target value in [pos, pos+count).
// If target is true, searches for set bits (1). If false, searches for zero bits (0).
// Returns -1 if no matching bit exists in range.
// Internal implementation - no validation.
func (b *Bitmap) nextBitInRange(pos, count int, target bool) int {
	if count == 0 {
		return -1
	}

	currentPos := pos
	for word, mask := range b.rangeWords(pos, count) {
		var matched uint64
		if target {
			// Looking for set bits
			matched = *word & mask
		} else {
			// Looking for zero bits (invert and mask)
			matched = (^*word) & mask
		}

		if matched != 0 {
			tz := bits.TrailingZeros64(matched)
			w := wordIdx(currentPos)
			bitPos := w*WordBits + tz
			if bitPos < pos+count {
				return bitPos
			}
			return -1
		}

		// Advance position by number of bits checked in this word
		bitsInMask := bits.OnesCount64(mask)
		currentPos += bitsInMask
	}

	return -1
}

// countZerosFrom counts consecutive zero bits starting at pos.
// Returns 0 if bit at pos is set.
// Stops at first set bit or end of bitmap.
// Internal implementation - no validation.
func (b *Bitmap) countZerosFrom(pos int) int {
	return b.countBitsFromInRange(pos, b.lenBits-pos, false)
}

// countOnesFrom counts consecutive set bits starting at pos.
// Returns 0 if bit at pos is clear.
// Stops at first zero bit or end of bitmap.
// Internal implementation - no validation.
func (b *Bitmap) countOnesFrom(pos int) int {
	return b.countBitsFromInRange(pos, b.lenBits-pos, true)
}

// countZerosFromInRange counts consecutive zero bits starting at pos within [pos, pos+count).
// Returns 0 if bit at pos is set.
// Stops at first set bit or end of range.
// Internal implementation - no validation.
func (b *Bitmap) countZerosFromInRange(pos, count int) int {
	return b.countBitsFromInRange(pos, count, false)
}

// countOnesFromInRange counts consecutive set bits starting at pos within [pos, pos+count).
// Returns 0 if bit at pos is clear.
// Stops at first zero bit or end of range.
// Internal implementation - no validation.
func (b *Bitmap) countOnesFromInRange(pos, count int) int {
	return b.countBitsFromInRange(pos, count, true)
}

// countBitsFromInRange counts consecutive bits matching target value starting at pos within [pos, pos+count).
// If target is true, counts set bits (1). If false, counts zero bits (0).
// Returns 0 if bit at pos doesn't match target.
// Stops at first non-matching bit or end of range.
// Internal implementation - no validation.
func (b *Bitmap) countBitsFromInRange(pos, count int, target bool) int {
	if count == 0 {
		return 0
	}

	// Check if starting bit matches target
	if b.test(pos) != target {
		return 0
	}

	bitCount := 0

	for word, mask := range b.rangeWords(pos, count) {
		var matched uint64
		if target {
			// Counting set bits
			matched = *word & mask
		} else {
			// Counting zero bits (invert and mask)
			matched = (^*word) & mask
		}

		if matched == 0 {
			// Entire masked region doesn't match target, stop
			break
		}

		// Find first non-matching bit
		inverted := (^matched) & mask
		tz := bits.TrailingZeros64(inverted)
		bitsInMask := bits.OnesCount64(mask)

		if tz < bitsInMask {
			// Found a non-matching bit
			bitCount += tz
			break
		}

		// All bits in this masked region match target
		bitCount += bitsInMask
	}

	return bitCount
}
