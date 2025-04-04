package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// Секретный ключ для подписывания JWT
	secretKey = []byte("my_secret_key")

	// Список пользователей для примера
	users = map[string]string{
		"admin": "password", // admin: password
		"user":  "password", // user: password
	}

	// Хранение ролей пользователей
	roles = map[string]string{
		"admin": "admin",
		"user":  "user",
	}
)

// Структура для представления пользователя
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Структура для представления токена
type Token struct {
	Token string `json:"token"`
}

// Проверка роли пользователя
func checkRole(username, role string) bool {
	userRole, exists := roles[username]
	return exists && userRole == role
}

// Генерация JWT для пользователя
func generateToken(username, role string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // срок действия 2 часа
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Обработчик для входа пользователя и получения токена
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if password, ok := users[user.Username]; ok && password == user.Password {
		role := roles[user.Username]
		tokenString, err := generateToken(user.Username, role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Token{Token: tokenString})
		return
	}

	http.Error(w, "невозможно войти", http.StatusUnauthorized)
}

// Middleware для проверки JWT и ролей
func authMiddleware(requiredRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "отсутствует токен", http.StatusUnauthorized)
			return
		}

		token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return secretKey, nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			username := claims["username"].(string)
			// := claims["role"].(string)

			if checkRole(username, requiredRole) {
				// Если роль подтверждена, продолжаем
				w.Write([]byte("Доступ к защищенному ресурсу для " + username))
				return
			}
		}

		http.Error(w, "доступ запрещен", http.StatusForbidden)
	}
}

// Главная функция для настройки маршрутов
func main() { //через postman
	//выполнить login через user password и admin password
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/admin", authMiddleware("admin")) // Только для администраторов
	http.HandleFunc("/user", authMiddleware("user"))   // Для обычных пользователей
	http.ListenAndServe(":8080", nil)
}
