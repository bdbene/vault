package handler

import (
	"log"

	"github.com/bdbene/vault/cipher"
	"github.com/bdbene/vault/config"
	"github.com/bdbene/vault/storage"
)

type write struct {
	identifier []byte
	password   []byte
	secret     []byte
	err        chan error
}

type query struct {
	identifier []byte
	password   []byte
	result     chan string
	err        chan error
}

// Handler provides an interface for the internal workings of the applcation.
type Handler struct {
	dataStore     storage.DataStore
	writeRequests chan write
	queryRequests chan query
}

// NewHandler returns a Handler with the correct configurations and DataStore.
func NewHandler(dataStore storage.DataStore, configs *config.HandlerConfig) *Handler {
	handler := new(Handler)
	handler.dataStore = dataStore
	handler.writeRequests = make(chan write, configs.WriteBufferSize)
	handler.queryRequests = make(chan query, configs.QueryBufferSize)

	return handler
}

// RequestWrite queues a request to write a secret. Returned channel is used as a future for when
// the request is complete. nil error if the write was succesful, error otherwise.
func (handler *Handler) RequestWrite(identifier, password, secret []byte) chan error {
	errFuture := make(chan error)
	handler.writeRequests <- write{identifier, password, secret, errFuture}

	return errFuture
}

// RequestQuery queues a request to query a secret. Returned channels will
// contain either the result or an error.
func (handler *Handler) RequestQuery(identifier, password []byte) (chan string, chan error) {
	resultFuture := make(chan string)
	errFuture := make(chan error)
	handler.queryRequests <- query{identifier, password, resultFuture, errFuture}

	return resultFuture, errFuture
}

// ProcessRequests will start processing all requests concurrently.
func (handler *Handler) ProcessRequests() {
	go handler.process()
}

func (handler *Handler) process() {
	for {
		select {
		case w := <-handler.writeRequests:
			log.Printf("Buffered requests: %d\n", len(handler.writeRequests))
			handler.processWrite(w)
		case q := <-handler.queryRequests:
			log.Printf("Buffered requests: %d\n", len(handler.queryRequests))
			handler.processQuery(q)
		}
	}
}

func (handler *Handler) processWrite(w write) {
	key := cipher.CreateKey(w.password)
	ciphertext, nonce, err := cipher.Encrypt(key, w.secret)
	if err != nil {
		w.err <- err
		return
	}

	err = handler.dataStore.Write(w.identifier, ciphertext, nonce)
	if err != nil {
		w.err <- err
		return
	}

	w.err <- nil
}

func (handler *Handler) processQuery(q query) {
	key := cipher.CreateKey(q.password)
	ciphertext, nonce, err := handler.dataStore.Read(q.identifier)
	if err != nil {
		q.err <- err
		return
	}

	secret, err := cipher.Decrypt(key, ciphertext, nonce)
	if err != nil {
		q.err <- err
		return
	}

	q.result <- string(secret)
}
