package btmp

import "math/bits"

const (
	wordMask  = 64 - 1 // mask 0b111111
	wordShift = 6      // log2(64), used for division by 64 using rightshift >> 6
)

// Bitmap is a growable bitset backed by 64-bit words.
type Bitmap struct {
	words   []uint64
	lenBits int
	// cache
	wordCount int    // number of words covering Len(); (lenBits + wordMask) >> wordShift
	tailMask  uint64 // mask for last used word; 0 if lenBits==0; ^0 if lenBits%64==0
}

// New returns an empty bitmap sized for nbits bits (Len==nbits).
func New(nbits uint) *Bitmap {
	b := &Bitmap{
		words:   make([]uint64, (nbits+wordMask)>>wordShift),
		lenBits: int(nbits),
	}
	b.computeCache()
	return b
}

// computeCache recomputes cache fields from lenBits only.
func (b *Bitmap) computeCache() {
	if b.lenBits == 0 {
		b.wordCount = 0
		b.tailMask = 0
		return
	}
	b.wordCount = int((b.lenBits + wordMask) >> wordShift)

	r := uint(b.lenBits) & wordMask // 0..63
	if r == 0 {
		b.tailMask = ^uint64(0)
		return
	}
	b.tailMask = (uint64(1) << r) - 1
}

// Len returns the logical length in bits.
func (b *Bitmap) Len() int { return b.lenBits }

// Words exposes the underlying words slice (length may exceed the logical need).
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
	// full words except the last
	lastIdx := b.wordCount - 1
	for i := range lastIdx {
		if b.words[i] != 0 {
			return true
		}
	}
	// masked last word
	return (b.words[lastIdx] & b.tailMask) != 0
}

// Count returns the number of set bits in [0, Len()).
func (b *Bitmap) Count() int {
	if b.lenBits == 0 {
		return 0
	}
	sum := 0
	lastIdx := b.wordCount - 1
	for i := range lastIdx {
		sum += bits.OnesCount64(b.words[i])
	}
	return sum + bits.OnesCount64(b.words[lastIdx]&b.tailMask)
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

	limit := b.wordCount

	// first word
	word := b.words[w] & (^uint64(0) << off)
	// if this first word is also the last logical word, mask its tail
	if w == limit-1 {
		word &= b.tailMask
	}
	if word != 0 {
		return (w << wordShift) + bits.TrailingZeros64(word)
	}

	// middle full words (if any)
	for w = w + 1; w < limit-1; w++ {
		if b.words[w] != 0 {
			return (w << wordShift) + bits.TrailingZeros64(b.words[w])
		}
	}

	// last word
	if w == limit-1 {
		last := b.words[w] & b.tailMask
		if last != 0 {
			return (w << wordShift) + bits.TrailingZeros64(last)
		}
	}

	return -1
}

// --- internal helpers ---

// maskTail zeros bits >= Len() in the last word using cached tailMask.
func maskTail(b *Bitmap) {
	if b.lenBits <= 0 || len(b.words) == 0 || b.wordCount == 0 {
		return
	}
	last := b.wordCount - 1
	b.words[last] &= b.tailMask
}

// finalize recomputes cache then applies tail masking to enforce invariants.
func finalize(b *Bitmap) {
	b.computeCache()
	maskTail(b)
}

// wordIndex converts a bit index to (wordIdx, bitOffset).
func wordIndex(i int) (w int, off uint) { return i >> wordShift, uint(i & wordMask) }
