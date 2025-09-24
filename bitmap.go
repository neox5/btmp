package btmp

import "math/bits"

type Bitmap struct {
	words   []uint64
	lenBits int
}

func New() *Bitmap { return &Bitmap{} }

func NewWithCap(capBits int) *Bitmap {
	if capBits < 0 {
		panic("btmp: negative capBits")
	}
	return &Bitmap{words: make([]uint64, 0, wordsFor(capBits))}
}

func (b *Bitmap) Len() int { return b.lenBits }

func (b *Bitmap) Words() []uint64 { return b.words }

func (b *Bitmap) Test(i int) bool {
	if i < 0 || i >= b.lenBits {
		panic("btmp: Test index out of range")
	}
	w := i >> 6
	return ((b.words[w] >> uint(i&63)) & 1) == 1
}

func (b *Bitmap) Any() bool {
	nw := wordsFor(b.lenBits)
	for i := range nw {
		if b.words[i] != 0 {
			return true
		}
	}
	return false
}

func (b *Bitmap) Count() int {
	nw := wordsFor(b.lenBits)
	sum := 0
	for i := range nw {
		sum += bits.OnesCount64(b.words[i])
	}
	return sum
}

func (b *Bitmap) NextSetBit(from int) int {
	if from < 0 || from > b.lenBits {
		panic("btmp: NextSetBit from out of range")
	}
	if from == b.lenBits {
		return -1
	}
	ws := from >> 6
	off := uint(from & 63)

	nw := wordsFor(b.lenBits)
	if ws >= nw {
		return -1
	}
	w := b.words[ws] & (^uint64(0) << off)
	if w != 0 {
		return (ws << 6) + bits.TrailingZeros64(w)
	}
	for i := ws + 1; i < nw; i++ {
		if b.words[i] != 0 {
			return (i << 6) + bits.TrailingZeros64(b.words[i])
		}
	}
	return -1
}
