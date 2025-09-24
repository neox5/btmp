package btmp

import "math/bits"

// Bitmap is a growable bitset backed by 64-bit words.
type Bitmap struct {
	words   []uint64
	lenBits int
}

// New returns an empty bitmap.
func New() *Bitmap {
	return &Bitmap{}
}

// NewWithCap returns an empty bitmap with capacity for capBits.
func NewWithCap(capBits int) *Bitmap {
	if capBits < 0 {
		panic("NewWithCap: negative capBits")
	}
	return &Bitmap{words: make([]uint64, wordsFor(capBits))}
}

// Len returns the logical length in bits.
func (b *Bitmap) Len() int { return b.lenBits }

// Words exposes the underlying words slice (length may exceed wordsFor(Len())).
func (b *Bitmap) Words() []uint64 { return b.words }

// Test reports whether bit i is set. Panics if i is out of [0, Len()).
func (b *Bitmap) Test(i int) bool {
	if i < 0 || i >= b.lenBits {
		panic("Test: index out of range")
	}
	w, off := wordIndex(i)
	return (b.words[w]>>off)&1 == 1
}

// Any reports whether any bit in [0, Len()) is set.
func (b *Bitmap) Any() bool {
	if b.lenBits == 0 {
		return false
	}
	limit := wordsFor(b.lenBits)
	for i := range limit-1 {
		if b.words[i] != 0 {
			return true
		}
	}
	// last word must respect tail mask
	last := limit - 1
	if b.lenBits&63 == 0 {
		return b.words[last] != 0
	}
	keep := uint(b.lenBits & 63)
	mask := (uint64(1) << keep) - 1
	return (b.words[last] & mask) != 0
}

// Count returns the number of set bits in [0, Len()).
func (b *Bitmap) Count() int {
	if b.lenBits == 0 {
		return 0
	}
	sum := 0
	limit := wordsFor(b.lenBits)
	for i := range limit-1 {
		sum += bits.OnesCount64(b.words[i])
	}
	last := limit - 1
	if b.lenBits&63 == 0 {
		return sum + bits.OnesCount64(b.words[last])
	}
	keep := uint(b.lenBits & 63)
	mask := (uint64(1) << keep) - 1
	return sum + bits.OnesCount64(b.words[last]&mask)
}

// NextSetBit returns the index of the first set bit >= from, or -1 if none.
// Panics if from is out of [0, Len()].
func (b *Bitmap) NextSetBit(from int) int {
	if from < 0 || from > b.lenBits {
		panic("NextSetBit: from out of range")
	}
	if from == b.lenBits {
		return -1
	}
	w, off := wordIndex(from)
	limit := wordsFor(b.lenBits)

	// first word
	word := b.words[w] & (^uint64(0) << off)
	if w == limit-1 && b.lenBits&63 != 0 {
		keep := uint(b.lenBits & 63)
		tailMask := (uint64(1) << keep) - 1
		word &= tailMask
	}
	if word != 0 {
		return (w << 6) + bits.TrailingZeros64(word)
	}

	// middle words
	for w = w + 1; w < limit-1; w++ {
		if b.words[w] != 0 {
			return (w << 6) + bits.TrailingZeros64(b.words[w])
		}
	}

	// last word
	if w == limit-1 {
		last := b.words[w]
		if b.lenBits&63 != 0 {
			keep := uint(b.lenBits & 63)
			last &= (uint64(1) << keep) - 1
		}
		if last != 0 {
			return (w << 6) + bits.TrailingZeros64(last)
		}
	}

	return -1
}

// --- internal helpers (signatures only) ---

// maskTail zeros bits >= Len() in the last word.
func maskTail(b *Bitmap) {
	if b.lenBits <= 0 || len(b.words) == 0 {
		return
	}
	if b.lenBits&63 == 0 {
		return
	}
	last := wordsFor(b.lenBits) - 1
	keep := uint(b.lenBits & 63)
	mask := (uint64(1) << keep) - 1
	b.words[last] &= mask
}

// finalize applies post-mutation cleanup and enforces invariants.
func finalize(b *Bitmap) { maskTail(b) }

// wordIndex converts a bit index to (wordIdx, bitOffset).
func wordIndex(i int) (w int, off uint) { return i >> 6, uint(i & 63) }

// helpers
func wordsFor(nbits int) int {
	if nbits <= 0 {
		return 0
	}
	return (nbits + 63) >> 6
}
