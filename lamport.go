package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/Distributed-Computing-Denver/lamport/node"
)

func main() {
	configFile := flag.String("configFile", "lamport.toml", "Lamport config file")
	flag.Parse()

	config, err := config.ReadConfig(*configFile)
	if err != nil {
		panic(fmt.Errorf("Error reading config file: %s", err))
	}

	sigCh := make(chan bool)
	go node.Run(config, sigCh)

	// handle SIGINT, notify node, wait for confirm to exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	sigCh <- true
	<-sigCh
}
