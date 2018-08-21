package fileio

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

// Fileio handles the read/write of ciphertext.
type Fileio struct {
	fileName string
}

// NewFileio is the factory method for Fileio.
func NewFileio(fileName string) *Fileio {
	io := new(Fileio)
	io.fileName = fileName

	_, err := os.Stat(fileName)

	// Check if file exists and has correct permissions.
	// Cause fast failure during startup.
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}

		file.Close()
	} else if os.IsPermission(err) {
		panic(err)
	}

	return io
}

// WriteToFile writes the ciphertext and corresponding nonce to a file in hex format seperated by comma.
func (fileio *Fileio) WriteToFile(ciphertext []byte, nonce []byte) {
	// file, err := os.OpenFile(fileio.fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	// TODO: Do not overwrite file.
	file, err := os.Create(fileio.fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	hexCipher := make([]byte, hex.EncodedLen(len(ciphertext)))
	hexNonce := make([]byte, hex.EncodedLen(len(nonce)))
	hex.Encode(hexCipher, ciphertext)
	hex.Encode(hexNonce, nonce)

	write(file, hexCipher)
	write(file, []byte{','})
	write(file, hexNonce)
	write(file, []byte{'\n'})
}

func write(file *os.File, payload []byte) {
	n, err := file.Write(payload)
	if err != nil {
		fmt.Printf("wrote %d bytes\n", n)
		panic(err)
	}
}

// ReadFromFile reads the ciphertext and corresponding nonce from a file.
func (fileio *Fileio) ReadFromFile() (ciphertext, nonce []byte) {
	file, err := os.Open(fileio.fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Scan() by default splits on "\n".
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	tokens := strings.Split(line, ",")
	ciphertext, _ = hex.DecodeString(tokens[0])
	nonce, _ = hex.DecodeString(tokens[1])

	return ciphertext, nonce
}
