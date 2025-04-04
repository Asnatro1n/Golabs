package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Функция для генерации случайных чисел и отправки их в канал
func generateRandomNumbers(ch chan<- int) {
	for {
		num := rand.Intn(100)   // Генерация случайного числа от 0 до 99
		ch <- num               // Отправка числа в канал
		time.Sleep(time.Second) // Задержка для имитации времени генерации
	}
}

// Функция для определения четности/нечетности числа и отправки сообщения в канал
func checkEvenOdd(ch <-chan int, resultCh chan<- string) {
	for num := range ch {
		if num%2 == 0 {
			resultCh <- fmt.Sprintf("%d - четное", num)
		} else {
			resultCh <- fmt.Sprintf("%d - нечетное", num)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел

	numCh := make(chan int)       // Канал для случайных чисел
	resultCh := make(chan string) // Канал для результатов четности/нечетности

	go generateRandomNumbers(numCh)  // Запуск горутины для генерации чисел
	go checkEvenOdd(numCh, resultCh) // Запуск горутины для проверки четности/нечетности

	// Используем select для приема данных из канала результатов
	for i := 0; i < 10; i++ { // Ограничиваем количество выводов
		select {
		case result := <-resultCh:
			fmt.Println(result) // Вывод результата
		}
	}

	close(numCh)    // Закрываем канал случайных чисел
	close(resultCh) // Закрываем канал результатов
}
