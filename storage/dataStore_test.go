package storage

import (
	"testing"

	"github.com/bdbene/vault/config"
)

func TestRegisterFactory_FactoryRegister(t *testing.T) {
	factoryName := "factoryName"
	err := RegisterDataStoreFactory(factoryName, NewFileio)

	if err != nil {
		t.Errorf("Error registering factory: %s", err.Error())
		return
	}

	_, ok := dataStoreFactoryRegistry[factoryName]

	if !ok {
		t.Error("Expected factory method to be registered, it is not.")
		return
	}
}

func TestRegisterExistingFactory_CannotRegister(t *testing.T) {
	factoryName := "factoryName"
	RegisterDataStoreFactory(factoryName, NewFileio)

	err := RegisterDataStoreFactory(factoryName, NewFileio)

	if err == nil {
		t.Error("Expected error when existing factory method registered again.")
		return
	}
}

func TestCreateDataStore_CorrectDataStoreCreated(t *testing.T) {
	storeName := "storeName"
	configs := &config.StorageConfig{"/dev/null", storeName}
	RegisterDataStoreFactory(storeName, NewFileio)

	_, err := CreateDataStore(configs)

	if err != nil {
		t.Errorf("Error creating datastore when none expected: %s", err.Error())
		return
	}
}
