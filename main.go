package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func main() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}

	text := []byte("Hello world")

	ciphertext, nonce := encrypt(key, text)
	fmt.Printf("%x\n", ciphertext)

	deciphered := decrypt(key, ciphertext, nonce)
	fmt.Printf("%s\n", deciphered)
}

func encrypt(key []byte, plaintext []byte) ([]byte, []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce
}

func decrypt(key []byte, ciphertext []byte, nonce []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return plaintext
}
