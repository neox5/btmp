package btmp

import "math/bits"

// test reports whether bit pos is set.
// Internal implementation - no validation.
func (b *Bitmap) test(pos int) bool {
	w, off := wordIndex(pos)
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

// anyRange reports whether any bit in [start, start+count) is set.
// Internal implementation - no validation.
func (b *Bitmap) anyRange(start, count int) bool {
	if count == 0 {
		return false
	}

	end := start + count
	w0, _ := wordIndex(start)
	w1, _ := wordIndex(end)

	headMask, tailMask := maskForRange(start, end)

	// Single word case
	if w0 == w1 {
		return (b.words[w0] & headMask) != 0
	}

	// Check head partial word
	if headMask != 0 && (b.words[w0]&headMask) != 0 {
		return true
	}

	// Check middle full words
	for w := w0 + 1; w < w1; w++ {
		if b.words[w] != 0 {
			return true
		}
	}

	// Check tail partial word
	if tailMask != 0 && (b.words[w1]&tailMask) != 0 {
		return true
	}

	return false
}

// allRange reports whether all bits in [start, start+count) are set.
// Internal implementation - no validation.
func (b *Bitmap) allRange(start, count int) bool {
	if count == 0 {
		return true // vacuously true for empty range
	}

	end := start + count
	w0, _ := wordIndex(start)
	w1, _ := wordIndex(end)

	headMask, tailMask := maskForRange(start, end)

	// Single word case
	if w0 == w1 {
		return (b.words[w0] & headMask) == headMask
	}

	// Check head partial word
	if headMask != 0 && (b.words[w0]&headMask) != headMask {
		return false
	}

	// Check middle full words
	for w := w0 + 1; w < w1; w++ {
		if b.words[w] != WordMask {
			return false
		}
	}

	// Check tail partial word
	if tailMask != 0 && (b.words[w1]&tailMask) != tailMask {
		return false
	}

	return true
}

// countRange returns the number of set bits in [start, start+count).
// Internal implementation - no validation.
func (b *Bitmap) countRange(start, count int) int {
	if count == 0 {
		return 0
	}

	end := start + count
	w0, _ := wordIndex(start)
	w1, _ := wordIndex(end)

	headMask, tailMask := maskForRange(start, end)

	// Single word case
	if w0 == w1 {
		return bits.OnesCount64(b.words[w0] & headMask)
	}

	sum := 0

	// Count head partial word
	if headMask != 0 {
		sum += bits.OnesCount64(b.words[w0] & headMask)
	}

	// Count middle full words
	for w := w0 + 1; w < w1; w++ {
		sum += bits.OnesCount64(b.words[w])
	}

	// Count tail partial word
	if tailMask != 0 {
		sum += bits.OnesCount64(b.words[w1] & tailMask)
	}

	return sum
}
