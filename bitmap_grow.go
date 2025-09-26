package btmp

import "slices"

// ensureBits grows the logical length to at least n bits without validation.
// Internal implementation - no bounds checking, no finalization.
// Caller must ensure n >= 0 and handle finalization.
func (b *Bitmap) ensureBits(n int) {
	if n <= b.lenBits {
		return
	}
	need := (n + IndexMask) >> WordShift

	if need > len(b.words) {
		old := len(b.words)
		// Ensure capacity >= need, then reslice and zero new words
		b.words = slices.Grow(b.words, need-old)[:need]
		clear(b.words[old:])
	}
	b.lenBits = n
}

// addBits grows the logical length by n bits without validation.
// Internal implementation - no bounds checking, no finalization.
// Caller must ensure n >= 0 and handle finalization.
func (b *Bitmap) addBits(n int) {
	if n == 0 {
		return
	}
	b.ensureBits(b.lenBits + n)
}
