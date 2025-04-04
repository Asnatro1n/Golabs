package main

import (
	"fmt"
	"os"
)

func main() {
	var FirstNumber int
	var SecondNumber int
	fmt.Print("Первое целое число: ")
	fmt.Fscan(os.Stdin, &FirstNumber)
	fmt.Print("Второе целое число: ")
	fmt.Fscan(os.Stdin, &SecondNumber)
	SredneZnach2(FirstNumber, SecondNumber)
}

func SredneZnach2(a int, b int) {
	var sredne float64
	sredne = float64(a+b) / 2
	fmt.Println("Среднее значение: ", sredne)
}
