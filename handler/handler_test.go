package handler

import (
	"encoding/hex"
	"testing"

	"github.com/bdbene/vault/config"
	"github.com/bdbene/vault/mocks"
	"github.com/golang/mock/gomock"
)

var handler *Handler
var mockStorage *mocks.MockDataStore

func setup(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockStorage = mocks.NewMockDataStore(mockCtrl)
	fakeConfig := &config.HandlerConfig{WriteBufferSize: 100, QueryBufferSize: 100}
	handler = NewHandler(mockStorage, fakeConfig)
}

type testError struct {
	errMsg string
}

func (err *testError) Error() string {
	return err.errMsg
}

func TestRequestWrite_ReturnsResultFuture(t *testing.T) {
	setup(t)
	identifier := []byte("Key")
	password := []byte("Password")
	secret := []byte("Secret")

	errFuture := handler.RequestWrite(identifier, password, secret)

	if errFuture == nil {
		t.Errorf("Expected future which will contain result.")
	}
}

func TestRequestWrite_CorrespondingFutureQueued(t *testing.T) {
	setup(t)
	identifier := []byte("Key")
	password := []byte("Password")
	secret := []byte("Secret")

	errFuture := handler.RequestWrite(identifier, password, secret)

	queuedWork := <-handler.writeRequests

	if errFuture != queuedWork.err {
		t.Error("Expected queued future to be the same as the callers's.")
	}
}

func TestProcessWrite_WriteToDataStore(t *testing.T) {
	setup(t)
	identifier := []byte("Key")
	password := []byte("Password")
	secret := []byte("Secret")
	writeWork := write{identifier: identifier, password: password, secret: secret, err: make(chan error)}

	mockStorage.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

	go handler.processWrite(writeWork)
	result := <-writeWork.err

	if result != nil {
		t.Errorf("Unexpected error when processing write: %s\n", result.Error())
	}
}

func TestProcessWrite_ErrorWritingToDataStore(t *testing.T) {
	setup(t)
	identifier := []byte("Key")
	password := []byte("Password")
	secret := []byte("Secret")
	writeWork := write{identifier: identifier, password: password, secret: secret, err: make(chan error)}
	err := &testError{"Test error."}

	mockStorage.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(err).Times(1)

	go handler.processWrite(writeWork)
	result := <-writeWork.err

	if result != err {
		t.Error("Expected error from data store, none given.")
	}
}

func TestRequestQuery_QueryWorkQueued(t *testing.T) {
	setup(t)
	identifier := []byte("Key")
	password := []byte("Password")

	resultFuture, errFuture := handler.RequestQuery(identifier, password)

	queuedWork := <-handler.queryRequests

	if resultFuture != queuedWork.result {
		t.Error("Expected result future to be the same as queued result future.")
	}

	if errFuture != queuedWork.err {
		t.Error("Expected error future to be the same as the queued error future.")
	}
}
func TestProcessQuery_QueuedWorkProcessed(t *testing.T) {
	setup(t)
	identifier := []byte("Key")
	password := []byte("NewPass5")
	secret := "Hello world!"
	ciphertext, _ := hex.DecodeString("d5ca155f688607d6f9bdcaca72f32c2a0d1b2efed03176a0d5835526")
	nonce, _ := hex.DecodeString("9345fef7a66fc7c67d47cfe1")
	queryWork := query{identifier: identifier, password: password, result: make(chan string), err: make(chan error)}

	mockStorage.EXPECT().Read(identifier).Return(ciphertext, nonce, nil).Times(1)

	go handler.processQuery(queryWork)

	select {
	case err := <-queryWork.err:
		t.Errorf("Unexpected error processing query: %s\n", err.Error())
	case result := <-queryWork.result:
		if result != secret {
			t.Errorf("Expected result: %s, received %s.\n", secret, result)
		}
	}
}

func TestProcessQuery_FailureFromDataStoreSet(t *testing.T) {
	setup(t)
	identifier := []byte("Key")
	password := []byte("NewPass5")
	queryWork := query{identifier: identifier, password: password, result: make(chan string), err: make(chan error)}

	mockStorage.EXPECT().Read(identifier).Return(nil, nil, &testError{"Test error."}).Times(1)

	go handler.processQuery(queryWork)

	select {
	case <-queryWork.err:
	case result := <-queryWork.result:
		t.Errorf("Unexpected result: %s, expecting error", result)
	}
}
