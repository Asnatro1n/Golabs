package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Функция для обработки ответа
func handleResponse(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("ошибка: %s, статус: %s", body, resp.Status)
	}
	return nil
}

// Функция для вывода пользователя
func printUser(user User) {
	fmt.Printf("ID: %d, Имя: %s\n", user.ID, user.Name)
}

// Функция для получения всех пользователей
func getAllUsers() []User {
	resp, err := http.Get("http://localhost:8080/users")
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return nil
	}
	defer resp.Body.Close()

	if err := handleResponse(resp); err != nil {
		fmt.Println(err)
		return nil
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		fmt.Println("Ошибка при декодировании ответа:", err)
		return nil
	}

	return users
}

// Функция для добавления пользователя
func addUser(name string) {
	user := User{Name: name}
	jsonData, _ := json.Marshal(user)

	resp, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	if err := handleResponse(resp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Пользователь добавлен.")
}

// Функция для удаления пользователя
func deleteUser(id int) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/users/%d", id), nil)
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	if err := handleResponse(resp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Пользователь удален.")
}

// Функция для обновления пользователя
func updateUser(id int, name string) {
	user := User{ID: id, Name: name}
	jsonData, _ := json.Marshal(user)

	req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/users/%d", id), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	if err := handleResponse(resp); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Пользователь обновлен.")
}

// Функция для отображения меню
func displayMenu() {
	fmt.Println("Выберите действие:")
	fmt.Println("1. Показать всех пользователей")
	fmt.Println("2. Добавить пользователя")
	fmt.Println("3. Удалить пользователя")
	fmt.Println("4. Обновить пользователя")
	fmt.Println("5. Выход")
}

// Основная функция
func main() {
	for {
		displayMenu()

		var choice int
		fmt.Print("Введите номер действия: ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		switch choice {
		case 1:
			users := getAllUsers()
			if users != nil {
				fmt.Println("Список пользователей:")
				for _, user := range users {
					printUser(user)
				}
			}
		case 2:
			var name string
			fmt.Print("Введите имя пользователя: ")
			fmt.Scan(&name)
			addUser(name)
		case 3:
			var id int
			fmt.Print("Введите ID пользователя для удаления: ")
			fmt.Scan(&id)
			deleteUser(id)
		case 4:
			var id int
			var name string
			fmt.Print("Введите ID пользователя для обновления: ")
			fmt.Scan(&id)
			fmt.Print("Введите новое имя пользователя: ")
			fmt.Scan(&name)
			updateUser(id, name)
		case 5:
			fmt.Println("Выход из программы.")
			os.Exit(0)
		default:
			fmt.Println("Неверный выбор. Пожалуйста, попробуйте снова.")
		}
	}
}
