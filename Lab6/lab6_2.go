package main

import (
	"fmt"
)

// Функция для генерации первых n чисел Фибоначчи и отправки их в канал
func fibonacci(n int, ch chan<- int) {
	a, b := 0, 1
	for i := 0; i < n; i++ {
		ch <- a       // Отправляем число в канал
		a, b = b, a+b // Генерируем следующее число Фибоначчи
	}
	close(ch) // Закрываем канал после отправки всех чисел
}

// Функция для чтения данных из канала и их вывода на экран
func printFibonacci(ch <-chan int) {
	for num := range ch { // Читаем из канала до его закрытия
		fmt.Println(num)
	}
}

func main() {
	ch := make(chan int) // Создаем новый канал

	go fibonacci(10, ch) // Запускаем горутину для генерации чисел Фибоначчи
	printFibonacci(ch)   // Запускаем чтение из канала
}
