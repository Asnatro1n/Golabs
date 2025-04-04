package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// Функция для генерации ключей и их сохранения
func generateKeys() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Ошибка генерации ключа: %v", err)
	}

	// Сохранение закрытого ключа в файл
	privateKeyFile, err := os.Create("private_key.pem")
	if err != nil {
		log.Fatalf("Ошибка создания файла закрытого ключа: %v", err)
	}
	defer privateKeyFile.Close()

	err = pem.Encode(privateKeyFile, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		log.Fatalf("Ошибка записи закрытого ключа в файл: %v", err)
	}

	// Сохранение открытого ключа в файл
	publicKeyFile, err := os.Create("public_key.pem")
	if err != nil {
		log.Fatalf("Ошибка создания файла открытого ключа: %v", err)
	}
	defer publicKeyFile.Close()

	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	err = pem.Encode(publicKeyFile, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err != nil {
		log.Fatalf("Ошибка записи открытого ключа в файл: %v", err)
	}

	fmt.Println("Ключи успешно сгенерированы и сохранены в файлы.")
}

// Функция для подписания сообщения
func signMessage(message []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// Функция для проверки подписи
func verifySignature(message []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(message)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
}

func main() {
	// Генерация ключей (в реальной ситуации ключи могут быть заранее сгенерированы и сохранены)
	generateKeys()

	// Загрузка закрытого и открытого ключей из файлов
	privateKeyData, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Fatalf("Ошибка чтения закрытого ключа: %v", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyData)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatalf("Ошибка парсинга закрытого ключа: %v", err)
	}

	publicKeyData, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Fatalf("Ошибка чтения открытого ключа: %v", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyData)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		log.Fatalf("Ошибка парсинга открытого ключа: %v", err)
	}

	// Отправитель подписывает сообщение
	message := []byte("Это тестовое сообщение для подписи.")
	signature, err := signMessage(message, privateKey)
	if err != nil {
		log.Fatalf("Ошибка подписи сообщения: %v", err)
	}

	fmt.Printf("Подпись: %x\n", signature)

	// Получатель проверяет подпись
	err = verifySignature(message, signature, publicKey)
	if err != nil {
		log.Fatalf("Ошибка проверки подписи: %v", err)
	}

	fmt.Println("Подпись успешно проверена.")
}
