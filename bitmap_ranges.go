package btmp

// SetRange sets bits in [start, start+count). Auto-grows to fit end.
// No-op if count == 0. Returns b.
//
// Invariant: after return, all bits >= Len() are zero, even when count == 0.
func (b *Bitmap) SetRange(start, count int) *Bitmap {
	defer finalize(b)
	if count < 0 {
		panic("SetRange: negative count")
	}
	if count == 0 {
		return b
	}
	end := checkedEnd(start, count)
	b.EnsureBits(end)

	// head partial
	w0, o0 := wordIndex(start)
	if o0 != 0 {
		n := min(count, int(64-o0))
		mask := ^uint64(0) << o0
		if n < int(64-o0) {
			mask &= (uint64(1) << (o0 + uint(n))) - 1
		}
		b.words[w0] |= mask
		start += n
		count -= n
	}

	// middle full words
	for count >= 64 {
		w, _ := wordIndex(start)
		b.words[w] = ^uint64(0)
		start += 64
		count -= 64
	}

	// tail
	if count > 0 {
		w, _ := wordIndex(start)
		mask := (uint64(1) << uint(count)) - 1
		b.words[w] |= mask
	}

	return b
}

// ClearRange clears bits in [start, start+count). In-bounds only.
// No-op if count == 0. Returns b.
//
// Invariant: after return, all bits >= Len() are zero, even when count == 0.
func (b *Bitmap) ClearRange(start, count int) *Bitmap {
	defer finalize(b)
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

	// head partial
	w0, o0 := wordIndex(start)
	if o0 != 0 {
		n := min(count, int(64-o0))
		mask := ^uint64(0) << o0
		if n < int(64-o0) {
			mask &= (uint64(1) << (o0 + uint(n))) - 1
		}
		b.words[w0] &^= mask
		start += n
		count -= n
	}

	// middle
	for count >= 64 {
		w, _ := wordIndex(start)
		b.words[w] = 0
		start += 64
		count -= 64
	}

	// tail
	if count > 0 {
		w, _ := wordIndex(start)
		mask := (uint64(1) << uint(count)) - 1
		b.words[w] &^= mask
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
	defer finalize(b)
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

	// destination growth
	dstEnd := checkedEnd(dstStart, count)
	b.EnsureBits(dstEnd)

	// choose direction considering overlap only when src == dst
	if src != b {
		copyForward(src, b, srcStart, dstStart, count)
		return b
	}
	if dstStart < srcStart || dstStart >= srcStart+count {
		copyForward(b, b, srcStart, dstStart, count)
	} else {
		copyBackward(b, b, srcStart, dstStart, count)
	}
	return b
}

// --- internal helpers for range ops ---

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

// copyForward copies left-to-right.
func copyForward(src, dst *Bitmap, srcStart, dstStart, count int) {
	// small ranges: scalar
	if count < 64 {
		for i := range count {
			sw, so := wordIndex(srcStart + i)
			dw, do := wordIndex(dstStart + i)
			bit := (src.words[sW(sw)]>>so)&1 == 1
			if bit {
				dst.words[dw] |= 1 << do
			} else {
				dst.words[dw] &^= 1 << do
			}
		}
		return
	}

	// align destination to word boundary
	for dstStart&63 != 0 && count > 0 {
		sw, so := wordIndex(srcStart)
		dw, do := wordIndex(dstStart)
		bit := (src.words[sW(sw)]>>so)&1 == 1
		if bit {
			dst.words[dw] |= 1 << do
		} else {
			dst.words[dw] &^= 1 << do
		}
		srcStart++
		dstStart++
		count--
	}

	// word-wise copies
	for count >= 64 {
		dw := dstStart >> 6
		soff := uint(srcStart & 63)
		sw := srcStart >> 6

		var v uint64
		if soff == 0 {
			v = src.words[sW(sw)]
		} else {
			lo := src.words[sW(sw)] >> soff
			hi := uint64(0)
			if sw+1 < len(src.words) {
				hi = src.words[sW(sw+1)] << (64 - soff)
			}
			v = lo | hi
		}
		dst.words[dw] = v
		srcStart += 64
		dstStart += 64
		count -= 64
	}

	// tail scalar
	for i := range count {
		sw, so := wordIndex(srcStart + i)
		dw, do := wordIndex(dstStart + i)
		bit := (src.words[sW(sw)]>>so)&1 == 1
		if bit {
			dst.words[dw] |= 1 << do
		} else {
			dst.words[dw] &^= 1 << do
		}
	}
}

// copyBackward copies right-to-left. Must not compute negative intermediate
// indices. Must handle count < 64 via a scalar path before any word loop.
func copyBackward(src, dst *Bitmap, srcStart, dstStart, count int) {
	// small ranges: scalar from end
	if count < 64 {
		for i := count - 1; i >= 0; i-- {
			sw, so := wordIndex(srcStart + i)
			dw, do := wordIndex(dstStart + i)
			bit := (src.words[sW(sw)]>>so)&1 == 1
			if bit {
				dst.words[dw] |= 1 << do
			} else {
				dst.words[dw] &^= 1 << do
			}
		}
		return
	}

	end := dstStart + count
	srcEnd := srcStart + count

	// align end down to word boundary using scalar from the end
	for end&63 != 0 && count > 0 {
		end--
		srcEnd--
		count--
		sw, so := wordIndex(srcEnd)
		dw, do := wordIndex(end)
		bit := (src.words[sW(sw)]>>so)&1 == 1
		if bit {
			dst.words[dw] |= 1 << do
		} else {
			dst.words[dw] &^= 1 << do
		}
	}

	// word-wise downward copies
	for count >= 64 {
		dw := (end - 1) >> 6
		soff := uint(srcEnd & 63)
		sw := (srcEnd - 1) >> 6

		var v uint64
		if soff == 0 {
			v = src.words[sW(sw)]
		} else {
			hi := src.words[sW(sw)] << (64 - soff)
			lo := uint64(0)
			if sw-1 >= 0 {
				lo = src.words[sW(sw-1)] >> soff
			}
			v = lo | hi
		}
		dst.words[dw] = v

		end -= 64
		srcEnd -= 64
		count -= 64
	}

	// head remainder scalar
	for i := count - 1; i >= 0; i-- {
		sw, so := wordIndex(srcStart + i)
		dw, do := wordIndex(dstStart + i)
		bit := (src.words[sW(sw)]>>so)&1 == 1
		if bit {
			dst.words[dw] |= 1 << do
		} else {
			dst.words[dw] &^= 1 << do
		}
	}
}

// sW bounds src word index safely for zero-length word slices.
func sW(i int) int {
	if i < 0 {
		return 0
	}
	return i
}
