package btmp

import "math"

// SetRange: grow if needed, set bits, return b.
func (b *Bitmap) SetRange(start, count int) *Bitmap {
	if count == 0 {
		return b
	}
	end := checkedEnd(start, count)
	if end > b.lenBits {
		ensureLen(b, end)
	}
	ws, we := wordSpan(start, count)
	head, tail, single := partialMasks(start, count)
	if single {
		b.words[ws] |= head & tail
		maskTail(b)
		return b
	}
	b.words[ws] |= head
	for i := ws + 1; i < we-1; i++ {
		b.words[i] = ^uint64(0)
	}
	b.words[we-1] |= tail
	maskTail(b)
	return b
}

// ClearRange: in-bounds only, return b.
func (b *Bitmap) ClearRange(start, count int) *Bitmap {
	if count == 0 {
		return b
	}
	end := checkedEnd(start, count)
	if start < 0 || end > b.lenBits {
		panic("btmp: ClearRange out of bounds")
	}
	ws, we := wordSpan(start, count)
	head, tail, single := partialMasks(start, count)
	if single {
		b.words[ws] &^= head & tail
		maskTail(b)
		return b
	}
	b.words[ws] &^= head
	for i := ws + 1; i < we-1; i++ {
		b.words[i] = 0
	}
	b.words[we-1] &^= tail
	maskTail(b)
	return b
}

// CopyRange: grow dst if needed, overlap-safe, return b.
func (b *Bitmap) CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap {
	if count == 0 {
		return b
	}
	if src == nil {
		panic("btmp: CopyRange nil src")
	}
	srcEnd := checkedEnd(srcStart, count)
	if srcStart < 0 || srcEnd > src.lenBits {
		panic("btmp: CopyRange source out of bounds")
	}
	dstEnd := checkedEnd(dstStart, count)
	if dstEnd > b.lenBits {
		ensureLen(b, dstEnd)
	}
	if src == b && rangesOverlap(srcStart, srcEnd, dstStart, dstEnd) {
		if dstStart < srcStart {
			copyForward(src, b, srcStart, dstStart, count)
		} else {
			copyBackward(src, b, srcStart, dstStart, count)
		}
	} else {
		copyForward(src, b, srcStart, dstStart, count)
	}
	maskTail(b)
	return b
}

/*** internal bit math and copy engines ***/

func checkedEnd(start, count int) int {
	if start < 0 || count < 0 {
		panic("btmp: negative start or count")
	}
	end := start + count
	if end < start || end > math.MaxInt {
		panic("btmp: start+count overflow")
	}
	return end
}

func wordsFor(bits int) int {
	if bits <= 0 {
		return 0
	}
	return (bits + 63) >> 6
}

func wordSpan(start, count int) (ws, we int) {
	return start >> 6, (start+count+63)>>6
}

func partialMasks(start, count int) (head, tail uint64, single bool) {
	sOff := uint(start & 63)
	end := start + count
	eOff := uint((end - 1) & 63)
	head = ^uint64(0) << sOff
	tail = lowMask(eOff + 1)
	single = (start>>6 == (end-1)>>6)
	return
}

func tailMask(lenBits int) uint64 {
	r := uint(lenBits & 63)
	if r == 0 {
		return ^uint64(0)
	}
	return (uint64(1) << r) - 1
}

func maskTail(b *Bitmap) {
	if b.lenBits == 0 {
		return
	}
	if (b.lenBits & 63) != 0 {
		last := wordsFor(b.lenBits) - 1
		b.words[last] &= tailMask(b.lenBits)
	}
}

func lowMask(n uint) uint64 {
	if n >= 64 {
		return ^uint64(0)
	}
	return (uint64(1) << n) - 1
}

func rangesOverlap(a0, a1, b0, b1 int) bool {
	return a0 < b1 && b0 < a1
}

func copyForward(src, dst *Bitmap, srcStart, dstStart, count int) {
	// Align dst to next word boundary for middle loop.
	doff := dstStart & 63
	if doff != 0 {
		headBits := min(count, 64-doff)
		val := readBits(src, srcStart, headBits)
		writeBits(dst, dstStart, headBits, val)
		srcStart += headBits
		dstStart += headBits
		count -= headBits
	}
	// Middle full words.
	for count >= 64 {
		wi := srcStart >> 6
		off := uint(srcStart & 63)
		var s0, s1 uint64
		if wi < len(src.words) {
			s0 = src.words[wi]
		}
		if wi+1 < len(src.words) {
			s1 = src.words[wi+1]
		}
		var dw uint64
		if off == 0 {
			dw = s0
		} else {
			dw = (s0 >> off) | (s1 << (64 - off))
		}
		dst.words[dstStart>>6] = dw
		srcStart += 64
		dstStart += 64
		count -= 64
	}
	// Tail.
	if count > 0 {
		val := readBits(src, srcStart, count)
		writeBits(dst, dstStart, count, val)
	}
}

func copyBackward(src, dst *Bitmap, srcStart, dstStart, count int) {
	srcEnd := srcStart + count
	dstEnd := dstStart + count

	// Align dstEnd to previous word boundary.
	tbits := dstEnd & 63
	if tbits != 0 {
		n := tbits
		val := readBits(src, srcEnd-n, n)
		writeBits(dst, dstEnd-n, n, val)
		srcEnd -= n
		dstEnd -= n
		count -= n
	}
	for count >= 64 {
		pos := srcEnd - 64
		wi := pos >> 6
		off := uint(pos & 63)
		var s0, s1 uint64
		if wi < len(src.words) {
			s0 = src.words[wi]
		}
		if wi+1 < len(src.words) {
			s1 = src.words[wi+1]
		}
		var dw uint64
		if off == 0 {
			dw = s0
		} else {
			dw = (s0 >> off) | (s1 << (64 - off))
		}
		dst.words[(dstEnd-64)>>6] = dw
		srcEnd -= 64
		dstEnd -= 64
		count -= 64
	}
	if count > 0 {
		val := readBits(src, srcStart, count)
		writeBits(dst, dstStart, count, val)
	}
}

func readBits(b *Bitmap, pos, n int) uint64 {
	wi := pos >> 6
	off := uint(pos & 63)
	var s0, s1 uint64
	if wi < len(b.words) {
		s0 = b.words[wi]
	}
	if wi+1 < len(b.words) {
		s1 = b.words[wi+1]
	}
	if n <= 64-int(off) {
		return (s0 >> off) & lowMask(uint(n))
	}
	x := (s0 >> off) | (s1 << (64 - off))
	return x & lowMask(uint(n))
}

func writeBits(b *Bitmap, pos, n int, val uint64) {
	wi := pos >> 6
	off := uint(pos & 63)

	if n <= 64-int(off) {
		mask := lowMask(uint(n)) << off
		b.words[wi] = (b.words[wi] &^ mask) | ((val << off) & mask)
		return
	}
	first := 64 - int(off)
	second := n - first

	mask1 := lowMask(uint(first)) << off
	b.words[wi] = (b.words[wi] &^ mask1) | ((val << off) & mask1)

	wi2 := wi + 1
	mask2 := lowMask(uint(second))
	b.words[wi2] = (b.words[wi2] &^ mask2) | ((val >> uint(first)) & mask2)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
