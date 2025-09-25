package btmp

import "math/bits"

const (
	wordBits         = 64
	wordShift        = 6            // log2(64), divide by 64 via >> 6
	indexMask        = wordBits - 1 // for i & indexMask
	wordMask  uint64 = ^uint64(0)   // 0xFFFFFFFFFFFFFFFF
)

// MaskFrom returns a 64-bit mask with ones in [off, 63] and zeros in [0, off).
// off ∈ [0, 64]. off==0 → ^uint64(0); off==64 → 0.
func MaskFrom(off uint) uint64 { return wordMask << off }

// MaskUpto returns a 64-bit mask with ones in [0, off) and zeros in [off, 63].
// off ∈ [0, 64]. off==0 → 0; off==64 → ^uint64(0).
func MaskUpto(off uint) uint64 { return (uint64(1) << off) - 1 }

// wordIndex converts a bit index to (wordIdx, bitOffset).
func wordIndex(i int) (w int, off uint) { return i >> wordShift, uint(i & indexMask) }

// Bitmap is a growable bitset backed by 64-bit words.
type Bitmap struct {
	words       []uint64
	lenBits     int
	lastWordIdx int    // index of last logical word; -1 if Len()==0
	tailMask    uint64 // mask for last logical word; 0 if Len()==0; wordMask if Len()%64==0
}

// New returns an empty bitmap sized for nbits bits (Len==nbits).
func New(nbits uint) *Bitmap {
	b := &Bitmap{
		words:   make([]uint64, (nbits+indexMask)>>wordShift),
		lenBits: int(nbits),
	}
	b.finalize()
	return b
}

// finalize recomputes cache then applies tail masking to enforce invariants.
func (b *Bitmap) finalize() {
	b.computeCache()
	b.maskTail()
}

// --- Public Bitmap API ---

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
	for i := range b.lastWordIdx {
		if b.words[i] != 0 {
			return true
		}
	}
	// masked last word
	return (b.words[b.lastWordIdx] & b.tailMask) != 0
}

// Count returns the number of set bits in [0, Len()).
func (b *Bitmap) Count() int {
	if b.lenBits == 0 {
		return 0
	}
	sum := 0
	for i := range b.lastWordIdx {
		sum += bits.OnesCount64(b.words[i])
	}
	return sum + bits.OnesCount64(b.words[b.lastWordIdx]&b.tailMask)
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

	// first word: keep bits from off..63, also tail-mask if it is the last
	first := b.words[w] & MaskFrom(off)
	if w == b.lastWordIdx {
		first &= b.tailMask
	}
	if first != 0 {
		return (w << wordShift) + bits.TrailingZeros64(first)
	}

	// middle full words (if any)
	for w = w + 1; w < b.lastWordIdx; w++ {
		if x := b.words[w]; x != 0 {
			return (w << wordShift) + bits.TrailingZeros64(x)
		}
	}

	// last word
	if w == b.lastWordIdx {
		if x := b.words[w] & b.tailMask; x != 0 {
			return (w << wordShift) + bits.TrailingZeros64(x)
		}
	}
	return -1
}

// --- internal methods ---

// computeCache recomputes cache fields from lenBits only.
func (b *Bitmap) computeCache() {
	if b.lenBits == 0 {
		b.lastWordIdx = -1
		b.tailMask = 0
		return
	}
	// ceil(lenBits/64) - 1
	b.lastWordIdx = int((b.lenBits+indexMask)>>wordShift) - 1

	r := uint(b.lenBits) & indexMask // bits used in last word, 0..63
	if r == 0 {
		b.tailMask = wordMask
		return
	}
	b.tailMask = MaskUpto(r)
}

// maskTail zeros bits >= Len() in the last word using cached tailMask.
func (b *Bitmap) maskTail() {
	if b.lenBits <= 0 || len(b.words) == 0 || b.lastWordIdx < 0 {
		return
	}
	b.words[b.lastWordIdx] &= b.tailMask
}
