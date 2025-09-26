package btmp

import "slices"

// EnsureBits grows the logical length to at least n bits. Newly added bits are zero.
// No-op if n <= Len(). Panics if n < 0. Returns b.
//
// Invariant: after return, all bits >= Len() are zero.
func (b *Bitmap) EnsureBits(n int) *Bitmap {
	defer b.computeCache() // we mask (zeroing) the new words if added.
	if n < 0 {
		panic("EnsureBits: negative length")
	}
	if n <= b.lenBits {
		return b
	}
	need := int((n + IndexMask) >> WordShift)

	if need > len(b.words) {
		old := len(b.words)
		// Ensure capacity >= need, then reslice (set to new len) and zero the new words.
		b.words = slices.Grow(b.words, need-old)[:need]
		clear(b.words[old:]) // zeros new words
	}
	b.lenBits = n
	return b
}

// AddBits grows the logical length by n bits. Newly added bits are zero.
// No-op if n == 0. Panic if n < 0.
//
// Uses EnsireBits for final execution.
func (b *Bitmap) AddBits(n int) *Bitmap {
	if n < 0 {
		panic("AddBits: negative length")
	}
	if n == 0 {
		return b
	}

	return b.EnsureBits(b.lenBits + n)
}
