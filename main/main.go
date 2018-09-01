package main

import (
	"fmt"
	"os"

	"github.com/bdbene/vault/cipher"
	"github.com/bdbene/vault/config"
	"github.com/bdbene/vault/storage"
)

func main() {
	// Get input.
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Invalid arguments")
		os.Exit(1)
	}
	password := args[0]

	// Read configs.
	var conf config.Config
	config.GetConfigs(&conf)

	text := []byte("Hello world!")

	// Create DataStore based on configurations.
	dataStore, err := storage.CreateDataStore(&conf.Storage)
	if err != nil {
		panic(err)
	}

	key := cipher.CreateKey(password)

	// Encrypt.
	{
		ciphertext, nonce := cipher.Encrypt(key, text)
		dataStore.Write(ciphertext, nonce)
	}

	ciphertext, nonce := dataStore.Read()

	deciphered := cipher.Decrypt(key, ciphertext, nonce)
	fmt.Printf("%s\n", deciphered)
}
