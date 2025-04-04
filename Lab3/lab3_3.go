package main

import (
	"fmt"
	"os"
	"stringutils"
)

func main() {
	var stroka string
	fmt.Print("Введите строку для переворота:")
	fmt.Fscan(os.Stdin, &stroka)
	reverse := stringutils.Perevorot(stroka)
	fmt.Println("Перевёрнутая строка:", reverse)
}
