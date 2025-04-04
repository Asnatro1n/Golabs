package main

import (
	"fmt"
	"os"
)

func main() {
	sum := 0
	var average float64
	size := 3
	m := make(map[string]int, size)
	for i := 0; i < size; i++ {
		var name string
		var age int
		fmt.Print("Введите имя: ")
		fmt.Fscan(os.Stdin, &name)
		fmt.Print("Введите возраст: ")
		fmt.Fscan(os.Stdin, &age)
		m[name] = age
		sum += age
	}
	average = Average(m, size, sum)
	fmt.Println(m, " - Average age: ", average)
}

func Average(m map[string]int, size int, sum int) float64 {
	return float64(sum) / float64(size)
}
