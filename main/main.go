package main

import (
	"crypto/rand"
	"fmt"

	"github.com/bdbene/vault/cipher"
	"github.com/bdbene/vault/fileio"
)

func main() {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}

	text := []byte("Hello world")

	ciphertext, nonce := cipher.Encrypt(key, text)
	fileio.WriteToFile(ciphertext, nonce)

	deciphered := cipher.Decrypt(key, ciphertext, nonce)
	fmt.Printf("%s\n", deciphered)
}
