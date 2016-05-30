package main

import (
	"flag"
	"fmt"
	"github.com/Distributed-Computing-Denver/lamport/server"
)

func main() {
	port := flag.String("port", "5936", "The port on which lamport will listen for incoming connections")
	host := flag.String("host", "127.0.0.1", "The host ip on which lamport will run")
	raftDir := flag.String("raftDir", ".raft", "The directory used for Raft storage")
	raftIp := flag.String("raftPort", "8500", "The port on which the raft protocol will communicate")

	raftNode, err := server.NewRaftNode(*host, *raftIp, *raftDir)
	if err != nil {
		panic(fmt.Errorf("Error creating raft node: %s", err))
	}
	flag.Parse()
	server.Run(*host, *port, raftNode)
}
