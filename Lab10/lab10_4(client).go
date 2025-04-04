package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var baseURL = "https://localhost:8080" // Измените на https
var sessionToken string

// Функция для авторизации пользователя
func authenticate(client *http.Client, username, password string) error {
	credentials := map[string]string{"username": username, "password": password}
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации данных: %v", err)
	}

	resp, err := client.Post(baseURL+"/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка при авторизации: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("ошибка авторизации: %s", body)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	sessionToken = result["token"] // Сохраняем токен
	fmt.Println("Успешная авторизация. Токен:", sessionToken)
	return nil
}

// Функция для добавления пользователя
func addUser(client *http.Client, username, password string) error {
	user := map[string]string{"username": username, "password": password}
	jsonData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("ошибка при сериализации данных: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, baseURL+"/adduser", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+sessionToken) // Добавляем токен в заголовок
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении пользователя: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("ошибка при добавлении пользователя: %s", body)
	}

	fmt.Println("Пользователь успешно добавлен.")
	return nil
}

// Функция для получения всех пользователей
func getAllUsers(client *http.Client) error {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/users", nil)
	if err != nil {
		return fmt.Errorf("ошибка при создании запроса: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+sessionToken) // Добавляем токен в заголовок

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при получении пользователей: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("ошибка при получении пользователей: %s", body)
	}

	var users []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	fmt.Println("Список пользователей:")
	for _, user := range users {
		fmt.Printf("ID: %v, Username: %v\n", user["id"], user["username"])
	}
	return nil
}

func main() {
	// Настройка TLS
	cert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		fmt.Printf("Ошибка загрузки сертификата клиента: %v\n", err)
		os.Exit(1)
	}

	// Загрузка CA сертификата
	caCert, err := ioutil.ReadFile("client.crt") // Замените на путь к вашему CA сертификату
	if err != nil {
		fmt.Printf("Ошибка загрузки CA сертификата: %v\n", err)
		os.Exit(1)
	}

	// Создание пула доверенных корневых сертификатов
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Создание кастомного HTTP клиента с поддержкой TLS
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool, // Устанавливаем CA сертификат
		},
	}

	client := &http.Client{Transport: tr}

	// Авторизация
	if err := authenticate(client, "user1", "password"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Добавление пользователя
	if err := addUser(client, "newuser", "newpassword"); err != nil {
		fmt.Println(err)
	}

	// Получение всех пользователей
	if err := getAllUsers(client); err != nil {
		fmt.Println(err)
	}
}
