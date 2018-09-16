package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/sha3"
)

// CreateKey generates a key based on the given secret
func CreateKey(password string) []byte {
	key := make([]byte, 32)

	sha3.ShakeSum256(key, []byte(password))
	return key
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
