package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var wg sync.WaitGroup

func handleConnection(ctx context.Context, conn net.Conn) {
	defer wg.Done()
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for {
		select {
		case <-ctx.Done():
			return // Завершаем обработку, если контекст отменен
		default:
			if scanner.Scan() {
				message := scanner.Text()
				fmt.Println("Получено сообщение:", message)
				_, err := conn.Write([]byte("Сообщение получено\n"))
				if err != nil {
					fmt.Println("Ошибка при отправке ответа:", err)
					return
				}
			} else {
				return // Завершаем, если произошла ошибка чтения
			}
		}
	}
}

func main() {
	port := ":8080"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Сервер запущен на порту", port)

	// Создаем контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Обработка сигналов для завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\nПолучен сигнал завершения. Ожидание завершения соединений...")
		cancel() // Отменяем контекст
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}

		wg.Add(1)
		fmt.Println("Новое соединение:", conn.RemoteAddr())
		go handleConnection(ctx, conn) // Обрабатываем соединение в отдельной горутине
	}

	// Ждем завершения всех горутин
	wg.Wait()
	fmt.Println("Все соединения завершены. Сервер остановлен.")
}
