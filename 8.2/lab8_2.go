package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// User представляет структуру пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Создаем соединение с базой данных
var db *sql.DB

func main() { //curl -X POST -H "Content-Type: application/json" -d "{\"name\": \"Alice\", \"age\": 30}" http://localhost:8080/users
	var err error
	// Подключаемся к базе данных
	dsn := "asnatroin:Lfq100he,ktq@tcp(127.0.0.1:3306)/userdb" // Замените username и password на ваши учетные данные
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler) // Обработка маршрутов с ID

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}

// usersHandler обрабатывает запросы к /users
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		addUser(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// userHandler обрабатывает запросы к /users/{id}
func userHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/users/"):]

	switch r.Method {
	case http.MethodGet:
		getUser(w, r, idStr)
	case http.MethodPut:
		updateUser(w, r, idStr)
	case http.MethodDelete:
		deleteUser(w, r, idStr)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// getUsers возвращает список всех пользователей
func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		http.Error(w, "Ошибка при получении пользователей", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	usersList := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			http.Error(w, "Ошибка при сканировании пользователя", http.StatusInternalServerError)
			return
		}
		usersList = append(usersList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usersList)
}

// getUser возвращает информацию о конкретном пользователе
func getUser(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow("SELECT id, name, age FROM users WHERE id=?", id).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при получении пользователя", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// addUser добавляет нового пользователя
func addUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Ошибка при получении ID пользователя", http.StatusInternalServerError)
		return
	}

	user.ID = int(id)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// updateUser обновляет информацию о пользователе
func updateUser(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	user.ID = id
	_, err = db.Exec("UPDATE users SET name=?, age=? WHERE id=?", user.Name, user.Age, user.ID)
	if err != nil {
		http.Error(w, "Ошибка при обновлении пользователя", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// deleteUser удаляет пользователя
func deleteUser(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
