package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret_key") // Секретный ключ для подписи токенов

// Структура для хранения данных пользователя
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// Структура для токена
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// Хранилище пользователей
var (
	validUsers = map[string]string{
		"user1": "password",
		"user2": "password",
		"user3": "password",
	}
	users  = make(map[int]User)
	nextID = 1
	mu     sync.Mutex
)

// Функция для авторизации пользователя
func login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Неверные данные"}`, http.StatusBadRequest)
		return
	}

	// Проверка учетных данных
	if validPassword, ok := validUsers[user.Username]; !ok || validPassword != user.Password {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Неверные учетные данные"}`, http.StatusUnauthorized)
		return
	}

	// Создание токена
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Ошибка при создании токена"}`, http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"token": tokenString}
	json.NewEncoder(w).Encode(response)
}

// Middleware для проверки токена
func tokenValid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Токен не предоставлен", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Функция для добавления пользователя
func addUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user.ID = nextID
	nextID++
	users[user.ID] = user
	mu.Unlock()

	// Создание токена для нового пользователя
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Ошибка при создании токена"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"user":  user,
		"token": tokenString,
	}
	json.NewEncoder(w).Encode(response)
}

// Функция для получения всех пользователей
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var userList []User
	for _, user := range users {
		userList = append(userList, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userList)
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/adduser", addUser)
	http.Handle("/users", tokenValid(http.HandlerFunc(getAllUsers)))

	fmt.Println("Сервер запущен на порту 8080")
	http.ListenAndServe(":8080", nil)
}
