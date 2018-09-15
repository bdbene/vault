package storage

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/bdbene/vault/config"
)

type dataStoreError struct {
	err string
}

func (e *dataStoreError) Error() string {
	return fmt.Sprintf("DataStore error: %s", e.err)
}

// Fileio handles the read/write of ciphertext.
type Fileio struct {
	fileName string
}

// NewFileio is the factory method for Fileio.
func NewFileio(conf *config.StorageConfig) (DataStore, error) {
	fileName := conf.Location
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

	return io, nil
}

// Writes the ciphertext and corresponding nonce to a file in hex format seperated by comma.
func (fileio *Fileio) Write(identifier []byte, ciphertext []byte, nonce []byte) error{
	
	// Check if the identifier already exists.
	exists, err := fileio.AlreadyExists(identifier)
	if err != nil {
		return err
	}
	if exists {
		return &dataStoreError{"Id already exists"}
	}

	file, err := os.OpenFile(fileio.fileName, os.O_APPEND | os.O_WRONLY, 0600)
	if err != nil {
		return &dataStoreError{err.Error()}
	}
	defer file.Close()

	hexCipher := make([]byte, hex.EncodedLen(len(ciphertext)))
	hexNonce := make([]byte, hex.EncodedLen(len(nonce)))
	hex.Encode(hexCipher, ciphertext)
	hex.Encode(hexNonce, nonce)

	// TODO: buffer this
	write(file, identifier)
	write(file, []byte{','})
	write(file, hexCipher)
	write(file, []byte{','})
	write(file, hexNonce)
	write(file, []byte{'\n'})

	return nil
}

func write(file *os.File, payload []byte)  {
	n, err := file.Write(payload)
	if err != nil {
		fmt.Printf("wrote %d bytes\n", n)
		panic(err)
	}
}

// Reads the ciphertext and corresponding nonce from a file.
func (fileio *Fileio) Read(identifier []byte) (ciphertext, nonce []byte, err error) {
	file, err := os.Open(fileio.fileName)
	if err != nil {
		return nil, nil, &DataStoreError{err.Error()}
	}
	defer file.Close()

	// Scan() by default splits on "\n".
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ",")
		lineId := tokens[0]
		
		if lineId == string(identifier) {
			ciphertext, err = hex.DecodeString(tokens[1])
			nonce, err = hex.DecodeString(tokens[2])
		}
	}

	return ciphertext, nonce, nil
}

func (fileio *Fileio) AlreadyExists(identifier []byte) (bool, error) {
	file, err := os.Open(fileio.fileName)
	if err != nil {
		return false, &DataStoreError{err.Error()}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, ",")
		lineId := tokens[0]
		
		if lineId == string(identifier) {
			return true, nil
		}
	}

	return false, nil
}

func init() {
	RegisterDataStoreFactory("FileIoDriver", NewFileio)
}
