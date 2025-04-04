package main

import (
	"fmt"
	"os"
)

type Shape interface {
	Area()
}

type Rectangle struct {
	length float32
	width  float32
}

func (a Rectangle) Area() {
	fmt.Println("Площадь прямоугольника: ", a.length*a.width)
}

type Circle struct {
	radius float32
}

func (b Circle) Area() {
	const pi float32 = 3.1415
	fmt.Println("Площадь круга: ", b.radius*b.radius*pi)
}

func main() {
	var radius float32
	var dlina float32
	var shirina float32
	fmt.Print("Введите длину прямоугольника: ")
	fmt.Fscan(os.Stdin, &dlina)
	fmt.Print("Введите ширину прямоугольника: ")
	fmt.Fscan(os.Stdin, &shirina)
	var rectangle Shape = Rectangle{dlina, shirina}
	fmt.Print("Введите радиус круга: ")
	fmt.Fscan(os.Stdin, &radius)
	var circle Shape = Circle{radius}
	rectangle.Area()
	circle.Area()
}
