package main

import "fmt"

func main() {
	x := 5
	y := 2
	plus := x + y
	minus := x - y
	multiplication := x * y
	var division float64 = float64(x) / float64(y)
	fmt.Println(plus)
	fmt.Println(minus)
	fmt.Println(multiplication)
	fmt.Println(division)
}
