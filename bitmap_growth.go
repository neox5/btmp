package btmp

import "slices"

// EnsureBits grows the logical length to at least n bits. Newly added bits are zero.
// No-op if n <= Len(). Panics if n < 0. Returns b.
//
// Invariant: after return, all bits >= Len() are zero.
func (b *Bitmap) EnsureBits(n int) *Bitmap {
	defer b.finalize()
	if n < 0 {
		panic("EnsureBits: negative length")
	}
	if n <= b.lenBits {
		return b
	}
	need := int((n + indexMask) >> wordShift)

	if need > len(b.words) {
		old := len(b.words)
		// Ensure capacity >= need, then reslice (set to new len) and zero the new words.
		b.words = slices.Grow(b.words, need-old)[:need]
		clear(b.words[old:]) // zeros new words
	}
	b.lenBits = n
	return b
}

// ReserveCap ensures capacity for at least n bits without changing Len().
// Panics if n < 0. Returns b.
func (b *Bitmap) ReserveCap(n int) *Bitmap {
	defer b.finalize()
	if n < 0 {
		panic("ReserveCap: negative")
	}
	needWords := int((n + indexMask) >> wordShift)
	if needWords > cap(b.words) {
		// Grow capacity to at least needWords; length stays the same.
		b.words = slices.Grow(b.words, needWords-cap(b.words))
	}
	return b
}

// Truncate reslices storage to the minimal number of words for Len().
// Capacity unchanged.
func (b *Bitmap) Truncate() *Bitmap {
	defer b.finalize()
	need := int((b.lenBits + indexMask) >> wordShift)
	if need < len(b.words) {
		b.words = b.words[:need]
	}
	return b
}

// Clip drops excess capacity after Truncate.
func (b *Bitmap) Clip() *Bitmap {
	defer b.finalize()
	b.words = slices.Clip(b.words)
	return b
}
