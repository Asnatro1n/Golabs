package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() { // telnet localhost 8080
	// Указываем адрес и порт, на котором будет слушать сервер
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на localhost:8080")

	for {
		// Ожидаем подключения клиента
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при подключении:", err)
			continue
		}

		go handleConnection(conn) // Обрабатываем подключение в отдельной горутине
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Читаем сообщение от клиента
	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении сообщения:", err)
		return
	}

	fmt.Print("Получено сообщение:", message)

	// Отправляем ответ клиенту
	response := "Сообщение получено\n"
	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Ошибка при отправке ответа:", err)
		return
	}
}
