package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

// User представляет структуру пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Временное хранилище пользователей
var (
	users  = make(map[int]User)
	nextID = 1
	mu     sync.Mutex
)

func main() { /*
	 curl -X POST -H "Content-Type: application/json" -d "{\"name\": \"Alice\", \"age\": 30}" http://localhost:8080/users
	 curl http://localhost:8080/users
	 curl http://localhost:8080/users/1
	 curl -X PUT -H "Content-Type: application/json" -d "{\"name\": \"Alice Updated\", \"age\": 31}" http://localhost:8080/users/1
	 curl -X DELETE http://localhost:8080/users/1
	*/
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler) // Обработка маршрутов с ID

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}

// usersHandler обрабатывает запросы к /users
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		addUser(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// userHandler обрабатывает запросы к /users/{id}
func userHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]

	switch r.Method {
	case http.MethodGet:
		getUser(w, r, idStr)
	case http.MethodPut:
		updateUser(w, r, idStr)
	case http.MethodDelete:
		deleteUser(w, r, idStr)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// getUsers возвращает список всех пользователей
func getUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	usersList := make([]User, 0, len(users))
	for _, user := range users {
		usersList = append(usersList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersList)
}

// getUser возвращает информацию о конкретном пользователе
func getUser(w http.ResponseWriter, r *http.Request, idStr string) {
	mu.Lock()
	defer mu.Unlock()

	id, err := strconv.Atoi(idStr)
	if err != nil || users[id].ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users[id])
}

// addUser добавляет нового пользователя
func addUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	user.ID = nextID
	nextID++
	users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// updateUser обновляет информацию о пользователе
func updateUser(w http.ResponseWriter, r *http.Request, idStr string) {
	mu.Lock()
	defer mu.Unlock()

	id, err := strconv.Atoi(idStr)
	if err != nil || users[id].ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	user.ID = id
	users[id] = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// deleteUser удаляет пользователя
func deleteUser(w http.ResponseWriter, r *http.Request, idStr string) {
	mu.Lock()
	defer mu.Unlock()

	id, err := strconv.Atoi(idStr)
	if err != nil || users[id].ID == 0 {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	delete(users, id)
	w.WriteHeader(http.StatusNoContent)
}
