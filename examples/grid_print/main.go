package main

import (
	"fmt"

	"github.com/neox5/btmp"
)

func main() {
	fmt.Println("=== Empty Grid ===\n")
	g0 := btmp.NewGridWithSize(0, 0)
	fmt.Println(g0.Print())
	fmt.Println("(empty output expected)\n")

	fmt.Println("=== Small Grid (5x3) ===\n")
	g1 := btmp.NewGridWithSize(5, 3)
	g1.SetRect(1, 0, 1, 1) // (1,0)
	g1.SetRect(3, 1, 1, 1) // (3,1)
	fmt.Println(g1.Print())
	fmt.Println()

	fmt.Println("=== Single Row ===\n")
	g2 := btmp.NewGridWithSize(8, 1)
	g2.SetRect(0, 0, 1, 1) // (0,0)
	g2.SetRect(3, 0, 1, 1) // (3,0)
	g2.SetRect(7, 0, 1, 1) // (7,0)
	fmt.Println(g2.Print())
	fmt.Println()

	fmt.Println("=== Single Column ===\n")
	g3 := btmp.NewGridWithSize(1, 8)
	g3.SetRect(0, 1, 1, 1) // (0,1)
	g3.SetRect(0, 4, 1, 1) // (0,4)
	g3.SetRect(0, 7, 1, 1) // (0,7)
	fmt.Println(g3.Print())
	fmt.Println()

	fmt.Println("=== 10x10 Grid with Pattern ===\n")
	g4 := btmp.NewGridWithSize(10, 10)
	// Diagonal
	for i := range 10 {
		g4.SetRect(i, i, 1, 1)
	}
	// Box at (3,3,4,4)
	g4.SetRect(3, 3, 4, 1) // top
	g4.SetRect(3, 6, 4, 1) // bottom
	g4.SetRect(3, 4, 1, 2) // left
	g4.SetRect(6, 4, 1, 2) // right
	fmt.Println(g4.Print())
	fmt.Println()

	fmt.Println("=== Large Grid (20x15) ===\n")
	g5 := btmp.NewGridWithSize(20, 15)
	// Fill some rectangles
	g5.SetRect(2, 2, 5, 3)
	g5.SetRect(10, 5, 8, 4)
	g5.SetRect(1, 11, 3, 2)
	g5.SetRect(15, 1, 4, 2)
	fmt.Println(g5.Print())
	fmt.Println()

	fmt.Println("=== Full Grid ===\n")
	g6 := btmp.NewGridWithSize(6, 4)
	g6.SetRect(0, 0, 6, 4) // Fill entire grid
	fmt.Println(g6.Print())
}
