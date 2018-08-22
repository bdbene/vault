package main

import (
	"fmt"
	"os"

	"github.com/bdbene/vault/cipher"
	"github.com/bdbene/vault/config"
	"github.com/bdbene/vault/fileio"
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

	// Set configurations.
	io := fileio.NewFileio(conf.Storage.Location)

	key := cipher.CreateKey(password)

	// Encrypt.
	{
		ciphertext, nonce := cipher.Encrypt(key, text)
		io.WriteToFile(ciphertext, nonce)
	}

	ciphertext, nonce := io.ReadFromFile()

	deciphered := cipher.Decrypt(key, ciphertext, nonce)
	fmt.Printf("%s\n", deciphered)
}
