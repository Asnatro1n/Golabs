package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Структура для хранения данных, полученных через POST запрос
type Data struct {
	Message string `json:"message"`
}

// Middleware для логирования входящих запросов
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r) // Передаем управление следующему обработчику
		duration := time.Since(start)
		fmt.Printf("Метод: %s, URL: %s, Время выполнения: %s\n", r.Method, r.URL.Path, duration)
	})
}

// Обработчик для GET запроса
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Привет! Это ваше приветственное сообщение.")
}

// Обработчик для POST запроса
func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var data Data
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&data); err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Полученные данные: %s\n", data.Message)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Данные успешно получены.")
}

func main() { /*
		curl -X POST http://localhost:8080/data -H "Content-Type: application/json" -d "{\"message\": \"Привет, сервер!\"}"
		curl -X POST http://localhost:8080/hello
	*/
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/data", dataHandler)

	port := ":8080"
	fmt.Println("Сервер запущен на порту", port)

	// Оборачиваем маршрутизатор в middleware
	if err := http.ListenAndServe(port, loggingMiddleware(mux)); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
