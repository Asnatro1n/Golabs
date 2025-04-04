package main

import (
	"fmt"
	"os"
)

func main() {
	var stroka string
	fmt.Print("Введите строку: ")
	fmt.Fscan(os.Stdin, &stroka)
	DlinaStroki(stroka)
}

func DlinaStroki(x string) {
	fmt.Println(len(x))
}
