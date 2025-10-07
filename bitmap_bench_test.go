package btmp_test

import (
	"fmt"
	"testing"

	"github.com/neox5/btmp"
)

// BenchmarkSetRange tests setting ranges of different sizes
func BenchmarkSetRange(b *testing.B) {
	sizes := []struct {
		name  string
		total int
		start int
		count int
	}{
		{"SingleWord", 1000, 10, 50},
		{"CrossWord", 1000, 60, 20},
		{"ExactWord", 1000, 64, 64},
		{"Small", 10000, 100, 1000},
		{"Medium", 100000, 1000, 10000},
		{"Large", 1000000, 10000, 100000},
		{"FullBitmap", 100000, 0, 100000},
		{"Unaligned", 100000, 37, 28391},
	}

	for _, sz := range sizes {
		b.Run(sz.name, func(b *testing.B) {
			bm := btmp.New(uint(sz.total))
			b.ResetTimer()
			for b.Loop() {
				bm.SetRange(sz.start, sz.count)
				bm.ClearRange(sz.start, sz.count)
			}
		})
	}
}

// BenchmarkClearRange tests clearing ranges
func BenchmarkClearRange(b *testing.B) {
	sizes := []struct {
		name  string
		total int
		count int
	}{
		{"Small", 10000, 100},
		{"Medium", 100000, 10000},
		{"Large", 1000000, 100000},
	}

	for _, sz := range sizes {
		b.Run(sz.name, func(b *testing.B) {
			bm := btmp.New(uint(sz.total))
			bm.SetAll() // Start with all bits set
			b.ResetTimer()
			for b.Loop() {
				bm.ClearRange(10, sz.count)
				bm.SetRange(10, sz.count) // Reset for next iteration
			}
		})
	}
}

// BenchmarkAnyRange tests any-bit checking with different scenarios
func BenchmarkAnyRange(b *testing.B) {
	scenarios := []struct {
		name   string
		size   int
		setBit int // Position of set bit (-1 for none)
	}{
		{"Empty_Small", 1000, -1},
		{"Empty_Large", 1000000, -1},
		{"FirstBit_Small", 1000, 0},
		{"FirstBit_Large", 1000000, 0},
		{"EarlyBit", 1000000, 100},
		{"MiddleBit", 1000000, 500000},
		{"LateBit", 1000000, 999999},
		{"LastBit_Small", 1000, 999},
		{"LastBit_Large", 1000000, 999999},
	}

	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			bm := btmp.New(uint(sc.size))
			if sc.setBit >= 0 {
				bm.SetBit(sc.setBit)
			}
			b.ResetTimer()
			for b.Loop() {
				_ = bm.AnyRange(0, sc.size)
			}
		})
	}
}

// BenchmarkAnyRange_Density tests with different bit densities
func BenchmarkAnyRange_Density(b *testing.B) {
	size := 100000
	densities := []struct {
		name    string
		percent float64
	}{
		{"Empty", 0},
		{"Sparse_0.1%", 0.001},
		{"Sparse_1%", 0.01},
		{"Medium_10%", 0.1},
		{"Dense_50%", 0.5},
		{"Dense_90%", 0.9},
		{"Full", 1.0},
	}

	for _, d := range densities {
		b.Run(d.name, func(b *testing.B) {
			bm := btmp.New(uint(size))
			setBits := int(float64(size) * d.percent)
			if setBits > 0 {
				step := size / setBits
				for i := range setBits {
					bm.SetBit(i * step)
				}
			}
			b.ResetTimer()
			for b.Loop() {
				_ = bm.AnyRange(0, size)
			}
		})
	}
}

// BenchmarkAllRange tests all-bits checking
func BenchmarkAllRange(b *testing.B) {
	scenarios := []struct {
		name  string
		size  int
		setup func(*btmp.Bitmap)
	}{
		{"AllSet_Small", 1000, func(bm *btmp.Bitmap) {
			bm.SetAll()
		}},
		{"AllSet_Large", 100000, func(bm *btmp.Bitmap) {
			bm.SetAll()
		}},
		{"MissingFirst", 100000, func(bm *btmp.Bitmap) {
			bm.SetAll()
			bm.ClearBit(0)
		}},
		{"MissingMiddle", 100000, func(bm *btmp.Bitmap) {
			bm.SetAll()
			bm.ClearBit(50000)
		}},
		{"MissingLast", 100000, func(bm *btmp.Bitmap) {
			bm.SetAll()
			bm.ClearBit(99999)
		}},
		{"HalfSet", 100000, func(bm *btmp.Bitmap) {
			bm.SetRange(0, 50000)
		}},
	}

	for _, sc := range scenarios {
		b.Run(sc.name, func(b *testing.B) {
			bm := btmp.New(uint(sc.size))
			sc.setup(bm)
			b.ResetTimer()
			for b.Loop() {
				_ = bm.AllRange(0, sc.size)
			}
		})
	}
}

// BenchmarkCountRange tests bit counting
func BenchmarkCountRange(b *testing.B) {
	densities := []struct {
		name    string
		size    int
		setBits int
	}{
		{"Sparse_Small", 1000, 10},
		{"Sparse_Large", 100000, 100},
		{"Medium_Small", 1000, 500},
		{"Medium_Large", 100000, 50000},
		{"Dense_Small", 1000, 900},
		{"Dense_Large", 100000, 90000},
		{"Full_Small", 1000, 1000},
		{"Full_Large", 100000, 100000},
	}

	for _, d := range densities {
		b.Run(d.name, func(b *testing.B) {
			bm := btmp.New(uint(d.size))
			if d.setBits > 0 {
				step := d.size / d.setBits
				for i := 0; i < d.setBits; i++ {
					bm.SetBit(i * step)
				}
			}
			b.ResetTimer()
			for b.Loop() {
				_ = bm.CountRange(0, d.size)
			}
		})
	}
}

