package main

import (
	"fmt"
	"sync"
)

// Определяем структуру для запроса на выполнение операции
type Request struct {
	Operation string
	Operand1  float64
	Operand2  float64
	ResultCh  chan float64
}

// Функция для обработки запросов
func calculator(requests <-chan Request, wg *sync.WaitGroup) {
	defer wg.Done()
	for req := range requests {
		var result float64
		switch req.Operation {
		case "+":
			result = req.Operand1 + req.Operand2
		case "-":
			result = req.Operand1 - req.Operand2
		case "*":
			result = req.Operand1 * req.Operand2
		case "/":
			if req.Operand2 != 0 {
				result = req.Operand1 / req.Operand2
			} else {
				fmt.Println("Ошибка: деление на ноль")
				result = 0
			}
		default:
			fmt.Println("Ошибка: неизвестная операция")
			result = 0
		}
		req.ResultCh <- result // Отправляем результат обратно
	}
}

func main() {
	var wg sync.WaitGroup
	requests := make(chan Request)

	// Запускаем горутину калькулятора
	wg.Add(1)
	go calculator(requests, &wg)

	// Отправляем запросы на выполнение операций
	for i := 0; i < 5; i++ {
		resultCh := make(chan float64)
		requests <- Request{Operation: "+", Operand1: float64(i), Operand2: float64(i + 1), ResultCh: resultCh}
		fmt.Printf("Запрос: %d + %d = %f\n", i, i+1, <-resultCh)

		resultCh = make(chan float64)
		requests <- Request{Operation: "-", Operand1: float64(i + 1), Operand2: float64(i), ResultCh: resultCh}
		fmt.Printf("Запрос: %d - %d = %f\n", i+1, i, <-resultCh)

		resultCh = make(chan float64)
		requests <- Request{Operation: "*", Operand1: float64(i), Operand2: float64(i + 2), ResultCh: resultCh}
		fmt.Printf("Запрос: %d * %d = %f\n", i, i+2, <-resultCh)

		resultCh = make(chan float64)
		requests <- Request{Operation: "/", Operand1: float64(i + 2), Operand2: float64(i + 1), ResultCh: resultCh}
		fmt.Printf("Запрос: %d / %d = %f\n", i+2, i+1, <-resultCh)
	}

	close(requests) // Закрываем канал запросов
	wg.Wait()       // Ожидаем завершения горутины
}
