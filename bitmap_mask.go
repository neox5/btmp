package btmp

// MaskFrom returns a mask with ones in [off, 63] and zeros in [0, off).
// If off >= 64, it returns 0.
func MaskFrom(off uint) uint64 {
	if off >= WordBits {
		return 0
	}
	return WordMask << off
}

// MaskUpto returns a mask with ones in [0, off) and zeros in [off, 63].
// If off >= 64, it returns WordMask. If off == 0, it returns 0.
func MaskUpto(off uint) uint64 {
	if off >= WordBits {
		return WordMask
	}
	if off == 0 {
		return 0
	}
	return (uint64(1) << off) - 1
}

// MaskRange returns a mask with ones in [lo, hi) and zeros elsewhere.
// If lo >= hi, it returns 0. Valid for 0 ≤ lo,hi ≤ 64.
func MaskRange(lo, hi uint) uint64 {
	if lo >= hi {
		return 0
	}
	return MaskFrom(lo) & MaskUpto(hi)
}
