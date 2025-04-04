package main

import (
	"fmt"
	"os"
)

func main() {
	var numpeople int
	fmt.Print("Введите кол-во людей в карте: ")
	fmt.Fscan(os.Stdin, &numpeople)
	m := make(map[string]int, numpeople)
	for i := 1; i <= numpeople; i++ {
		var people string
		var age int
		fmt.Printf("%d: Введите имя: ", i)
		fmt.Fscan(os.Stdin, &people)
		fmt.Printf("%d: Введите возраст: ", i)
		fmt.Fscan(os.Stdin, &age)
		m[people] = age
	}
	m["Anton"] = 20
	fmt.Println(m)
}
