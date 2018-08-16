package main

import (
	"crypto/rand"
	"fmt"

	"github.com/bdbene/vault/cipher"
	"github.com/bdbene/vault/config"
	"github.com/bdbene/vault/fileio"
)

func main() {
	var conf config.Config
	config.GetConfigs(&conf)

	// Create key.
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}

	text := []byte("Hello world!")

	// Set configurations.
	io := fileio.NewFileio(conf.Storage.Location)

	// Encrypt.
	{
		ciphertext, nonce := cipher.Encrypt(key, text)
		io.WriteToFile(ciphertext, nonce)
	}

	ciphertext, nonce := io.ReadFromFile()

	deciphered := cipher.Decrypt(key, ciphertext, nonce)
	fmt.Printf("%s\n", deciphered)
}
