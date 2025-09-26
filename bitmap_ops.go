package btmp

// and performs bitwise AND with other bitmap.
// Internal implementation - no validation, no finalization.
// Assumes same length and sufficient capacity.
func (b *Bitmap) and(other *Bitmap) {
	if b.lenBits == 0 {
		return
	}

	// Process full words
	for i := range b.lastWordIdx {
		b.words[i] &= other.words[i]
	}

	// Process last partial word with proper masking
	b.words[b.lastWordIdx] = (b.words[b.lastWordIdx] & other.words[b.lastWordIdx]) & b.tailMask
}

// or performs bitwise OR with other bitmap.
// Internal implementation - no validation, no finalization.
// Assumes same length and sufficient capacity.
func (b *Bitmap) or(other *Bitmap) {
	if b.lenBits == 0 {
		return
	}

	// Process full words
	for i := range b.lastWordIdx {
		b.words[i] |= other.words[i]
	}

	// Process last partial word with proper masking
	b.words[b.lastWordIdx] = (b.words[b.lastWordIdx] | other.words[b.lastWordIdx]) & b.tailMask
}

// xor performs bitwise XOR with other bitmap.
// Internal implementation - no validation, no finalization.
// Assumes same length and sufficient capacity.
func (b *Bitmap) xor(other *Bitmap) {
	if b.lenBits == 0 {
		return
	}

	// Process full words
	for i := range b.lastWordIdx {
		b.words[i] ^= other.words[i]
	}

	// Process last partial word with proper masking
	b.words[b.lastWordIdx] = (b.words[b.lastWordIdx] ^ other.words[b.lastWordIdx]) & b.tailMask
}

// not performs bitwise NOT (flips all bits in [0, Len())).
// Internal implementation - no validation, no finalization.
func (b *Bitmap) not() {
	if b.lenBits == 0 {
		return
	}

	// Process full words
	for i := range b.lastWordIdx {
		b.words[i] = ^b.words[i]
	}

	// Process last partial word with proper masking
	b.words[b.lastWordIdx] = (^b.words[b.lastWordIdx]) & b.tailMask
}
