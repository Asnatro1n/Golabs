package main

import (
	"fmt"
	"os"
)

func main() {
	var value int
	//var granica int
	fmt.Print("Введите кол-во вводимых чисел: ")
	fmt.Fscan(os.Stdin, &value)
	massive := make([]int, value)
	for i := 0; i < value; i++ {
		var chislo int
		fmt.Print("Введите число: ")
		fmt.Fscan(os.Stdin, &chislo)
		massive[i] = chislo
	}
	/*if value%2 == 1 {
		granica = (value - 1) / 2
	} else {
		granica = value / 2
	}
	value--
	for i := 0; i < granica; i++ {
		buffer := massive[i]
		massive[i] = massive[value]
		massive[value] = buffer
		value--
	}
	for i := 0; i < len(massive); i++ {
		fmt.Println(massive[i])
	}*/
	for i := len(massive) - 1; i >= 0; i-- {
		fmt.Println(massive[i])
	}
}
