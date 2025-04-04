package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const baseURL = "http://localhost:8080/users"

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	for {
		fmt.Println("Выберите операцию:")
		fmt.Println("1. Получить всех пользователей")
		fmt.Println("2. Создать пользователя")
		fmt.Println("3. Получить пользователя по ID")
		fmt.Println("4. Обновить пользователя")
		fmt.Println("5. Удалить пользователя")
		fmt.Println("6. Выход")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			getUsers()
		case 2:
			createUser()
		case 3:
			getUser()
		case 4:
			updateUser()
		case 5:
			deleteUser()
		case 6:
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Неверный выбор. Пожалуйста, попробуйте снова.")
		}
	}
}

func handleResponse(resp *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("ошибка: %s", string(body))
	}

	return body, nil
}

func formatUser(user User) {
	fmt.Printf("ID: %d, Имя: %s\n", user.ID, user.Name)
}

func formatUsers(users []User) {
	fmt.Println("Список пользователей:")
	for _, user := range users {
		formatUser(user)
	}
}

func getUsers() {
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	body, err := handleResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		fmt.Println("Ошибка декодирования:", err)
		return
	}

	formatUsers(users)
}

func createUser() {
	var name string
	fmt.Print("Введите имя пользователя: ")
	fmt.Scan(&name)

	user := User{Name: name}
	jsonData, _ := json.Marshal(user)

	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	body, err := handleResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	var createdUser User
	if err := json.Unmarshal(body, &createdUser); err != nil {
		fmt.Println("Ошибка декодирования:", err)
		return
	}

	fmt.Println("Пользователь успешно создан:")
	formatUser(createdUser)
}

func getUser() {
	var id int
	fmt.Print("Введите ID пользователя: ")
	fmt.Scan(&id)

	resp, err := http.Get(fmt.Sprintf("%s/%d", baseURL, id))
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	body, err := handleResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		fmt.Println("Ошибка декодирования:", err)
		return
	}

	fmt.Println("Пользователь найден:")
	formatUser(user)
}

func updateUser() {
	var id int
	fmt.Print("Введите ID пользователя для обновления: ")
	fmt.Scan(&id)

	var name string
	fmt.Print("Введите новое имя пользователя: ")
	fmt.Scan(&name)

	user := User{Name: name}
	jsonData, _ := json.Marshal(user)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%d", baseURL, id), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	body, err := handleResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	var updatedUser User
	if err := json.Unmarshal(body, &updatedUser); err != nil {
		fmt.Println("Ошибка декодирования:", err)
		return
	}

	fmt.Println("Пользователь успешно обновлен:")
	formatUser(updatedUser)
}

func deleteUser() {
	var id int
	fmt.Print("Введите ID пользователя для удаления: ")
	fmt.Scan(&id)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", baseURL, id), nil)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	body, err := handleResponse(resp)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Используем body для вывода сообщения
	fmt.Println("Ответ сервера:", string(body))
	fmt.Println("Пользователь успешно удален.")
}
