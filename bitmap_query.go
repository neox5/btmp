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
