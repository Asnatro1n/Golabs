package main

import (
	"fmt"
	"sync"
)

// Общая переменная-счетчик
var counter int

// Мьютекс для синхронизации доступа к счетчику
var mu sync.Mutex

// Функция для увеличения счетчика
func increment(wg *sync.WaitGroup) {
	defer wg.Done() // Уменьшаем счетчик WaitGroup по завершении функции
	for i := 0; i < 1000; i++ {
		mu.Lock()   // Блокируем мьютекс перед доступом к счетчику
		counter++   // Увеличиваем счетчик
		mu.Unlock() // Освобождаем мьютекс
	}
}

func main() {
	var wg sync.WaitGroup

	// Запускаем 5 горутин для увеличения счетчика
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go increment(&wg)
	}

	wg.Wait() // Ожидаем завершения всех горутин
	fmt.Println("Счетчик с мьютексом:", counter)

	// Сброс счетчика
	counter = 0

	// Запускаем 5 горутин без мьютекса
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter++ // Увеличиваем счетчик без блокировки
			}
		}()
	}

	wg.Wait() // Ожидаем завершения всех горутин
	fmt.Println("Счетчик без мьютекса:", counter)
}
