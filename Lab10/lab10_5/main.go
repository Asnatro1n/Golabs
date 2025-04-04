package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
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

	w.Write([]byte(tokenString))
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
			fmt.Print("Correct")
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Проверка ролей
		for _, role := range allowedRoles {
			if claims.Role == role {
				return
			}
		}

		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func main() { /*
		curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{\"username\":\"admin\", \"password\":\"adminpass\"}"
		curl -X POST http://localhost:8080/login -H "Content-Type: application/json" -d "{\"username\":\"user\", \"password\":\"userpass\"}"
		curl -X GET http://localhost:8080/admin -H "Authorization: YOUR_ADMIN_TOKEN"
		curl -X GET http://localhost:8080/user -H "Authorization: YOUR_USER_TOKEN"
	*/
	r := mux.NewRouter()

	r.HandleFunc("/login", GenerateToken).Methods("POST")
	r.HandleFunc("/admin", AuthMiddleware("admin")).Methods("GET")
	r.HandleFunc("/user", AuthMiddleware("user", "admin")).Methods("GET")

	http.ListenAndServe(":8080", r)
}
