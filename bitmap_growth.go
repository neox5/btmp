package btmp

// EnsureBits grows the logical length to at least n bits. Newly added bits are zero.
// No-op if n <= Len(). Panics if n < 0. Returns b.
//
// Invariant: after return, all bits >= Len() are zero.
func (b *Bitmap) EnsureBits(n int) *Bitmap {
	defer finalize(b)
	if n < 0 {
		panic("EnsureBits: negative length")
	}
	if n <= b.lenBits {
		return b
	}
	need := int((n + wordMask) >> wordShift)
	if need > len(b.words) {
		ws := make([]uint64, need)
		copy(ws, b.words)
		b.words = ws
	}
	b.lenBits = n
	return b
}

// ReserveCap ensures capacity for at least n bits without changing Len().
// Panics if n < 0. Returns b.
func (b *Bitmap) ReserveCap(n int) *Bitmap {
	if n < 0 {
		panic("ReserveCap: negative")
	}
	need := int((n + wordMask) >> wordShift)
	if need > len(b.words) {
		ws := make([]uint64, need)
		copy(ws, b.words)
		b.words = ws
	}
	return b
}

// Trim reslices storage to the minimal number of words that hold Len() bits.
// Tail remains masked. Returns b.
func (b *Bitmap) Trim() *Bitmap {
	defer finalize(b)
	need := int((b.lenBits + wordMask) >> wordShift)
	if need < len(b.words) {
		b.words = b.words[:need]
	}
	return b
}
