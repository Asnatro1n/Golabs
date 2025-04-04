package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	defer conn.Close() // Закрываем соединение после завершения работы с ним

	// Создаем сканер для чтения сообщений от клиента
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()                   // Читаем сообщение
		fmt.Println("Получено сообщение:", message) // Выводим сообщение на экран

		// Отправляем ответ клиенту
		_, err := conn.Write([]byte("Сообщение получено\n"))
		if err != nil {
			fmt.Println("Ошибка при отправке ответа:", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при чтении сообщения:", err)
	}
}

func main() { //telnet localhost 8080
	// Указываем порт для прослушивания
	port := ":8080"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
	defer listener.Close() // Закрываем слушатель при завершении работы

	fmt.Println("Сервер запущен на порту", port)

	// Бесконечный цикл для обработки входящих соединений
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}

		fmt.Println("Новое соединение:", conn.RemoteAddr())
		go handleConnection(conn) // Обрабатываем соединение в отдельной горутине
	}
}
