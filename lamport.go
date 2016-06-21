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

	if config.ElectionLibrary == "Raft" {
		raftNode, err := server.NewRaftNode(config.Host, config.RaftPort, config.LamportPort, config.RaftDir)
		if err != nil {
			panic(fmt.Errorf("Error creating raft node: %s", err))
		}
		server.RunRaftServer(config.Host, config.LamportPort, raftNode, *joinServer)
	} else if config.ElectionLibrary == "Zookeeper" {
		ch := make(chan bool)
		server.RunZkServer(config.Host, config.LamportPort, ch)
	} else {
		panic(fmt.Errorf("Unsupported election library! Must be 'Raft' or 'Zookeeper', not '%s'!", config.ElectionLibrary))
	}
}
