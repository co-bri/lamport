package main

import (
	"flag"
	"fmt"
	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/Distributed-Computing-Denver/lamport/server"
)

func main() {
	tomlConfigFile := flag.String("tomlConfigFile", "lamport.toml", "The TOML file used to configure lamport")
	flag.Parse()

	config, err := config.ReadConfig(*tomlConfigFile)
	if err != nil {
		panic(fmt.Errorf("Error reading config file: %s", err))
	}

	raftNode, err := server.NewRaftNode(config.Host, config.RaftPort, config.JoinPort, config.RaftDir)
	if err != nil {
		panic(fmt.Errorf("Error creating raft node: %s", err))
	}

	server.Run(config.Host, config.JoinPort, raftNode)
}