// BenchmarkWordBoundaries tests operations at word boundaries
func BenchmarkWordBoundaries(b *testing.B) {
	tests := []struct {
		name  string
		start int
		count int
	}{
		{"Within_First_Word", 5, 50},
		{"Cross_Single_Boundary", 60, 10},
		{"Exact_Word", 64, 64},
		{"Exact_Two_Words", 0, 128},
		{"Multiple_Words", 64, 1000},
		{"Unaligned_Start", 37, 283},
		{"Unaligned_Large", 17, 100000},
		{"Last_Partial", 50, 14}, // Within single word
		{"Cross_Many", 3, 500},   // Crosses multiple words
	}

	bm := btmp.New(200000)

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				bm.SetRange(tt.start, tt.count)
				bm.ClearRange(tt.start, tt.count)
			}
		})
	}
}

// BenchmarkCopyRange tests copying ranges
func BenchmarkCopyRange(b *testing.B) {
	tests := []struct {
		name     string
		size     int
		srcStart int
		dstStart int
		count    int
	}{
		{"Small_NoOverlap", 10000, 100, 5000, 1000},
		{"Large_NoOverlap", 100000, 1000, 50000, 10000},
		{"Small_Overlap_Forward", 10000, 100, 500, 1000},
		{"Small_Overlap_Backward", 10000, 500, 100, 1000},
		{"Large_Overlap_Forward", 100000, 1000, 5000, 10000},
		{"Large_Overlap_Backward", 100000, 5000, 1000, 10000},
		{"SamePosition", 100000, 1000, 1000, 10000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			bm := btmp.New(uint(tt.size))
			// Set pattern in source range
			for i := tt.srcStart; i < tt.srcStart+tt.count && i < tt.size; i += 3 {
				bm.SetBit(i)
			}
			b.ResetTimer()
			for b.Loop() {
				bm.CopyRange(bm, tt.srcStart, tt.dstStart, tt.count)
			}
		})
	}
}

// BenchmarkMoveRange tests moving ranges
func BenchmarkMoveRange(b *testing.B) {
	tests := []struct {
		name     string
		size     int
		srcStart int
		dstStart int
		count    int
	}{
		{"Small_NoOverlap", 10000, 100, 5000, 1000},
		{"Large_NoOverlap", 100000, 1000, 50000, 10000},
		{"Small_Overlap", 10000, 100, 500, 1000},
		{"Large_Overlap", 100000, 1000, 5000, 10000},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			bm := btmp.New(uint(tt.size))
			template := btmp.New(uint(tt.size))
			// Create template pattern
			for i := tt.srcStart; i < tt.srcStart+tt.count && i < tt.size; i += 3 {
				template.SetBit(i)
			}
			b.ResetTimer()
			for b.Loop() {
				bm.CopyRange(template, 0, 0, tt.size) // Reset bitmap
				bm.MoveRange(tt.srcStart, tt.dstStart, tt.count)
			}
		})
	}
}

// BenchmarkSetAll tests setting all bits
func BenchmarkSetAll(b *testing.B) {
	sizes := []int{64, 1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			bm := btmp.New(uint(size))
			b.ResetTimer()
			for b.Loop() {
				bm.SetAll()
				bm.ClearAll() // Reset for next iteration
			}
		})
	}
}

// BenchmarkClearAll tests clearing all bits
func BenchmarkClearAll(b *testing.B) {
	sizes := []int{64, 1000, 10000, 100000, 1000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Size_%d", size), func(b *testing.B) {
			bm := btmp.New(uint(size))
			b.ResetTimer()
			for b.Loop() {
				bm.SetAll()   // Prepare
				bm.ClearAll() // Measure
			}
		})
	}
}

// BenchmarkRangeSizes compares performance across different range sizes
func BenchmarkRangeSizes(b *testing.B) {
	bitmapSize := 1000000
	rangeSizes := []int{1, 10, 64, 100, 1000, 10000, 100000, 500000, 1000000}

	bm := btmp.New(uint(bitmapSize))

	for _, size := range rangeSizes {
		b.Run(fmt.Sprintf("SetRange_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				bm.SetRange(0, size)
				bm.ClearRange(0, size)
			}
		})
	}

	// Test with counting
	bm.SetAll()
	for _, size := range rangeSizes {
		b.Run(fmt.Sprintf("CountRange_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_ = bm.CountRange(0, size)
			}
		})
	}
}

// BenchmarkPatterns tests specific bit patterns
func BenchmarkPatterns(b *testing.B) {
	size := 100000
	patterns := []struct {
		name  string
		setup func(*btmp.Bitmap)
	}{
		{"Alternating", func(bm *btmp.Bitmap) {
			for i := 0; i < size; i += 2 {
				bm.SetBit(i)
			}
		}},
		{"Blocks_64", func(bm *btmp.Bitmap) {
			for i := 0; i < size; i += 128 {
				bm.SetRange(i, 64)
			}
		}},
		{"Sparse_1_in_100", func(bm *btmp.Bitmap) {
			for i := 0; i < size; i += 100 {
				bm.SetBit(i)
			}
		}},
		{"Dense_99_in_100", func(bm *btmp.Bitmap) {
			bm.SetAll()
			for i := 0; i < size; i += 100 {
				bm.ClearBit(i)
			}
		}},
	}

	for _, p := range patterns {
		bm := btmp.New(uint(size))
		p.setup(bm)

		b.Run(p.name+"_Count", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_ = bm.CountRange(0, size)
			}
		})

		b.Run(p.name+"_Any", func(b *testing.B) {
			b.ResetTimer()
			for b.Loop() {
				_ = bm.AnyRange(0, size)
			}
		})
	}
}
