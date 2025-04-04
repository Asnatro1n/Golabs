package main

import (
	"fmt"
	"os"
)

type Rectangle struct {
	length float64
	width  float64
}

func main() {
	var dlina float64
	var shirina float64
	fmt.Print("Введите длину прямоугольника: ")
	fmt.Fscan(os.Stdin, &dlina)
	fmt.Print("Введите ширину прямоугольника: ")
	fmt.Fscan(os.Stdin, &shirina)
	ploshad := Rectangle{dlina, shirina}
	Square(ploshad.length, ploshad.width)
}

func Square(a float64, b float64) {
	ploshad := a * b
	fmt.Println("Площадь прямоугольника: ", ploshad)
}
