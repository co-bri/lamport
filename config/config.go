package config

import (
	"github.com/BurntSushi/toml"
)

// Config stores configuration options for Lamport
type Config struct {
	Host      string
	Port      string
	Bootstrap string
	Zookeeper []string
	RaftDir   string
	RaftPort  string
}

// ReadConfig returns a Config created from the supplied config file
func ReadConfig(configFile string) (Config, error) {
	var config Config

	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return config, err
	}
	return config, nil
}
