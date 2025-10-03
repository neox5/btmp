package main

import (
	"fmt"

	"github.com/neox5/btmp"
)

func main() {
	fmt.Println("=== Basic Binary Printing ===\n")

	// Small bitmap
	b1 := btmp.New(16)
	b1.SetBit(0)
	b1.SetBit(2)
	b1.SetBit(7)
	b1.SetBit(15)
	fmt.Printf("16 bits with pattern:\n%s\n\n", b1.Print())

	// Range printing
	b2 := btmp.New(100)
	b2.SetRange(10, 20)
	fmt.Printf("Bits [10, 30) set:\n%s\n\n", b2.PrintRange(0, 40))

	fmt.Println("=== Hexadecimal Printing ===\n")

	// Hex format
	b3 := btmp.New(16)
	b3.SetBits(0, 16, 0xABCD)
	fmt.Printf("0xABCD in binary:\n%s\n", b3.Print())
	fmt.Printf("0xABCD in hex:\n%s\n\n", b3.PrintFormat(16, false, 0, ""))

	// Hex with partial word
	b4 := btmp.New(12)
	b4.SetBits(0, 12, 0xABC)
	fmt.Printf("0xABC (12 bits) in binary:\n%s\n", b4.Print())
	fmt.Printf("0xABC (12 bits) in hex:\n%s\n\n", b4.PrintFormat(16, false, 0, ""))

	fmt.Println("=== Grouped Binary ===\n")

	// Binary grouped by 4
	b5 := btmp.New(16)
	b5.SetAll()
	fmt.Printf("16 bits all set, grouped by 4:\n%s\n\n", b5.PrintFormat(2, true, 4, "_"))

	// Binary grouped by 8
	b6 := btmp.New(32)
	b6.SetBits(0, 32, 0xDEADBEEF)
	fmt.Printf("0xDEADBEEF grouped by 8 bits:\n%s\n\n", b6.PrintFormat(2, true, 8, " "))

	fmt.Println("=== Grouped Hexadecimal ===\n")

	// Hex grouped by 2
	b7 := btmp.New(32)
	b7.SetBits(0, 32, 0xDEADBEEF)
	fmt.Printf("0xDEADBEEF in hex:\n%s\n", b7.PrintFormat(16, false, 0, ""))
	fmt.Printf("0xDEADBEEF in hex, grouped by 2:\n%s\n\n", b7.PrintFormat(16, true, 2, " "))

	// Hex grouped by 4
	b8 := btmp.New(64)
	b8.SetAll()
	fmt.Printf("64 bits all set, hex grouped by 4:\n%s\n\n", b8.PrintFormat(16, true, 4, "_"))

	fmt.Println("=== Large Ranges (> 64 bits) ===\n")

	// Large binary
	b9 := btmp.New(130)
	b9.SetBit(0)
	b9.SetBit(64)
	b9.SetBit(128)
	fmt.Printf("130 bits with bits 0, 64, 128 set:\n%s\n\n", b9.Print())

	// Large binary grouped
	b10 := btmp.New(128)
	b10.SetAll()
	fmt.Printf("128 bits all set, grouped by 8:\n%s\n\n", b10.PrintFormat(2, true, 8, "_"))

	// Large hex
	b11 := btmp.New(96)
	b11.SetAll()
	fmt.Printf("96 bits all set in hex:\n%s\n", b11.PrintFormat(16, false, 0, ""))
	fmt.Printf("96 bits all set in hex, grouped by 4:\n%s\n\n", b11.PrintFormat(16, true, 4, " "))

	fmt.Println("=== Range Formatting ===\n")

	// Extract and format specific range
	b12 := btmp.New(200)
	b12.SetRange(50, 20)
	fmt.Printf("Bits [50, 70) in range [40, 80):\n%s\n\n", b12.PrintRangeFormat(40, 40, 2, true, 4, "_"))

	// Hex range
	b13 := btmp.New(100)
	b13.SetRange(16, 32)
	fmt.Printf("Bits [16, 48) as hex, grouped by 2:\n%s\n", b13.PrintRangeFormat(16, 32, 16, true, 2, " "))
}
