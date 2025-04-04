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

func getUsers() {
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Статус ответа: %s\n", resp.Status)
	if resp.StatusCode == http.StatusOK {
		var users []User
		json.Unmarshal(body, &users)
		fmt.Printf("Пользователи: %+v\n", users)
	} else {
		fmt.Println("Ошибка:", string(body))
	}
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

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Статус ответа: %s\n", resp.Status)
	if resp.StatusCode == http.StatusCreated {
		var createdUser User
		json.Unmarshal(body, &createdUser)
		fmt.Printf("Пользователь успешно создан: %+v\n", createdUser)
	} else {
		fmt.Println("Ошибка:", string(body))
	}
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

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Статус ответа: %s\n", resp.Status)
	if resp.StatusCode == http.StatusOK {
		var user User
		json.Unmarshal(body, &user)
		fmt.Printf("Пользователь: %+v\n", user)
	} else {
		fmt.Println("Ошибка:", string(body))
	}
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

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Статус ответа: %s\n", resp.Status)
	if resp.StatusCode == http.StatusOK {
		var updatedUser User
		json.Unmarshal(body, &updatedUser)
		fmt.Printf("Пользователь успешно обновлен: %+v\n", updatedUser)
	} else {
		fmt.Println("Ошибка:", string(body))
	}
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

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Статус ответа: %s\n", resp.Status)
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Пользователь успешно удален.")
	} else {
		fmt.Println("Ошибка:", string(body))
	}
}
