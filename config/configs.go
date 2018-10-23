package config

// Config stores all the parsed configurations for the applcation.
type Config struct {
	Storage StorageConfig
	Server  ServiceConfig
	Handler HandlerConfig
}

// StorageConfig stores the configurations for storing data.
type StorageConfig struct {
	Location string
	Driver   string
}

// ServiceConfig stores the configurations for running restful service.
type ServiceConfig struct {
	Port       string
	ServerCert string
	ServerKey  string
	TLSEnabled bool
}

// HandlerConfig stores the configurations for handling requests.
type HandlerConfig struct {
	WriteBufferSize int
	QueryBufferSize int
}
