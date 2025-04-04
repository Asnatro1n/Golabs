package main

import (
	"fmt"
)

// Определяем структуру Book
type Book struct {
	Title  string
	Author string
	Year   int
}

// Реализуем метод String() для структуры Book
func (b Book) String() string {
	return fmt.Sprintf("'%s' by %s (%d)", b.Title, b.Author, b.Year)
}

func main() {
	// Создаем экземпляр книги
	book := Book{
		Title:  "1984",
		Author: "George Orwell",
		Year:   1949,
	}

	// Выводим информацию о книге, используя метод String()
	fmt.Println(book)
}
