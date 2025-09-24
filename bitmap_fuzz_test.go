package btmp_test

import (
	"math/rand"
	"testing"

	btmp "github.com/neox5/btmp"
)

func FuzzBitmapAgainstRef(f *testing.F) {
	// Seed cases: (seed, opsN, maxLen)
	f.Add(int64(1), int64(200), int64(4096))
	f.Add(int64(42), int64(400), int64(2048))
	f.Add(int64(7), int64(50), int64(512))

	f.Fuzz(func(t *testing.T, seed, opsN, maxLen int64) {
		if opsN <= 0 {
			opsN = 200
		}
		if maxLen <= 0 {
			maxLen = 4096
		}
		if maxLen > 1<<20 {
			maxLen = 1 << 20
		}

		rng := rand.New(rand.NewSource(seed))
		b := btmp.New()
		r := newRef(0)

		for range int(opsN) {
			switch rng.Intn(3) {
			case 0: // set
				s := rng.Intn(int(maxLen))
				c := rng.Intn(256)
				b = b.SetRange(s, c)
				r.ensure(s + c)
				r.setRange(s, c)

			case 1: // clear within bounds
				n := r.len
				if n == 0 {
					continue
				}
				s := rng.Intn(n)
				c := rng.Intn(n - s)
				b = b.ClearRange(s, c)
				r.clearRange(s, c)

			case 2: // copy self with auto-grow on dst
				n := r.len
				if n == 0 {
					continue
				}
				ss := rng.Intn(n)
				cc := rng.Intn(n - ss)
				ds := rng.Intn(int(maxLen))
				b = b.CopyRange(b, ss, ds, cc)
				// ensure reference size
				end := ds + cc
				if end > r.len {
					r.ensure(end)
				}
				r.copyRange(r, ss, ds, cc)
			}
		}

		eqBitmapRef(t, b, r)
	})
}
