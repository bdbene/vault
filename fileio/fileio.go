package fileio

import (
	"encoding/hex"
	"fmt"
	"os"
)

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
