package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock server for testing client functions
func mockServer() *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/login":
			var user User
			json.NewDecoder(r.Body).Decode(&user)
			if user.Username == "user1" && user.Password == "password" {
				tokenResponse := TokenResponse{Token: "mock_token"}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(tokenResponse)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		case "/adduser":
			var user User
			json.NewDecoder(r.Body).Decode(&user)
			if r.Header.Get("Authorization") == "mock_token" {
				user.ID = 1 // Mock ID assignment
				response := struct {
					User  User   `json:"user"`
					Token string `json:"token"`
				}{
					User:  user,
					Token: "mock_token",
				}
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(response)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		case "/users":
			if r.Header.Get("Authorization") == "mock_token" {
				usersList := []User{
					{ID: 1, Username: "user1"},
					{ID: 2, Username: "newuser"},
				}
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(usersList)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	return httptest.NewServer(handler)
}

func TestClientLogin(t *testing.T) {
	server := mockServer()
	defer server.Close()

	err := login("user1", "password")
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
	if token == "" {
		t.Fatal("Expected token to be set but it is empty")
	}
}

func TestClientAddUser(t *testing.T) {
	server := mockServer()
	defer server.Close()

	token = "mock_token" // Set the mock token directly

	err := addUser("newuser", "newpassword")
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
}

func TestClientGetAllUsers(t *testing.T) {
	server := mockServer()
	defer server.Close()

	token = "mock_token" // Set the mock token directly

	err := getAllUsers()
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}
}
