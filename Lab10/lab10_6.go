package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type User struct {
	Username string
	Password string
	Role     string // "admin" или "user"
}

var users = map[string]User{
	"admin": {Username: "admin", Password: "adminpass", Role: "admin"},
	"user":  {Username: "user", Password: "userpass", Role: "user"},
}

var jwtKey = []byte("my_secret_key")
var store = sessions.NewCookieStore([]byte("session_key"))

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	storedUser, exists := users[user.Username]
	if !exists || storedUser.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Генерация токена CSRF
	session, _ := store.Get(r, "session-name")
	session.Values["csrf_token"] = generateCSRFToken()
	session.Save(r, w)

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		Role:     storedUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-CSRF-Token", session.Values["csrf_token"].(string))
	w.Write([]byte(tokenString))
}

func generateCSRFToken() string {
	// Генерация токена CSRF (можно использовать более сложный механизм)
	return "random_csrf_token"
}

func AuthMiddleware(allowedRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println("Доступ к защищенному ресурсу")

		// Проверка ролей
		for _, role := range allowedRoles {
			if claims.Role == role {
				return
			}
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func CSRFProtection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропустить проверку CSRF для маршрута /login
		if r.Method == http.MethodPost && r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			session, _ := store.Get(r, "session-name")
			csrfToken := r.Header.Get("X-CSRF-Token")
			if csrfToken != session.Values["csrf_token"] {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func main() { /*
		curl -X GET http://localhost:8080/admin -H "Content-Type: application/json" -d "{\"username\": \"admin\", \"password\": \"adminpass\"}"
	*/
	r := mux.NewRouter()

	r.HandleFunc("/login", GenerateToken).Methods("POST")
	r.HandleFunc("/admin", AuthMiddleware("admin")).Methods("GET")
	r.HandleFunc("/user", AuthMiddleware("user", "admin")).Methods("GET")

	// Обернуть маршруты в защиту от CSRF
	r.Use(CSRFProtection)

	http.ListenAndServe(":8080", r)
}
