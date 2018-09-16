package cipher

import (
	"encoding/hex"
	"testing"
)

func TestCreateKey_32ByteKeyCreated(t *testing.T) {
	key := CreateKey([]byte("password"))
	expectedLength := 32

	keyLength := len(key)

	if keyLength != expectedLength {
		t.Errorf("Key should be %d bytes, got: %d", expectedLength, keyLength)
	}
}

func TestCreateKey_SamePasswordSameKey(t *testing.T) {
	password := []byte("password")
	key1 := CreateKey((password))
	key2 := CreateKey([]byte("password"))

	if !equalSlice(key1, key2) {
		t.Error("Keys generated from the same password should be equal.")
		return
	}
}

func TestEncrypt_NonceShouldBeRandom(t *testing.T) {
	key := CreateKey([]byte("password"))
	secret := []byte("mySecret")

	cipher1, nonce1, err := Encrypt(key, secret)
	if err != nil {
		t.Error("Unexpected error during testing.")
		return
	}

	cipher2, nonce2, err := Encrypt(key, secret)
	if err != nil {
		t.Error("Unexpected error during testing.")
		return
	}

	if equalSlice(nonce1, nonce2) {
		t.Error("Nonces should be randomly generated and unique.")
		return
	}

	if equalSlice(cipher1, cipher2) {
		t.Error("Ciphers should be unique if random nonce is used.")
		return
	}
}

func TestDecrypt_AssumeCipherAndNonce(t *testing.T) {
	hexCipherText := "d5ca155f688607d6f9bdcaca72f32c2a0d1b2efed03176a0d5835526"
	hexNonce := "9345fef7a66fc7c67d47cfe1"
	password := []byte("NewPass5")
	secret := []byte("Hello world!")

	key := CreateKey(password)
	cipherText, _ := hex.DecodeString(hexCipherText)
	nonce, _ := hex.DecodeString(hexNonce)

	plaintext, err := Decrypt(key, cipherText, nonce)

	if err != nil {
		t.Errorf("Failed to decrypt: %s", err.Error())
		return
	}

	if !equalSlice(secret, plaintext) {
		t.Errorf("Deciphered text expected %s, got %s.", secret, plaintext)
		return
	}
}

func equalSlice(a, b []byte) bool {

	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
