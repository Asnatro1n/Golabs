package main

import (
	"fmt"
	mathut1ls "mathutils"
	"os"
)

func main() {
	var chislo int
	fmt.Print("Введите целое число для вычисления его факториала: ")
	fmt.Fscan(os.Stdin, &chislo)
	mathut1ls.Factorial(chislo)
	factor := mathut1ls.Factorial(chislo)
	fmt.Println("Факториал данного числа:", factor)
}
