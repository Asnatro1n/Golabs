package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	users  = make(map[int]User)
	nextID = 1
	mu     sync.Mutex
)

func main() {
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	http.ListenAndServe(":8080", nil)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w)
	case http.MethodPost:
		createUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

	switch r.Method {
	case http.MethodGet:
		getUser(w, id)
	case http.MethodPut:
		updateUser(w, r, id)
	case http.MethodDelete:
		deleteUser(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter) {
	mu.Lock()
	defer mu.Unlock()

	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}
	json.NewEncoder(w).Encode(userList)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.ID = nextID
	nextID++
	users[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, id string) {
	mu.Lock()
	defer mu.Unlock()

	userID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user, exists := users[userID]; exists {
		json.NewEncoder(w).Encode(user)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func updateUser(w http.ResponseWriter, r *http.Request, id string) {
	mu.Lock()
	defer mu.Unlock()

	userID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, exists := users[userID]; exists {
		user.ID = userID
		users[userID] = user
		json.NewEncoder(w).Encode(user)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteUser(w http.ResponseWriter, id string) {
	mu.Lock()
	defer mu.Unlock()

	userID, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, exists := users[userID]; exists {
		delete(users, userID)
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}
