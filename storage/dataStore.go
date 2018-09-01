package storage

import "github.com/bdbene/vault/config"

// DataStore specifies the interface for storing data.
type DataStore interface {
	Write(ciphertext, nonce []byte)
	Read() (ciphertext, nonce []byte)
}

// DataStoreFactory is function type for creating creating DataStores
type DataStoreFactory func(conf *config.StorageConfig) (DataStore, error)

var dataStoreFactoryRegistry = make(map[string]DataStoreFactory)

// RegisterDataStoreFactory allows factory methods for DataStores to be saved
// so they can be referenced by name.
func RegisterDataStoreFactory(name string, factory DataStoreFactory) {
	_, ok := dataStoreFactoryRegistry[name]
	if ok {
		panic("Attempting to register existing DriverFactory: " + name)
	}

	dataStoreFactoryRegistry[name] = factory
}

// CreateDataStore is a generic wrapper around DataStoreFactory functions to
// handle configurations.
func CreateDataStore(conf *config.StorageConfig) (DataStore, error) {
	dataStoreName := conf.Driver
	dataStoreFactory, ok := dataStoreFactoryRegistry[dataStoreName]
	if !ok {
		panic("Cannot create DataStore with alias: " + dataStoreName)
	}

	return dataStoreFactory(conf)
}
