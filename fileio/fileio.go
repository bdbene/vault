package fileio

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// WriteToFile writes the ciphertext and corresponding nonce to a file in hex format seperated by comma.
func WriteToFile(ciphertext []byte, nonce []byte) {
	hexCipher := make([]byte, hex.EncodedLen(len(ciphertext)))
	hexNonce := make([]byte, hex.EncodedLen(len(nonce)))
	hex.Encode(hexCipher, ciphertext)
	hex.Encode(hexNonce, nonce)

	file, err := os.Create("test")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	write(file, hexCipher)
	write(file, []byte{','})
	write(file, hexNonce)
}

func write(file *os.File, payload []byte) {
	n, err := file.Write(payload)
	if err != nil {
		fmt.Printf("wrote %d bytes\n", n)
		panic(err.Error())
	}
}

// ReadFromFile reads the ciphertext and corresponding nonce from a file.
func ReadFromFile() (ciphertext, nonce []byte) {
	file, err := os.Open("test")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	// Scan() by default splits on "\n"
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	tokens := strings.Split(line, ",")
	ciphertext, _ = hex.DecodeString(tokens[0])
	nonce, _ = hex.DecodeString(tokens[1])

	return ciphertext, nonce
}
