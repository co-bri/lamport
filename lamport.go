package main

import (
	"flag"
	"fmt"

	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/Distributed-Computing-Denver/lamport/server"
)

func main() {
	joinServer := flag.String("join", "", "A Lamport server to join")
	tomlConfigFile := flag.String("tomlConfigFile", "lamport.toml", "The TOML file used to configure lamport")
	flag.Parse()

	config, err := config.ReadConfig(*tomlConfigFile)
	if err != nil {
		panic(fmt.Errorf("Error reading config file: %s", err))
	}

	server.Run(*joinServer, config)
}
