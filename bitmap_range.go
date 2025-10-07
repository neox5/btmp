package btmp

import (
	"iter"
	"math/bits"
)

// ========================================
// Range Helper Functions
// ========================================

// wordIdx returns the word index for a bit position.
func wordIdx(pos int) int {
	return pos >> WordShift
}

// bitOffset returns the bit offset within a word [0, 63].
func bitOffset(pos int) int {
	return pos & IndexMask
}

// rangeWordIndices returns the first and last word indices for a bit range.
func rangeWordIndices(start, count int) (w0, w1 int) {
	if count == 0 {
		return 0, -1
	}
	w0 = wordIdx(start)
	w1 = wordIdx(start + count - 1)
	return w0, w1
}

// headMaskForRange returns the mask for the first word of a range.
func headMaskForRange(start, count int) uint64 {
	if count == 0 {
		return 0
	}
	startBit := bitOffset(start)
	endBit := min(WordBits, startBit+count) - 1
	return MaskRange(uint(startBit), uint(endBit+1))
}

// tailMaskForRange returns the mask for the last word of a range.
func tailMaskForRange(start, count int) uint64 {
	if count == 0 {
		return 0
	}
	endBit := bitOffset(start + count - 1)
	return MaskUpto(uint(endBit + 1))
}

// rangeWords returns an iterator over words in a range with masks.
func (b *Bitmap) rangeWords(start, count int) iter.Seq2[*uint64, uint64] {
	return func(yield func(*uint64, uint64) bool) {
		if count == 0 {
			return
		}

		w0, w1 := rangeWordIndices(start, count)

		if w0 == w1 {
			yield(&b.words[w0], headMaskForRange(start, count))
			return
		}

		if !yield(&b.words[w0], headMaskForRange(start, count)) {
			return
		}

		for w := w0 + 1; w < w1; w++ {
			if !yield(&b.words[w], WordMask) {
				return
			}
		}

		yield(&b.words[w1], tailMaskForRange(start, count))
	}
}

// ========================================
// Range Operation Implementations
// ========================================

// setRange sets bits in [start, start+count) to 1.
// Internal implementation - no validation, no auto-growth, no finalization.
func (b *Bitmap) setRange(start, count int) {
	for word, mask := range b.rangeWords(start, count) {
		*word |= mask
	}
}

// clearRange clears bits in [start, start+count) to 0.
// Internal implementation - no validation, no auto-growth, no finalization.
func (b *Bitmap) clearRange(start, count int) {
	for word, mask := range b.rangeWords(start, count) {
		*word &^= mask
	}
}

// anyRange reports whether any bit in [start, start+count) is set.
// Internal implementation - no validation.
func (b *Bitmap) anyRange(start, count int) bool {
	if count == 0 {
		return false
	}

	for word, mask := range b.rangeWords(start, count) {
		if (*word & mask) != 0 {
			return true
		}
	}
	return false
}

// allRange reports whether all bits in [start, start+count) are set.
// Internal implementation - no validation.
func (b *Bitmap) allRange(start, count int) bool {
	if count == 0 {
		return true // vacuously true for empty range
	}

	for word, mask := range b.rangeWords(start, count) {
		if (*word & mask) != mask {
			return false
		}
	}
	return true
}

// countRange returns the number of set bits in [start, start+count).
// Internal implementation - no validation.
func (b *Bitmap) countRange(start, count int) int {
	if count == 0 {
		return 0
	}

	sum := 0
	for word, mask := range b.rangeWords(start, count) {
		sum += bits.OnesCount64(*word & mask)
	}
	return sum
}

// copyRange copies count bits from src[srcStart:] to dst[dstStart:].
// Internal implementation - no validation, no auto-growth, no finalization.
// Overlap-safe with memmove semantics.
func (b *Bitmap) copyRange(src *Bitmap, srcStart, dstStart, count int) {
	if count == 0 || srcStart == dstStart {
		return
	}

	// Determine copy direction for overlap safety
	backward := needsBackwardCopy(srcStart, dstStart, count)

	// Perform bit-level copy
	copyBitRange(b, src, srcStart, dstStart, count, backward)
}

// needsBackwardCopy determines if backward iteration is needed for safe overlapping copy.
func needsBackwardCopy(srcStart, dstStart, count int) bool {
	srcEnd := srcStart + count
	dstEnd := dstStart + count
	// Overlap exists AND dst > src requires backward copy
	return srcStart < dstEnd && dstStart < srcEnd && dstStart > srcStart
}

// copyBitRange performs the actual bit copying with proper direction handling.
func copyBitRange(dst, src *Bitmap, srcStart, dstStart, count int, backward bool) {
	remaining := count
	sp := srcStart // source position
	dp := dstStart // dest position

	if backward {
		sp += count - WordBits
		dp += count - WordBits
	}

	for remaining > 0 {
		n := min(remaining, WordBits) // bits to process this iteration

		if backward && n < WordBits {
			// Adjust position for final partial chunk
			adj := WordBits - n
			sp += adj
			dp += adj
		}

		// Extract bits from source using getBits
		bits := src.getBits(sp, n)

		// Insert bits into destination using setBits
		dst.setBits(dp, n, bits)

		remaining -= n
		if backward {
			sp -= WordBits // always step by full word size
			dp -= WordBits
		} else {
			sp += n // step by actual bits processed
			dp += n
		}
	}
}

// moveRange moves bits from [srcStart, srcStart+count) to [dstStart, dstStart+count).
// Internal implementation - no validation, no auto-growth, no finalization.
func (b *Bitmap) moveRange(srcStart, dstStart, count int) {
	b.copyRange(b, srcStart, dstStart, count)
	if count > 0 && srcStart != dstStart {
		// Clear non-overlapping parts of source
		srcEnd := srcStart + count
		dstEnd := dstStart + count

		// Clear before overlap
		if srcStart < dstStart {
			clearEnd := min(srcEnd, dstStart)
			b.clearRange(srcStart, clearEnd-srcStart)
		}

		// Clear after overlap
		if srcEnd > dstEnd {
			clearStart := max(srcStart, dstEnd)
			b.clearRange(clearStart, srcEnd-clearStart)
		}
	}
}

// setAll sets all bits in [0, Len()) to 1.
// Internal implementation - no validation, no finalization.
func (b *Bitmap) setAll() {
	if b.lenBits == 0 {
		return
	}

	// Set all full words
	for i := range b.lastWordIdx {
		b.words[i] = WordMask
	}

	// Set masked last word
	b.words[b.lastWordIdx] = b.tailMask
}

// clearAll clears all bits in [0, Len()) to 0.
// Internal implementation - no validation, no finalization.
func (b *Bitmap) clearAll() {
	if b.lenBits == 0 {
		return
	}

	// Clear all words up to and including last logical word
	for i := range b.lastWordIdx + 1 {
		b.words[i] = 0
	}
}
