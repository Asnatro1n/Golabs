package main

import (
	"fmt"
	"os"
)

type Person struct {
	name []string
	age  []int
}

func main() {
	var value int
	var person Person
	var number int
	fmt.Print("Введите количество записей в списке Person: ")
	fmt.Fscan(os.Stdin, &value)
	person.age = make([]int, value)
	person.name = make([]string, value)
	for i := 0; i < value; i++ {
		var name string
		var age int
		fmt.Print("Введите имя человека: ")
		fmt.Fscan(os.Stdin, &name)
		fmt.Print("Введите его возраст: ")
		fmt.Fscan(os.Stdin, &age)
		person.age[i] = age
		person.name[i] = name
	}
	fmt.Print("Введите номер человека в списке: ")
	fmt.Fscan(os.Stdin, &number)
	person.print(number)
}

func (a Person) print(b int) {
	fmt.Println("Имя: ", a.name[b-1])
	fmt.Println("Возраст: ", a.age[b-1])
}
