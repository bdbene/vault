package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

// GetConfigs loads configs from config.tml into conf
func GetConfigs(conf *Config) error {
	data, err := readConfigurations()
	if err != nil {
		return err
	}

	toml.Decode(data, conf)
	return nil
}

func readConfigurations() (string, error) {
	configurations, err := ioutil.ReadFile("./config.tml")
	if err != nil {
		return "", &ConfigError{err.Error()}
	}

	return string(configurations[:]), nil
}
