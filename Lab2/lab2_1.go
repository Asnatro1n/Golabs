package main

import (
	"fmt"
	"os"
)

func main() {
	var chislo int
	fmt.Print("Введите целое число: ")
	fmt.Fscan(os.Stdin, &chislo)
	if chislo%2 == 0 {
		fmt.Print("Введённое число чётное")
	} else {
		fmt.Print("Введённое число нечётное")
	}
}
