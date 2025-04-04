package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users = make(map[int]User)
	mu    sync.Mutex
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Здесь можно добавить логику для получения всех пользователей
		w.WriteHeader(http.StatusNotImplemented) // Заглушка
	case http.MethodPost:
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		mu.Lock()
		users[user.ID] = user
		mu.Unlock()
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID пользователя из URL
	idStr := r.URL.Path[len("/users/"):]
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, exists := users[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID пользователя из URL
	idStr := r.URL.Path[len("/users/"):]
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	_, exists := users[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	users[id] = user
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID пользователя из URL
	idStr := r.URL.Path[len("/users/"):]
	var id int
	_, err := fmt.Sscanf(idStr, "%d", &id)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	_, exists := users[id]
	if !exists {
		mu.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	delete(users, id)
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/users", usersHandler)
	http.ListenAndServe(":8080", nil)
}
