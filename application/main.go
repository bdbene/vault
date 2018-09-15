package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"encoding/hex"

	"github.com/bdbene/vault/cipher"
	"github.com/bdbene/vault/config"
	"github.com/bdbene/vault/storage"
	"github.com/gorilla/mux"
)

type writeStruct struct {
	Password string
	Text string
	Identifier string
}

type queryStruct struct {
	Identifier string
	Password string
}

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


	writeHandler := func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("Write\n")
	
		decoder := json.NewDecoder(request.Body)
		var t writeStruct
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		key, _ := cipher.CreateKey(t.Password)
		PrintHex(key)
		ciphertext, nonce, _ := cipher.Encrypt(key, []byte(t.Text))
		dataStore.Write([]byte(t.Identifier), ciphertext, nonce)
		fmt.Fprintf(writer, "Success: %s", html.EscapeString(request.URL.Path))
	}

	queryHandler := func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("Query\n")

		decoder := json.NewDecoder(request.Body)
		var t queryStruct
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		key, _ := cipher.CreateKey(t.Password)
		ciphertext, nonce, err := dataStore.Read([]byte(t.Identifier))
		if err != nil {
			panic(err)
		}

		PrintHex(key)
		fmt.Printf("%s %s\n", ciphertext, nonce)

		deciphered, _ := cipher.Decrypt(key, ciphertext, nonce)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(writer, "Secret: %s", deciphered)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/secret/write", writeHandler)
	router.HandleFunc("/secret/query", queryHandler)


	fmt.Printf("Running server...\n")
	log.Fatal(http.ListenAndServe(":8080", router))
	

	// Encrypt.
	/*{
		ciphertext, nonce, _ := cipher.Encrypt(key, text)
		dataStore.Write([]byte(identifier), ciphertext, nonce)
		ciphertext, nonce, _ = cipher.Encrypt(key, text)
		dataStore.Write([]byte(identifier + "2"), ciphertext, nonce)
	}

	ciphertext, nonce, _ := dataStore.Read([]byte(identifier + "2"))

	deciphered, _ := cipher.Decrypt(key, ciphertext, nonce)
	fmt.Printf("%s\n", deciphered)*/
}

func PrintHex(txt []byte) {
	hexTxt := make([]byte, hex.EncodedLen(len(txt)))
	hex.Encode(hexTxt, txt)
	fmt.Printf("%s\n", hexTxt)
}