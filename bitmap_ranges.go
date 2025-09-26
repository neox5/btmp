package btmp

// SetRange sets bits in [start, start+count) to 1. Auto-grows if needed.
// No-op if count == 0. Returns b for chaining.
func (b *Bitmap) SetRange(start, count int) *Bitmap {
	defer b.finalize()
	if count == 0 {
		return b
	}
	b.EnsureBits(start + count)
	end := start + count

	w0, off0 := wordIndex(start)
	w1, off1 := wordIndex(end)

	// Single word.
	if w0 == w1 {
		b.words[w0] |= MaskRange(uint(off0), uint(off1))
		return b
	}

	// Head partial word.
	if off0 != 0 {
		b.words[w0] |= MaskFrom(uint(off0))
	}

	// Middle full words.
	for w := w0 + 1; w < w1; w++ {
		b.words[w] = WordMask
	}

	// Tail partial word.
	if off1 != 0 {
		b.words[w1] |= MaskUpto(uint(off1))
	}

	return b
}

// ClearRange clears bits in [start, start+count) to 0. In-bounds only.
// No-op if count == 0. Returns b for chaining.
func (b *Bitmap) ClearRange(start, count int) *Bitmap {
	defer b.finalize()
	if count == 0 {
		return b
	}
	end := b.checkedEnd(start, count)

	w0, off0 := wordIndex(start)
	w1, off1 := wordIndex(end)

	// Single word.
	if w0 == w1 {
		b.words[w0] &^= MaskRange(uint(off0), uint(off1))
		return b
	}

	// Head partial word.
	if off0 != 0 {
		b.words[w0] &^= MaskFrom(uint(off0))
	}

	// Middle full words.
	for w := w0 + 1; w < w1; w++ {
		b.words[w] = 0
	}

	// Tail partial word.
	if off1 != 0 {
		b.words[w1] &^= MaskUpto(uint(off1))
	}

	return b
}

// CopyRange copies count bits from src[srcStart:] to dst[dstStart:].
// In-bounds only - no auto-grow. Overlap-safe with memmove semantics.
// No-op if count == 0. Returns b for chaining.
func (b *Bitmap) CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap {
	defer b.finalize()
	if count == 0 || srcStart == dstStart {
		return b
	}
	
	// Validate bounds - no auto-grow
	_ = src.checkedEnd(srcStart, count)
	_ = b.checkedEnd(dstStart, count)
	
	// Determine copy direction for overlap safety
	backward := needsBackwardCopy(srcStart, dstStart, count)
	
	// Perform bit-level copy
	copyBitRange(b, src, srcStart, dstStart, count, backward)
	return b
}

// MoveRange moves bits from [src, src+count) to [dst, dst+count).
// Equivalent to CopyRange followed by clearing the source range.
// Auto-grows destination if needed. Returns b for chaining.
func (b *Bitmap) MoveRange(src, dst, count int) *Bitmap {
	defer b.finalize()
	b.CopyRange(b, src, dst, count)
	if count > 0 && src != dst {
		b.ClearRange(src, count)
	}
	return b
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
	sp := srcStart  // source position
	dp := dstStart  // dest position
	
	if backward {
		sp += count - WordBits
		dp += count - WordBits
	}
	
	for remaining > 0 {
		n := min(remaining, WordBits)  // bits to process this iteration
		
		if backward && n < WordBits {
			// Adjust position for final partial chunk
			adj := WordBits - n
			sp += adj
			dp += adj
		}
		
		// Extract bits from source
		bits := extractBits(src, sp, n)
		
		// Insert bits into destination  
		insertBits(dst, dp, n, bits)
		
		remaining -= n
		if backward {
			sp -= WordBits  // always step by full word size
			dp -= WordBits
		} else {
			sp += n  // step by actual bits processed
			dp += n
		}
	}
}

// extractBits extracts n bits starting from pos, returned right-aligned.
func extractBits(src *Bitmap, pos int, n int) uint64 {
	w, off := wordIndex(pos)
	
	// Single word case - most common
	if off+n <= WordBits {
		word := src.words[w]
		return (word >> off) & MaskUpto(uint(n))
	}
	
	// Spans two words case:
	// Source:     Word w: [ . . . . high ] [ low, off <- pos ]
	//             Word w+1:[ high, 0 <- ... ] [ . . . . . . . . ]
	// Transform:  low  = w >> off          (shift right by off)
	//             high = w+1 & mask        (mask low bitsH bits)
	// Result:     [ . . . . . . . ] [ high | low ]
	wL := src.words[w]
	bitsL := WordBits - off  // bits from first word (low bits of result)
	bitsH := n - bitsL       // bits from second word (high bits of result)
	
	// Extract low bits from first word (shift right to position 0)
	low := wL >> off
	
	// Extract high bits from second word and position them after low bits
	wH := src.words[w+1]  // checkedEnd guarantees this exists
	high := wH & MaskUpto(uint(bitsH))
	
	return low | (high << bitsL)
}

// insertBits inserts the low n bits of val into dst starting at pos.
func insertBits(dst *Bitmap, pos int, n int, val uint64) {
	w, off := wordIndex(pos)
	
	// Mask val to exactly n bits
	maskedVal := val & MaskUpto(uint(n)) 
	
	// Single word case
	if off+n <= WordBits {
		mask := MaskUpto(uint(n)) << off
		v := maskedVal << off
		dst.words[w] = (dst.words[w] &^ mask) | v
		return
	}
	
	// Spans two words case:
	// Source:     Value: [ . . . . . . . ] [ high | low ]
	// Transform:  lowVal  = maskedVal << off     (shift left by off)
	//             highVal = maskedVal >> bitsL   (shift right by bitsL)
	// Target:     Word w: [ . . . . . . . ] [ lowVal, off <- pos ]
	//             Word w+1:[ highVal, 0 <- ... ] [ . . . . . . . . ]
	bitsL := WordBits - off  // bits going to first word
	bitsH := n - bitsL       // bits going to second word
	
	// First word: insert low bits of value
	maskL := MaskUpto(uint(bitsL)) << off
	lowVal := maskedVal << off
	dst.words[w] = (dst.words[w] &^ maskL) | lowVal
	
	// Second word: insert high bits of value
	maskH := MaskUpto(uint(bitsH))
	highVal := maskedVal >> bitsL
	dst.words[w+1] = (dst.words[w+1] &^ maskH) | highVal
}
