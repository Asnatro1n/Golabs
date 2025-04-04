package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	//"sync"
	"testing"
)

/*var (
    users = make(map[int]User)
    mu    sync.Mutex
)*/

func TestUsersHandler(t *testing.T) {
	// Тестирование POST запроса
	user := User{ID: 1, Name: "John Doe"}
	jsonData, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(usersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Проверка, что пользователь добавлен
	mu.Lock()
	if u, exists := users[1]; !exists || u != user {
		t.Errorf("user not added: got %v want %v", u, user)
	}
	mu.Unlock()
}

func TestGetUser(t *testing.T) {
	user := User{ID: 1, Name: "John Doe"}
	mu.Lock()
	users[user.ID] = user
	mu.Unlock()

	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedUser User
	json.NewDecoder(rr.Body).Decode(&returnedUser)

	if returnedUser != user {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedUser, user)
	}
}

func TestUpdateUser(t *testing.T) {
	user := User{ID: 1, Name: "John Doe"}
	mu.Lock()
	users[user.ID] = user
	mu.Unlock()

	updatedUser := User{ID: 1, Name: "Jane Doe"}
	jsonData, _ := json.Marshal(updatedUser)

	req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	mu.Lock()
	if users[1].Name != updatedUser.Name {
		t.Errorf("user not updated: got %v want %v", users[1].Name, updatedUser.Name)
	}
	mu.Unlock()
}

func TestDeleteUser(t *testing.T) {
	user := User{ID: 1, Name: "John Doe"}
	mu.Lock()
	users[user.ID] = user
	mu.Unlock()

	req, err := http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteUser)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	mu.Lock()
	if _, exists := users[1]; exists {
		t.Errorf("user not deleted")
	}
	mu.Unlock()
}
