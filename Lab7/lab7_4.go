package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Структура для хранения данных, полученных через POST
type Data struct {
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Привет, мир!"))
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var data Data
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Ошибка разбора JSON", http.StatusBadRequest)
			return
		}

		// Выводим полученные данные в консоль
		fmt.Printf("Полученные данные: %s\n", data.Message)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Данные получены"))
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/data", dataHandler)

	fmt.Println("Сервер запущен на порту 8080") // проверка через запрос curl -X POST http://localhost:8080/data -H "Content-Type: application/json" -d "{\"message\": \"Привет, сервер!\"}"
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
