package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"hash/fnv"
	"io"
	mathRand "math/rand"
)

// CreateKey uses the password to create an AES key
func CreateKey(password string) []byte {
	hash := fnv.New64a()
	hash.Write([]byte(password))
	mathRand.Seed(int64(hash.Sum64()))

	// Create key.
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err.Error())
	}

	return key
}

// Encrypt the plaintext using the given AES key
func Encrypt(key []byte, plaintext []byte) ([]byte, []byte) {
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

// Decrypt the ciphertext using the given key and nonce
func Decrypt(key []byte, ciphertext []byte, nonce []byte) []byte {
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
