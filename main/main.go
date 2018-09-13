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
	if len(args) != 2 {
		fmt.Println("Invalid arguments")
		os.Exit(1)
	}

	identifier := args[0]
	password := args[1]

	// Read configs.
	var conf config.Config
	err := config.GetConfigs(&conf)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		fmt.Printf("Cannot continue, shutting down.\n")
		os.Exit(1)
	}


	text := []byte("Hello world!")

	// Create DataStore based on configurations.
	dataStore, err := storage.CreateDataStore(&conf.Storage)
	if err != nil {
		panic(err)
	}

	key, _ := cipher.CreateKey(password)

	// Encrypt.
	{
		ciphertext, nonce, _ := cipher.Encrypt(key, text)
		dataStore.Write([]byte(identifier), ciphertext, nonce)
		ciphertext, nonce, _ = cipher.Encrypt(key, text)
		dataStore.Write([]byte(identifier + "2"), ciphertext, nonce)
	}

	ciphertext, nonce, _ := dataStore.Read([]byte(identifier + "2"))

	deciphered, _ := cipher.Decrypt(key, ciphertext, nonce)
	fmt.Printf("%s\n", deciphered)
}
