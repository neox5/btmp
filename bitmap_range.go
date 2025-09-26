package btmp

// maskForRange calculates bit masks for range [startBit, endBit).
// Returns headMask for first partial word, tailMask for last partial word.
// For single-word ranges, only headMask is set.
func maskForRange(startBit, endBit int) (headMask, tailMask uint64) {
	w0, off0 := wordIndex(startBit)
	w1, off1 := wordIndex(endBit)

	if w0 == w1 {
		headMask = MaskRange(uint(off0), uint(off1))
		return
	}

	if off0 != 0 {
		headMask = MaskFrom(uint(off0))
	}
	if off1 != 0 {
		tailMask = MaskUpto(uint(off1))
	}
	return
}

// setRange sets bits in [start, start+count) to 1.
// Internal implementation - no validation, no auto-growth, no finalization.
func (b *Bitmap) setRange(start, count int) {
	if count == 0 {
		return
	}
	end := start + count
	w0, _ := wordIndex(start)
	w1, _ := wordIndex(end)

	headMask, tailMask := maskForRange(start, end)

	// Single word case
	if w0 == w1 {
		b.words[w0] |= headMask
		return
	}

	// Head partial word
	if headMask != 0 {
		b.words[w0] |= headMask
	}

	// Middle full words
	for w := w0 + 1; w < w1; w++ {
		b.words[w] = WordMask
	}

	// Tail partial word
	if tailMask != 0 {
		b.words[w1] |= tailMask
	}
}

// clearRange clears bits in [start, start+count) to 0.
// Internal implementation - no validation, no auto-growth, no finalization.
func (b *Bitmap) clearRange(start, count int) {
	if count == 0 {
		return
	}
	end := start + count
	w0, _ := wordIndex(start)
	w1, _ := wordIndex(end)

	headMask, tailMask := maskForRange(start, end)

	// Single word case
	if w0 == w1 {
		b.words[w0] &^= headMask
		return
	}

	// Head partial word
	if headMask != 0 {
		b.words[w0] &^= headMask
	}

	// Middle full words
	for w := w0 + 1; w < w1; w++ {
		b.words[w] = 0
	}

	// Tail partial word
	if tailMask != 0 {
		b.words[w1] &^= tailMask
	}
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
// Uses getBits/setBits from bitmap_bits.go for bit extraction and insertion.
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
// Equivalent to copyRange followed by clearing the non-overlapping source range.
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
