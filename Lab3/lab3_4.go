package main

import (
	"fmt"
	"os"
)

func main() {
	var element int
	var massiv [5]int
	for i := 0; i < 5; i++ {
		fmt.Print("Введите следующий элемент массива:")
		fmt.Fscan(os.Stdin, &element)
		massiv[i] = element
	}
	fmt.Println("Созданный вами массив:", massiv)
}
