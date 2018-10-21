package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bdbene/vault/config"
	"github.com/bdbene/vault/handler"
	"github.com/gorilla/mux"
)

type writeStruct struct {
	Password   string
	Secret     string
	Identifier string
}

type queryStruct struct {
	Identifier string
	Password   string
}

// RestServer that exposes writing and querying across network.
type RestServer struct {
	port       string
	handler    *handler.Handler
	router     *mux.Router
	serverCert string
	serverKey  string
	tlsEnabled bool
}

// NewServer configures a new restful service.
func NewServer(configs *config.ServiceConfig, handler *handler.Handler) *RestServer {
	server := new(RestServer)
	server.port = configs.Port
	server.tlsEnabled = configs.TLSEnabled
	server.serverCert = configs.ServerCert
	server.serverKey = configs.ServerKey
	server.handler = handler

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/secret/write", server.WriteSecret)
	router.HandleFunc("/secret/query", server.QuerySecret)
	server.router = router

	return server
}

// WriteSecret allows a client to encrypt and store a secret.
func (server *RestServer) WriteSecret(writer http.ResponseWriter, request *http.Request) {
	log.Print("INFO: Writing secret...")

	decoder := json.NewDecoder(request.Body)
	var payload writeStruct
	err := decoder.Decode(&payload)
	if err != nil {
		log.Printf("ERROR Bad request: %s", err.Error())
		fmt.Fprintf(writer, "Failure to decode payload.")
		return
	}

	errFuture := server.handler.RequestWrite(
		[]byte(payload.Identifier),
		[]byte(payload.Password),
		[]byte(payload.Secret))

	errResult := <-errFuture

	if errResult != nil {
		log.Printf("ERROR failure: %s", errResult.Error())
		fmt.Fprintf(writer, "Failure occured")
		return
	}

	fmt.Fprintf(writer, "Success!")
	log.Print("INFO: Secret successfully written.")
}

// QuerySecret allows a client to lookup a stored secret.
func (server *RestServer) QuerySecret(writer http.ResponseWriter, request *http.Request) {
	log.Print("INFO: Request for secret...")

	decoder := json.NewDecoder(request.Body)
	var payload queryStruct
	err := decoder.Decode(&payload)
	if err != nil {
		log.Print("ERROR: Failure to decode payload.")
		fmt.Fprintf(writer, "ERROR: bad request.")
		return
	}

	secretFuture, errFuture := server.handler.RequestQuery(
		[]byte(payload.Identifier),
		[]byte(payload.Password))

	select {
	case secret := <-secretFuture:
		fmt.Fprintf(writer, "%s", secret)
		log.Print("INFO: Secret successfully retreived.")
	case err := <-errFuture:
		log.Printf("Error: %s\n", err.Error())
		fmt.Fprintf(writer, "ERROR: failure.")
	}
}

// Listen for rest calls from clients.
func (server *RestServer) Listen() {
	log.Print("Running server...")

	if server.tlsEnabled {
		log.Fatal(http.ListenAndServeTLS(server.port, server.serverCert, server.serverKey, server.router))
	} else {
		log.Fatal(http.ListenAndServe(server.port, server.router))
	}
}
