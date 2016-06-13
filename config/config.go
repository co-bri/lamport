package config

import (
	"github.com/BurntSushi/toml"
)

// Config holds configuration options for lamport
type Config struct {
	Host, LamportPort, ElectionLibrary, RaftDir, RaftPort string
}

// ReadConfig reads configuration options out of a .toml file
func ReadConfig(tomlFilename string) (Config, error) {
	var config Config

	if _, err := toml.DecodeFile(tomlFilename, &config); err != nil {
		return config, err
	}
	return config, nil
}
