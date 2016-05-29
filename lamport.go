package main

import (
	"flag"
	"github.com/Distributed-Computing-Denver/lamport/server"
)

func main() {
	ip := flag.String("port", "5936", "The port on which lamport will listen for incoming connections")
	host := flag.String("host", "127.0.0.1", "The host ip on which lamport will run")

	raftDir := flag.String("raftDir", ".raft", "The directory used for Raft storage")
	raftNode := server.NewRaftNode(*raftDir)
	flag.Parse()
	server.Run(*ip, *host, raftNode)
}
