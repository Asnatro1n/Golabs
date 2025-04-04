package main

import (
	"fmt"
	"os"
)

func main() {
	var chislo float64
	fmt.Print("Введите число: ")
	fmt.Fscan(os.Stdin, &chislo)
	KakoeChislo(chislo)
}

func KakoeChislo(x float64) {
	if x > 0 {
		fmt.Print("Positive")
	} else if x < 0 {
		fmt.Print("Negative")
	} else {
		fmt.Print("Zero")
	}
}
