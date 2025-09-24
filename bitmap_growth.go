package btmp

// EnsureBits grows the logical length to at least n bits.
// Newly added bits are zero. May grow capacity. Never shrinks.
// No-op if n <= Len(). Panics if n < 0. Returns b for chaining.
func (b *Bitmap) EnsureBits(n int) *Bitmap {
	if n < 0 {
		panic("btmp: EnsureBits negative")
	}
	if n > b.lenBits {
		ensureLen(b, n)
	}
	return b
}

// ReserveCap ensures capacity for at least n bits without changing Len().
// Panics if n < 0. Returns b for chaining.
func (b *Bitmap) ReserveCap(n int) *Bitmap {
	if n < 0 {
		panic("btmp: ReserveCap negative")
	}
	needW := wordsFor(n)
	if needW <= cap(b.words) {
		return b
	}
	newCap := growCap(len(b.words), needW)
	nb := make([]uint64, len(b.words), newCap)
	copy(nb, b.words)
	b.words = nb
	return b
}

// Trim reslices storage to the minimal number of words that hold Len() bits.
// Capacity may remain >= length. Tail remains masked. Returns b for chaining.
func (b *Bitmap) Trim() *Bitmap {
	want := wordsFor(b.lenBits)
	if want < len(b.words) {
		b.words = b.words[:want]
	}
	return b
}

/*** internal growth ***/

const (
	growMinWords     = 8
	growSoftCapWords = 1 << 20 // words; ~8 MiB of uint64s
)

func ensureLen(b *Bitmap, needBits int) {
	needW := wordsFor(needBits)
	oldW := len(b.words)
	if needW > oldW {
		if needW <= cap(b.words) {
			old := b.words
			b.words = old[:needW]
			for i := oldW; i < needW; i++ {
				b.words[i] = 0
			}
		} else {
			newCap := growCap(oldW, needW)
			nb := make([]uint64, needW, newCap)
			copy(nb, b.words)
			b.words = nb
		}
	}
	b.lenBits = needBits
	maskTail(b)
}

func growCap(cur, need int) int {
	if cur < 0 {
		cur = 0
	}
	var inc int
	if cur >= growSoftCapWords {
		inc = cur / 5 // ~1.2x
	} else {
		inc = cur / 2 // ~1.5x
	}
	if inc < growMinWords {
		inc = growMinWords
	}
	return max(cur+inc, need)
}
