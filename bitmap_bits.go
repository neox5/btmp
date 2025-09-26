package btmp

// setBit sets bit at position i to 1.
// Internal implementation - no validation, no finalization.
func (b *Bitmap) setBit(i int) {
	w, off := wordIndex(i)
	b.words[w] |= uint64(1) << off
}

// clearBit sets bit at position i to 0.
// Internal implementation - no validation, no finalization.
func (b *Bitmap) clearBit(i int) {
	w, off := wordIndex(i)
	b.words[w] &^= uint64(1) << off
}

// flipBit toggles bit at position i.
// Internal implementation - no validation, no finalization.
func (b *Bitmap) flipBit(i int) {
	w, off := wordIndex(i)
	b.words[w] ^= uint64(1) << off
}

// getBits extracts n bits starting from pos, returned right-aligned.
// No validation performed - caller must ensure bounds.
//
// Fast paths:
//   - Full word aligned: direct word access
//   - Single word unaligned: shift and mask
//   - Cross word: combine low bits from first word + high bits from second word
func (b *Bitmap) getBits(pos, n int) uint64 {
	// Fast path: full word aligned read
	if n == WordBits && (pos&IndexMask) == 0 {
		w := pos >> WordShift
		return b.words[w]
	}

	w, off := wordIndex(pos)

	// Fast path: single word case
	if off+n <= WordBits {
		word := b.words[w]
		return (word >> off) & MaskUpto(uint(n))
	}

	// Cross-word case: spans exactly two words
	bitsFromFirst := WordBits - off
	bitsFromSecond := n - bitsFromFirst

	// Extract low bits from first word
	firstWord := b.words[w]
	lowBits := firstWord >> off

	// Extract high bits from second word and position them
	secondWord := b.words[w+1]
	highBits := secondWord & MaskUpto(uint(bitsFromSecond))

	return lowBits | (highBits << bitsFromFirst)
}

// setBits inserts the low n bits of val into the bitmap starting at pos.
// Preserves surrounding bits. No validation performed - caller must ensure bounds.
//
// Fast paths:
//   - Full word aligned: direct word assignment
//   - Single word unaligned: mask and insert
//   - Cross word: split value between two words
func (b *Bitmap) setBits(pos, n int, val uint64) {
	// Fast path: full word aligned write
	if n == WordBits && (pos&IndexMask) == 0 {
		w := pos >> WordShift
		b.words[w] = val
		return
	}

	w, off := wordIndex(pos)

	// Mask val to exactly n bits to prevent overflow
	maskedVal := val & MaskUpto(uint(n))

	// Fast path: single word case
	if off+n <= WordBits {
		mask := MaskUpto(uint(n)) << off
		insertVal := maskedVal << off
		b.words[w] = (b.words[w] &^ mask) | insertVal
		return
	}

	// Cross-word case: spans exactly two words
	bitsToFirst := WordBits - off
	bitsToSecond := n - bitsToFirst

	// First word: insert low bits of value
	maskFirst := MaskUpto(uint(bitsToFirst)) << off
	lowVal := maskedVal << off
	b.words[w] = (b.words[w] &^ maskFirst) | lowVal

	// Second word: insert high bits of value
	maskSecond := MaskUpto(uint(bitsToSecond))
	highVal := maskedVal >> bitsToFirst
	b.words[w+1] = (b.words[w+1] &^ maskSecond) | highVal
}
