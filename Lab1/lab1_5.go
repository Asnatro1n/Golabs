package main

import "fmt"

func main() {
	var a float64 = 6.52
	var b float64 = 17.25
	PlusMinus(a, b)
}

func PlusMinus(x float64, y float64) {
	var plus = x + y
	var minus = x - y
	fmt.Println("Summa=", plus)
	fmt.Println("Raznost=", minus)
}
