package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Функция для вычисления факториала
func factorial(n int, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup по завершении функции
	result := 1
	time.Sleep(time.Second) // Имитация задержки
	for i := 1; i <= n; i++ {
		result *= i
	}
	fmt.Printf("Факториал %d: %d\n", n, result)
}

// Функция для генерации случайных чисел
func randomNumbers(count int, wg *sync.WaitGroup) {
	defer wg.Done()             // Уменьшаем счетчик WaitGroup по завершении функции
	time.Sleep(time.Second * 2) // Имитация задержки
	for i := 0; i < count; i++ {
		num := rand.Intn(100) // Генерация случайного числа от 0 до 99
		fmt.Printf("Случайное число: %d\n", num)
	}
}

// Функция для вычисления суммы числового ряда
func sumSeries(n int, wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup по завершении функции
	sum := 0
	time.Sleep(time.Second * 3) // Имитация задержки
	for i := 1; i <= n; i++ {
		sum += i
	}
	fmt.Printf("Сумма числового ряда до %d: %d\n", n, sum)
}

func main() {
	var wg sync.WaitGroup // Создаем WaitGroup для ожидания завершения горутин

	// Запускаем горутины
	wg.Add(3) // Увеличиваем счетчик WaitGroup
	go factorial(5, &wg)
	go randomNumbers(5, &wg)
	go sumSeries(5, &wg)

	wg.Wait() // Ждем завершения всех горутин
	fmt.Println("Все горутины завершили выполнение.")
}
