package main

import "fmt"

func main() {
	x := 1.0
	y := 2.0
	z := 1.0
	AverageZnach(x, y, z)
}

func AverageZnach(x float64, y float64, z float64) {
	fmt.Println("SredneZnach=", (x+y+z)/3)
}
