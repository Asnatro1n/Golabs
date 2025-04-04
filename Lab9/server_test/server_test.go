package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ID       int    `json:"id,omitempty"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func TestLogin(t *testing.T) {
	tests := []struct {
		username string
		password string
		expected int
	}{
		{"user1", "password", http.StatusOK},
		{"user1", "wrongpassword", http.StatusUnauthorized},
		{"unknownuser", "password", http.StatusUnauthorized},
	}

	for _, test := range tests {
		user := User{Username: test.username, Password: test.password}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		login(w, req)

		res := w.Result()
		if res.StatusCode != test.expected {
			t.Errorf("Expected status %d but got %d for user %s", test.expected, res.StatusCode, test.username)
		}
	}
}

func TestAddUser(t *testing.T) {
	tests := []struct {
		username string
		password string
		expected int
	}{
		{"newuser", "newpassword", http.StatusCreated},
	}

	for _, test := range tests {
		user := User{Username: test.username, Password: test.password}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest("POST", "/adduser", bytes.NewBuffer(body))

		// Simulate a login to get a valid token for adding a user
		loginUser := User{Username: "user1", Password: "password"}
		loginBody, _ := json.Marshal(loginUser)
		loginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
		loginW := httptest.NewRecorder()
		login(loginW, loginReq)

		var tokenResponse TokenResponse
		json.NewDecoder(loginW.Body).Decode(&tokenResponse)

		req.Header.Set("Authorization", tokenResponse.Token)

		w := httptest.NewRecorder()

		addUser(w, req)

		res := w.Result()
		if res.StatusCode != test.expected {
			t.Errorf("Expected status %d but got %d for adding user %s", test.expected, res.StatusCode, test.username)
		}
	}
}

func TestGetAllUsers(t *testing.T) {
	// Simulate a login to get a valid token for getting users
	loginUser := User{Username: "user1", Password: "password"}
	loginBody, _ := json.Marshal(loginUser)
	loginReq := httptest.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	loginW := httptest.NewRecorder()
	login(loginW, loginReq)

	var tokenResponse TokenResponse
	json.NewDecoder(loginW.Body).Decode(&tokenResponse)

	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", tokenResponse.Token)

	w := httptest.NewRecorder()

	getAllUsers(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d but got %d for getting all users", http.StatusOK, res.StatusCode)
	}
}
