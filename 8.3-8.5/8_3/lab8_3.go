package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

// User представляет структуру пользователя
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name" validate:"required,min=3,max=100"`
	Age  int    `json:"age" validate:"required,min=0,max=120"`
}

// Создаем соединение с базой данных
var db *sql.DB
var validate *validator.Validate

func main() { //curl -X POST -H "Content-Type: application/json" -d "{\"name\": \"Alice\", \"age\": 35}" http://localhost:8080/users
	//curl -X DELETE http://localhost:8080/users/1
	//curl -X POST -H "Content-Type: application/json" -d "{\"name\": \"Alice\", \"age\": -5}" http://localhost:8080/users
	//INSERT INTO users (name, age) VALUES ('Alice', -5);
	var err error
	// Подключаемся к базе данных
	dsn := "root:Lfq100he,ktq@tcp(127.0.0.1:3306)/userdb"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	defer db.Close()

	validate = validator.New()

	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler) // Обработка маршрутов с ID

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}

// errorHandler централизованная обработка ошибок
func errorHandler(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
}

// usersHandler обрабатывает запросы к /users
func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		addUser(w, r)
	default:
		errorHandler(w, fmt.Errorf("Метод не поддерживается"), http.StatusMethodNotAllowed)
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
		errorHandler(w, fmt.Errorf("Метод не поддерживается"), http.StatusMethodNotAllowed)
	}
}

// getUsers возвращает список всех пользователей
func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		errorHandler(w, fmt.Errorf("Ошибка при получении пользователей: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	usersList := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			errorHandler(w, fmt.Errorf("Ошибка при сканировании пользователя: %v", err), http.StatusInternalServerError)
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
		errorHandler(w, fmt.Errorf("Неверный ID: %v", err), http.StatusBadRequest)
		return
	}

	var user User
	err = db.QueryRow("SELECT id, name, age FROM users WHERE id=?", id).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			errorHandler(w, fmt.Errorf("Пользователь не найден"), http.StatusNotFound)
		} else {
			errorHandler(w, fmt.Errorf("Ошибка при получении пользователя: %v", err), http.StatusInternalServerError)
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
		errorHandler(w, fmt.Errorf("Неверный формат данных: %v", err), http.StatusBadRequest)
		return
	}

	// Валидация данных
	if err := validate.Struct(user); err != nil {
		errorHandler(w, fmt.Errorf("Ошибка валидации: %v", err), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		errorHandler(w, fmt.Errorf("Ошибка при добавлении пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		errorHandler(w, fmt.Errorf("Ошибка при получении ID пользователя: %v", err), http.StatusInternalServerError)
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
		errorHandler(w, fmt.Errorf("Неверный ID: %v", err), http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errorHandler(w, fmt.Errorf("Неверный формат данных: %v", err), http.StatusBadRequest)
		return
	}

	// Валидация данных
	if err := validate.Struct(user); err != nil {
		errorHandler(w, fmt.Errorf("Ошибка валидации: %v", err), http.StatusBadRequest)
		return
	}

	user.ID = id
	_, err = db.Exec("UPDATE users SET name=?, age=? WHERE id=?", user.Name, user.Age, user.ID)
	if err != nil {
		errorHandler(w, fmt.Errorf("Ошибка при обновлении пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// deleteUser удаляет пользователя
func deleteUser(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		errorHandler(w, fmt.Errorf("Неверный ID: %v", err), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM users WHERE id=?", id)
	if err != nil {
		errorHandler(w, fmt.Errorf("Ошибка при удалении пользователя: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
