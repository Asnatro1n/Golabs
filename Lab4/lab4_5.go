package main

import (
	"fmt"
	"os"
)

func main() {
	var value int
	var summa float32 = 0.0
	fmt.Print("Введите кол-во вводимых чисел: ")
	fmt.Fscan(os.Stdin, &value)
	//Chisla := make([]float32, value)
	for i := 0; i < value; i++ {
		var chislo float32
		fmt.Print("Введите число: ")
		fmt.Fscan(os.Stdin, &chislo)
		summa += chislo
	}
	fmt.Println("Сумма всех введённых чисел: ", summa)
}
