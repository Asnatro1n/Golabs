package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
)

func main() {
	var input string
	var key string

	fmt.Println("Введите строку для шифрования:")
	fmt.Scanln(&input)

	fmt.Println("Введите секретный ключ (16, 24 или 32 байта):")
	fmt.Scanln(&key)

	// Проверка длины ключа
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		log.Fatal("Ключ должен быть длиной 16, 24 или 32 байта.")
	}

	encrypted, err := encrypt(input, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Зашифрованная строка:", encrypted)

	decrypted, err := decrypt(encrypted, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Расшифрованная строка:", decrypted)
}

// encrypt шифрует текст с использованием AES
func encrypt(plainText, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Создаем вектор инициализации
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Паддинг
	paddedText := pad([]byte(plainText), aes.BlockSize)

	// Шифруем текст
	cipherText := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, paddedText)

	// Возвращаем зашифрованный текст в виде base64
	return base64.StdEncoding.EncodeToString(append(iv, cipherText...)), nil
}

// decrypt расшифровывает текст с использованием AES
func decrypt(encryptedText, key string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Извлекаем вектор инициализации
	iv := data[:aes.BlockSize]
	cipherText := data[aes.BlockSize:]

	// Расшифровываем текст
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	// Удаляем паддинг
	unpaddedText, err := unpad(cipherText)
	if err != nil {
		return "", err
	}

	return string(unpaddedText), nil
}

// pad добавляет паддинг к тексту
func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// unpad удаляет паддинг из текста
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("unpad: empty input")
	}
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("unpad: invalid padding")
	}
	return src[:(length - unpadding)], nil
}
