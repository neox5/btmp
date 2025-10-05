package main

import (
	"fmt"

	"github.com/neox5/btmp"
)

func main() {
	fmt.Println("=== Empty Grid ===")
	fmt.Println()
	g0 := btmp.NewGridWithSize(0, 0)
	fmt.Println(g0.Print())
	fmt.Println("(empty output expected)")
	fmt.Println()

	fmt.Println("=== Small Grid (3x5) ===")
	fmt.Println()
	g1 := btmp.NewGridWithSize(3, 5)
	g1.SetRect(0, 1, 1, 1) // (0,1)
	g1.SetRect(1, 3, 1, 1) // (1,3)
	fmt.Println(g1.Print())
	fmt.Println()

	fmt.Println("=== Single Row ===")
	fmt.Println()
	g2 := btmp.NewGridWithSize(1, 8)
	g2.SetRect(0, 0, 1, 1) // (0,0)
	g2.SetRect(0, 3, 1, 1) // (0,3)
	g2.SetRect(0, 7, 1, 1) // (0,7)
	fmt.Println(g2.Print())
	fmt.Println()

	fmt.Println("=== Single Column ===")
	fmt.Println()
	g3 := btmp.NewGridWithSize(8, 1)
	g3.SetRect(1, 0, 1, 1) // (1,0)
	g3.SetRect(4, 0, 1, 1) // (4,0)
	g3.SetRect(7, 0, 1, 1) // (7,0)
	fmt.Println(g3.Print())
	fmt.Println()

	fmt.Println("=== 10x10 Grid with Pattern ===")
	fmt.Println()
	g4 := btmp.NewGridWithSize(10, 10)
	// Diagonal
	for i := range 10 {
		g4.SetRect(i, i, 1, 1)
	}
	// Box at (3,3,4,4)
	g4.SetRect(3, 3, 1, 4) // top
	g4.SetRect(6, 3, 1, 4) // bottom
	g4.SetRect(4, 3, 2, 1) // left
	g4.SetRect(4, 6, 2, 1) // right
	fmt.Println(g4.Print())
	fmt.Println()

	fmt.Println("=== Large Grid (15x20) ===")
	fmt.Println()
	g5 := btmp.NewGridWithSize(15, 20)
	// Fill some rectangles
	g5.SetRect(2, 2, 3, 5)
	g5.SetRect(5, 10, 4, 8)
	g5.SetRect(11, 1, 2, 3)
	g5.SetRect(1, 15, 2, 4)
	fmt.Println(g5.Print())
	fmt.Println()

	fmt.Println("=== Full Grid ===")
	fmt.Println()
	g6 := btmp.NewGridWithSize(4, 6)
	g6.SetRect(0, 0, 4, 6) // Fill entire grid
	fmt.Println(g6.Print())
}
