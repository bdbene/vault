package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"hash/fnv"
	"io"
	mathRand "math/rand"
	"fmt"
)

// CreateKey uses the password to create an AES key.
func CreateKey(password string) ([]byte, error) {
	hash := fnv.New64a()
	hash.Write([]byte(password))
	fmt.Printf("hash: %d\n", hash.Sum64())
	mathRand.Seed(int64(hash.Sum64()))

	// Create key.
	key := make([]byte, 64)
	_, err := rand.Read(key)
	if err != nil {
		return nil, &CipherError{"Key generation", err.Error()}
	}

	return key, nil
}

// Encrypt the plaintext using the given AES key, returns random nonce used as well.
func Encrypt(key []byte, plaintext []byte) ([]byte, []byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, &CipherError{"Encryption", err.Error()}
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, &CipherError{"Encryption", err.Error()}
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, &CipherError{"Encryption", err.Error()}
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// Decrypt the ciphertext using the given key and nonce
func Decrypt(key []byte, ciphertext []byte, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, &CipherError{"Decryption", err.Error()}
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, &CipherError{"Decryption", err.Error()}
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return plaintext, nil
}