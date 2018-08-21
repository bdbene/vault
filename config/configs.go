package config

type Config struct {
	Storage StorageConfig
}

type StorageConfig struct {
	Location string
}
