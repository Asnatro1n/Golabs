package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var users []User
var nextID = 1

// Обработчик для получения пользователей с пагинацией и фильтрацией
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	nameFilter := r.URL.Query().Get("name")
	ageFilterStr := r.URL.Query().Get("age")

	// Устанавливаем значения по умолчанию
	page := 1
	limit := 10

	// Парсим параметры пагинации
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			http.Error(w, "Invalid page number", http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			http.Error(w, "Invalid limit number", http.StatusBadRequest)
			return
		}
	}

	// Фильтруем пользователей
	var filteredUsers []User
	for _, user := range users {
		if (nameFilter == "" || strings.Contains(user.Name, nameFilter)) &&
			(ageFilterStr == "" || strconv.Itoa(user.Age) == ageFilterStr) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	// Пагинация
	start := (page - 1) * limit
	end := start + limit
	if start > len(filteredUsers) {
		start = len(filteredUsers)
	}
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	// Возвращаем отфильтрованный и пагинированный список пользователей
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredUsers[start:end])
}

func main() { /*
	   curl -X GET "http://localhost:8080/users"
	   curl -X GET "http://localhost:8080/users?page=2&limit=2"
	   curl -X GET "http://localhost:8080/users?name=John"
	   curl -X GET "http://localhost:8080/users?age=30"
	   curl -X GET "http://localhost:8080/users?page=1&limit=2&name=Jane"
	*/
	// Инициализация пользователей для тестирования
	users = append(users, User{ID: nextID, Name: "John Doe", Age: 30})
	nextID++
	users = append(users, User{ID: nextID, Name: "Jane Smith", Age: 25})
	nextID++
	users = append(users, User{ID: nextID, Name: "Alice Johnson", Age: 28})
	nextID++

	http.HandleFunc("/users", getUsers)
	http.ListenAndServe(":8080", nil)
}
