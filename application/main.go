package main

import (
	"fmt"
	"os"

	"github.com/bdbene/vault/handler"
	"github.com/bdbene/vault/server"

	"github.com/bdbene/vault/config"
	"github.com/bdbene/vault/storage"
)

func main() {

	// Read configs.
	var conf config.Config
	err := config.GetConfigs(&conf)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		fmt.Printf("Cannot continue, shutting down.\n")
		os.Exit(1)
	}

	// Create DataStore based on configurations.
	dataStore, err := storage.CreateDataStore(&conf.Storage)
	if err != nil {
		panic(err)
	}

	handler := handler.NewHandler(dataStore, &conf.Handler)
	server := server.NewServer(&conf.Server, handler)

	handler.ProcessRequests()
	server.Listen()
}
