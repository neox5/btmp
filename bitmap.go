// Package btmp provides a compact, growable bitmap optimized for fast range
// updates and overlap-safe copies, plus a zero-copy row-major 2D Grid view.
//
// Conventions:
//   - Length is in bits (Len).
//   - Storage is []uint64 words, exposed via Words() for read-only inspection.
//   - Ranges use (start, count).
//   - All operations are in-bounds only - no auto-growth.
//   - All mutating methods return *Bitmap for chaining.
//
// Invariant:
//   - After any public mutator returns, all bits at indexes >= Len() are zero,
//     even when count == 0.
package btmp

import "math/bits"

const (
	WordBits         = 64
	WordShift        = 6            // log2(64), divide by 64 via >> 6
	IndexMask        = WordBits - 1 // for i & IndexMask
	WordMask  uint64 = ^uint64(0)   // 0xFFFFFFFFFFFFFFFF
)

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
	b.computeCache()
	return b
}

// Len returns the logical length in bits.
func (b *Bitmap) Len() int { return b.lenBits }

// Words exposes the underlying words slice (length may exceed the logical need).
func (b *Bitmap) Words() []uint64 { return b.words }

// Test reports whether bit i is set. Panics if i is out of [0, Len()).
func (b *Bitmap) Test(i int) bool {
	validatePosition(i)
	b.validateInBounds(i)

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

// SetBit sets bit i to 1. Panics if i < 0 or i >= Len().
func (b *Bitmap) SetBit(i int) *Bitmap {
	validatePosition(i)
	b.validateInBounds(i)

	b.setBit(i)
	return b
}

// ClearBit sets bit i to 0. Panics if i < 0 or i >= Len().
func (b *Bitmap) ClearBit(i int) *Bitmap {
	validatePosition(i)
	b.validateInBounds(i)

	b.clearBit(i)
	return b
}

// FlipBit toggles bit i. Panics if i < 0 or i >= Len().
func (b *Bitmap) FlipBit(i int) *Bitmap {
	validatePosition(i)
	b.validateInBounds(i)

	b.flipBit(i)
	return b
}

// SetRange sets bits in [start, start+count) to 1. In-bounds only.
// Returns *Bitmap for chaining. Panics on negative inputs, overflow, or out-of-bounds.
func (b *Bitmap) SetRange(start, count int) *Bitmap {
	b.validateRange(start, count)

	b.setRange(start, count)
	return b
}

// ClearRange clears bits in [start, start+count) to 0. In-bounds only.
// Returns *Bitmap for chaining. Panics on negative inputs, overflow, or out-of-bounds.
func (b *Bitmap) ClearRange(start, count int) *Bitmap {
	b.validateRange(start, count)

	b.clearRange(start, count)
	return b
}

// CopyRange copies count bits from src[srcStart:] to dst[dstStart:].
// In-bounds only for both src and dst. Overlap-safe with memmove semantics.
// Returns *Bitmap for chaining. Panics on negative inputs, nil src, or out-of-bounds.
func (b *Bitmap) CopyRange(src *Bitmap, srcStart, dstStart, count int) *Bitmap {
	if src == nil {
		panic("CopyRange: nil source")
	}

	src.validateRange(srcStart, count)
	b.validateRange(dstStart, count)

	b.copyRange(src, srcStart, dstStart, count)
	return b
}

// EnsureBits grows the logical length to at least n bits. No-op if n <= Len().
// Returns *Bitmap for chaining. Panics if n < 0.
func (b *Bitmap) EnsureBits(n int) *Bitmap {
	validatePosition(n) // reuse for non-negative validation

	if n > b.lenBits {
		b.ensureBits(n)
		b.computeCache()
	}
	return b
}

// AddBits grows the logical length by n bits.
// Returns *Bitmap for chaining. Panics if n < 0.
func (b *Bitmap) AddBits(n int) *Bitmap {
	validateCount(n)

	if n > 0 {
		b.addBits(n)
		b.computeCache()
	}
	return b
}

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

// wordIndex converts a bit index to (wordIdx, bitOffset).
func wordIndex(i int) (w int, off int) {
	return i >> WordShift, i & IndexMask
}
