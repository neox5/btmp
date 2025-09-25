package btmp

// SetRange sets bits in [start, start+count). Auto-grows to fit end.
// No-op if count == 0. Returns b.
//
// Invariant: after return, all bits >= Len() are zero, even when count == 0.
func (b *Bitmap) SetRange(start, count int) *Bitmap {
	defer b.finalize()
	if count < 0 {
		panic("SetRange: negative count")
	}
	if count == 0 {
		return b
	}
	end := checkedEnd(start, count)
	b.EnsureBits(end)

	// head partial word
	w, off := wordIndex(start)
	if off != 0 {
		n := min(count, wordBits-int(off))
		mask := MaskFrom(off)
		if n < wordBits-int(off) {
			mask &= MaskUpto(off + uint(n))
		}
		b.words[w] |= mask
		start += n
		count -= n
	}

	// middle full words
	for count >= wordBits {
		w, _ := wordIndex(start)
		b.words[w] = wordMask
		start += wordBits
		count -= wordBits
	}

	// tail partial word
	if count > 0 {
		w, _ := wordIndex(start)
		b.words[w] |= MaskUpto(uint(count))
	}

	return b
}

// ClearRange clears bits in [start, start+count). In-bounds only.
// No-op if count == 0. Returns b.
//
// Invariant: after return, all bits >= Len() are zero, even when count == 0.
func (b *Bitmap) ClearRange(start, count int) *Bitmap {
	defer b.finalize()
	if count < 0 {
		panic("ClearRange: negative count")
	}
	if count == 0 {
		return b
	}
	end := checkedEnd(start, count)
	if end > b.lenBits {
		panic("ClearRange: out of bounds")
	}

	// head partial word
	w, off := wordIndex(start)
	if off != 0 {
		n := min(count, wordBits-int(off))
		mask := MaskFrom(off)
		if n < wordBits-int(off) {
			mask &= MaskUpto(off + uint(n))
		}
		b.words[w] &^= mask
		start += n
		count -= n
	}

	// middle full words
	for count >= wordBits {
		w, _ := wordIndex(start)
		b.words[w] = 0
		start += wordBits
		count -= wordBits
	}

	// tail partial word
	if count > 0 {
		w, _ := wordIndex(start)
		b.words[w] &^= MaskUpto(uint(count))
	}

	return b
}

// CopyRange copies count bits from src at srcStart to b at dstStart with memmove semantics.
// Auto-grows destination to fit. Panics on nil src or out-of-bounds source. Returns b.
//
// Overlap rule: if src == b and ranges overlap, the copy behaves like memmove.
// Direction selection: forward if dstStart < srcStart or non-overlap, else backward.
//
// Invariant: after return, all bits >= Len() are zero, even when count == 0.
func (b *Bitmap) CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap {
	defer b.finalize()
	if count < 0 {
		panic("CopyRange: negative count")
	}
	if count == 0 {
		return b
	}
	if src == nil {
		panic("CopyRange: nil src")
	}
	srcEnd := checkedEnd(srcStart, count)
	if srcEnd > src.lenBits {
		panic("CopyRange: source out of bounds")
	}

	dstEnd := checkedEnd(dstStart, count)
	b.EnsureBits(dstEnd)

	if src != b || dstStart < srcStart || dstStart >= srcEnd {
		copyForward(src, b, srcStart, dstStart, count)
	} else {
		copyBackward(b, b, srcStart, dstStart, count)
	}
	return b
}

// checkedEnd validates start and count and returns start+count or panics.
func checkedEnd(start, count int) int {
	if start < 0 {
		panic("negative start")
	}
	if count < 0 {
		panic("negative count")
	}
	end := start + count
	if end < start {
		panic("overflow")
	}
	return end
}

// copyBit copies a single bit from src to dst.
func copyBit(src, dst *Bitmap, srcIdx, dstIdx int) {
	sw, so := wordIndex(srcIdx)
	dw, do := wordIndex(dstIdx)
	if sw < len(src.words) && (src.words[sw]>>so)&1 == 1 {
		dst.words[dw] |= 1 << do
	} else {
		dst.words[dw] &^= 1 << do
	}
}

// extractWord extracts a 64-bit word from src starting at bit position.
func extractWord(src *Bitmap, bitPos int) uint64 {
	sw, soff := wordIndex(bitPos)
	if sw >= len(src.words) {
		return 0
	}

	var v uint64
	if soff == 0 {
		v = src.words[sw]
	} else {
		v = src.words[sw] >> soff
		if sw+1 < len(src.words) {
			v |= src.words[sw+1] << (wordBits - soff)
		}
	}
	return v
}

// copyForward copies left-to-right.
func copyForward(src, dst *Bitmap, srcStart, dstStart, count int) {
	// small ranges: bit-by-bit
	if count < wordBits {
		for i := range count {
			copyBit(src, dst, srcStart+i, dstStart+i)
		}
		return
	}

	// align destination to word boundary
	for dstStart&indexMask != 0 && count > 0 {
		copyBit(src, dst, srcStart, dstStart)
		srcStart++
		dstStart++
		count--
	}

	// word-aligned copies
	for count >= wordBits {
		dst.words[dstStart>>wordShift] = extractWord(src, srcStart)
		srcStart += wordBits
		dstStart += wordBits
		count -= wordBits
	}

	// tail bits
	for i := range count {
		copyBit(src, dst, srcStart+i, dstStart+i)
	}
}

// copyBackward copies right-to-left.
func copyBackward(src, dst *Bitmap, srcStart, dstStart, count int) {
	// small ranges: bit-by-bit from end
	if count < wordBits {
		for i := count - 1; i >= 0; i-- {
			copyBit(src, dst, srcStart+i, dstStart+i)
		}
		return
	}

	end := dstStart + count
	srcEnd := srcStart + count

	// align end down to word boundary
	for end&indexMask != 0 && count > 0 {
		end--
		srcEnd--
		count--
		copyBit(src, dst, srcEnd, end)
	}

	// word-aligned copies backward
	for count >= wordBits {
		end -= wordBits
		srcEnd -= wordBits
		dst.words[end>>wordShift] = extractWord(src, srcEnd)
		count -= wordBits
	}

	// head remainder bits
	for i := count - 1; i >= 0; i-- {
		copyBit(src, dst, srcStart+i, dstStart+i)
	}
}
