package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// Функция для реверсирования строки
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Структура для задания
type Task struct {
	Line string
}

// Функция воркера
func worker(tasks <-chan Task, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		reversed := reverse(task.Line)
		results <- reversed // Отправляем результат обратно
	}
}

// Функция для чтения строк из файла и распределения задач
func processFile(filename string, numWorkers int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	tasks := make(chan Task)
	results := make(chan string)

	var wg sync.WaitGroup

	// Запускаем воркеров
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(tasks, results, &wg)
	}

	// Читаем строки из файла и отправляем задачи
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			tasks <- Task{Line: scanner.Text()}
		}
		close(tasks) // Закрываем канал задач после завершения чтения
	}()

	// Ожидаем завершения воркеров
	go func() {
		wg.Wait()
		close(results) // Закрываем канал результатов после завершения всех воркеров
	}()

	// Выводим результаты в порядке их поступления
	for result := range results {
		fmt.Println(result)
	}
}

func main() {
	var filename string
	var numWorkers int

	fmt.Print("Введите имя файла: ")
	fmt.Scan(&filename)
	fmt.Print("Введите количество воркеров: ")
	fmt.Scan(&numWorkers)

	processFile(filename, numWorkers)
}
