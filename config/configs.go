package config

// Config stores all the parsed configurations for the applcation.
type Config struct {
	Storage StorageConfig
}

// StorageConfig stores the configurations for storing data.
type StorageConfig struct {
	Location string
	Driver   string
}
