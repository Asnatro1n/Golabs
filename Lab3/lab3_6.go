package main

import (
	"fmt"
	"os"
)

func main() {
	var elements int
	var stroka string
	fmt.Print("Введите кол-во строк в срезе:")
	fmt.Fscan(os.Stdin, &elements)
	var srez []string = make([]string, elements)
	var tekyshee int = 0
	for e := 0; e < elements; e++ {
		fmt.Print("Введите строку:")
		fmt.Fscan(os.Stdin, &stroka)
		srez[e] = stroka
	}
	for i := 0; i < elements-1; i++ {
		if len(srez[i]) >= len(srez[i+1]) {

		} else {
			tekyshee = i + 1
		}
	}
	fmt.Println("Самая длинная строка:", srez[tekyshee])
}
