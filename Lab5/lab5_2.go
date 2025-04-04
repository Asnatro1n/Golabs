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
	var nomer int
	fmt.Print("Введите количество записей в списке Person: ")
	fmt.Fscan(os.Stdin, &value)
	person.age = make([]int, value)
	person.name = make([]string, value)
	for i := 0; i < value; i++ {
		var name string
		var age int
		fmt.Print("Введите имя человека:")
		fmt.Fscan(os.Stdin, &name)
		fmt.Print("Введите его возраст: ")
		fmt.Fscan(os.Stdin, &age)
		person.age[i] = age
		person.name[i] = name
	}
	fmt.Print("Введите номер человека в списке, у которого день рождения: ")
	fmt.Fscan(os.Stdin, &nomer)
	fmt.Print("До дня рождения: ")
	person.VivodPerson(nomer)
	person.birthday(nomer)
	fmt.Print("После дня рождения: ")
	person.VivodPerson(nomer)
}

func (a Person) birthday(b int) {
	if len(a.name) < b {
		fmt.Print("Указан несуществующий номер в списке Person")
	} else {
		a.age[b-1] = a.age[b-1] + 1
	}
}

func (c Person) VivodPerson(d int) {
	if len(c.name) < d {
		fmt.Print("Указан несуществующий номер в списке Person")
	} else {
		fmt.Println(c.name[d-1], " - ", c.age[d-1], "years")
	}
}
