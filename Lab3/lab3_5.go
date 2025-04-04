package main

import "fmt"

func main() {
	var delete = 4
	massiv := [5]string{"Sten", "William", "Willson", "Robert", "Porshe"}
	srez := massiv[:]
	srez = append(srez, "Anton")
	fmt.Println(srez)
	srez = append(srez[:delete], srez[delete+1:]...)
	fmt.Println(srez)
}
