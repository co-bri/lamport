package server

import (
	"fmt"

	"github.com/Distributed-Computing-Denver/lamport/config"
)

// Run runs a lamport server, initializing it based upon the config file
// and a server to join
func Run(joinServerAddr string, config config.Config) {
	if config.ElectionLibrary == "Raft" {
		raftNode, err := NewRaftNode(config.Host, config.RaftPort, config.LamportPort, config.RaftDir)
		if err != nil {
			panic(fmt.Errorf("Error creating raft node: %s", err))
		}
		RunRaftServer(config.Host, config.LamportPort, raftNode, joinServerAddr)
	} else if config.ElectionLibrary == "Zookeeper" {
		ch := make(chan bool)
		RunZkServer(config.Host, config.LamportPort, ch)
	} else {
		panic(fmt.Errorf("Unsupported election library! Must be 'Raft' or 'Zookeeper', not '%s'!", config.ElectionLibrary))
	}

}
