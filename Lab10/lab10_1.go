package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func hashString(input string, algorithm string) string {
	var hash []byte

	switch algorithm {
	case "md5":
		h := md5.New()
		h.Write([]byte(input))
		hash = h.Sum(nil)
	case "sha256":
		h := sha256.New()
		h.Write([]byte(input))
		hash = h.Sum(nil)
	case "sha512":
		h := sha512.New()
		h.Write([]byte(input))
		hash = h.Sum(nil)
	default:
		fmt.Println("Некорректный алгоритм. Доступные варианты: md5, sha256, sha512")
		return ""
	}

	return hex.EncodeToString(hash)
}

func checkIntegrity(input string, expectedHash string, algorithm string) bool {
	actualHash := hashString(input, algorithm)
	return actualHash == expectedHash
}

func main() {
	var input string
	var algorithm string

	fmt.Println("Введите строку для хэширования:")
	fmt.Scanln(&input)

	fmt.Println("Выберите алгоритм хэширования (md5, sha256, sha512):")
	fmt.Scanln(&algorithm)

	hash := hashString(input, algorithm)
	if hash != "" {
		fmt.Printf("Хэш строки с использованием %s: %s\n", algorithm, hash)
	}

	// Проверка целостности данных
	var providedHash string
	fmt.Println("Введите хэш для проверки целостности:")
	fmt.Scanln(&providedHash)

	// Запрашиваем строку для проверки
	var checkInput string
	fmt.Println("Введите строку для проверки целостности:")
	fmt.Scanln(&checkInput)

	if checkIntegrity(checkInput, providedHash, algorithm) {
		fmt.Println("Хэш соответствует введенной строке.")
	} else {
		fmt.Println("Хэш не соответствует введенной строке.")
	}
}
