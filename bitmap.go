package btmp

import "math/bits"

const (
	WordBits         = 64
	WordShift        = 6            // log2(64), divide by 64 via >> 6
	IndexMask        = WordBits - 1 // for i & IndexMask
	WordMask  uint64 = ^uint64(0)   // 0xFFFFFFFFFFFFFFFF
)

// MaskFrom returns a mask with ones in [off, 63] and zeros in [0, off).
// If off >= 64, it returns 0.
func MaskFrom(off uint) uint64 {
	if off >= WordBits { return 0 }
	return WordMask << off
}

// MaskUpto returns a mask with ones in [0, off) and zeros in [off, 63].
// If off >= 64, it returns WordMask. If off == 0, it returns 0.
func MaskUpto(off uint) uint64 {
	if off >= WordBits { return WordMask }
	if off == 0 { return 0 }
	return (uint64(1) << off) - 1
}

// MaskRange returns a mask with ones in [lo, hi) and zeros elsewhere.
// If lo >= hi, it returns 0. Valid for 0 ≤ lo,hi ≤ 64.
func MaskRange(lo, hi uint) uint64 {
	if lo >= hi { return 0 }
	return MaskFrom(lo) & MaskUpto(hi)
}

// wordIndex converts a bit index to (wordIdx, bitOffset).
func wordIndex(i int) (w int, off int) { return i >> WordShift, i & IndexMask }

// Bitmap is a growable bitset backed by 64-bit words.
type Bitmap struct {
	words       []uint64
	lenBits     int
	lastWordIdx int    // index of last logical word; -1 if Len()==0
	tailMask    uint64 // mask for last logical word; 0 if Len()==0; WordMask if Len()%64==0
}

// New returns an empty bitmap sized for nbits bits (Len==nbits).
func New(nbits uint) *Bitmap {
	b := &Bitmap{
		words:   make([]uint64, (nbits+IndexMask)>>WordShift),
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

// --- internal methods ---

// computeCache recomputes cache fields from lenBits only.
func (b *Bitmap) computeCache() {
	if b.lenBits == 0 {
		b.lastWordIdx = -1
		b.tailMask = 0
		return
	}
	// ceil(lenBits/64) - 1
	b.lastWordIdx = int((b.lenBits+IndexMask)>>WordShift) - 1

	r := uint(b.lenBits) & IndexMask // bits used in last word, 0..63
	if r == 0 {
		b.tailMask = WordMask
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

// checkedEnd validates start/count, returns end, and panics if end > lenBits.
// No growth.
func (b *Bitmap) checkedEnd(start, count int) int {
	if start < 0 {
		panic("Bitmap: negative start")
	}
	if count < 0 {
		panic("Bitmap: negative count")
	}
	end := start + count
	if end < start {
		panic("Bitmap: integer overflow on end")
	}
	if end > b.lenBits {
		panic("Bitmap: out of bounds")
	}
	return end
}
