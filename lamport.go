package main

import (
	"flag"
	"fmt"

	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/Distributed-Computing-Denver/lamport/server"
)

func main() {
	configFile := flag.String("configFile", "lamport.toml", "Lamport config file")
	flag.Parse()

	config, err := config.ReadConfig(*configFile)
	if err != nil {
		panic(fmt.Errorf("Error reading config file: %s", err))
	}

	server.Run(config, make(chan bool))
}
