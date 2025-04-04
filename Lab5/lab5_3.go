package main

import (
	"fmt"
	"os"
)

type Circle struct {
	radius float32
}

func main() {
	var circle Circle
	fmt.Print("Введите радиус круга: ")
	fmt.Fscan(os.Stdin, &circle.radius)
	circle.SquareCircle()
}

func (a Circle) SquareCircle() {
	const pi float32 = 3.1415
	S := a.radius * a.radius * pi
	fmt.Println("Площадь данного круга: ", S)
}
