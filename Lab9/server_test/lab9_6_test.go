package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestikLogin(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(login))
	defer server.Close()

	// Успешный вход
	user := User{Username: "user1", Password: "password"}
	jsonData, _ := json.Marshal(user)

	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидали статус 200, получили %d", resp.StatusCode)
	}

	// Неверные учетные данные
	user = User{Username: "user1", Password: "wrongpassword"}
	jsonData, _ = json.Marshal(user)

	resp, err = http.Post(server.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Ожидали статус 401, получили %d", resp.StatusCode)
	}

	// Пустые поля
	user = User{Username: "", Password: ""}
	jsonData, _ = json.Marshal(user)

	resp, err = http.Post(server.URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка при отправке запроса: %v", err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Ожидали статус 400, получили %d", resp.StatusCode)
	}
}
