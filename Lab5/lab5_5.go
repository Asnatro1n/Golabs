package main

import (
	"fmt"
	"math"
)

// Определяем интерфейс Shape
type Shape interface {
	Area() float64
}

// Структура для круга
type Circle struct {
	radius float64
}

// Метод для вычисления площади круга
func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.radius, 2)
}

// Структура для прямоугольника
type Rectangle struct {
	width, height float64
}

// Метод для вычисления площади прямоугольника
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

// Функция, которая принимает срез фигур и выводит их площади
func PrintAreas(shapes []Shape) {
	for _, shape := range shapes {
		fmt.Printf("Area: %.2f\n", shape.Area())
	}
}

func main() {
	// Создаем объекты фигур
	circle := Circle{radius: 5}
	rectangle := Rectangle{width: 4, height: 6}

	// Создаем срез фигур
	shapes := []Shape{circle, rectangle}

	// Выводим площади фигур
	PrintAreas(shapes)
}
