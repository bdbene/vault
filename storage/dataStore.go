package storage

import "github.com/bdbene/vault/config"

//go:generate mockgen -destination=../mocks/mock_storage.go -package=mocks github.com/bdbene/vault/storage DataStore

// DataStore specifies the interface for storing data.
type DataStore interface {
	Write(identifier, ciphertext, nonce []byte) error
	Read(identifier []byte) (ciphertext, nonce []byte, err error)
	AlreadyExists(identifier []byte) (bool, error)
}

// DataStoreFactory is function type for creating creating DataStores
type DataStoreFactory func(conf *config.StorageConfig) (DataStore, error)

var dataStoreFactoryRegistry = make(map[string]DataStoreFactory)

// RegisterDataStoreFactory allows factory methods for DataStores to be saved
// so they can be referenced by name.
func RegisterDataStoreFactory(name string, factory DataStoreFactory) error {
	_, ok := dataStoreFactoryRegistry[name]
	if ok {
		return &DataStoreError{"DataStore named " + name + " already exists, cannot register it."}
	}

	dataStoreFactoryRegistry[name] = factory
	return nil
}

// CreateDataStore is a generic wrapper around DataStoreFactory functions to
// handle configurations.
func CreateDataStore(conf *config.StorageConfig) (DataStore, error) {
	dataStoreName := conf.Driver
	dataStoreFactory, ok := dataStoreFactoryRegistry[dataStoreName]
	if !ok {
		return nil, &DataStoreError{"Cannot load DataStore named" + dataStoreName + ". DataStore not registered."}
	}

	return dataStoreFactory(conf)
}
