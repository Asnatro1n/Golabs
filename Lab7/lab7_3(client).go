package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Указываем адрес и порт сервера
	serverAddress := "localhost:8080"

	// Подключаемся к серверу
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Подключено к серверу", serverAddress)

	// Создаем сканер для чтения ввода пользователя
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите сообщение (или 'exit' для выхода): ")
		scanner.Scan()
		message := scanner.Text()

		if message == "exit" {
			break // Выход из цикла, если пользователь ввел 'exit'
		}

		// Отправляем сообщение на сервер
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Ошибка при отправке сообщения:", err)
			return
		}

		// Читаем ответ от сервера
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Ошибка при чтении ответа:", err)
			return
		}

		fmt.Print("Ответ от сервера:", response)
	}
}
