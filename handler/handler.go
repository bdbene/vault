package handler

import (
	"github.com/bdbene/vault/cipher"
	"github.com/bdbene/vault/storage"
)

// Handler provides an interface for the internal workings of the applcation.
type Handler struct {
	DataStore storage.DataStore
}

// WriteSecret encrypts secret using password, then stores it using identifier as a lookup key.
func (handler *Handler) WriteSecret(identifier, password, secret []byte) error {
	key := cipher.CreateKey(password)
	ciphertext, nonce, err := cipher.Encrypt(key, secret)
	if err != nil {
		return err
	}

	err = handler.DataStore.Write(identifier, ciphertext, nonce)
	if err != nil {
		return err
	}

	return nil
}

// QuerySecret returns the deciphered secret from storage.
func (handler *Handler) QuerySecret(identifier, password []byte) ([]byte, error) {
	key := cipher.CreateKey(password)
	ciphertext, nonce, err := handler.DataStore.Read(identifier)
	if err != nil {
		return nil, err
	}

	secret, err := cipher.Decrypt(key, ciphertext, nonce)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
