package main

import (
	"flag"

	"github.com/Distributed-Computing-Denver/lamport/server"
)

func main() {
	port := flag.String("port", "5936", "The port on which lamport will listen for incoming connections")
	ip := flag.String("host", "127.0.0.1", "The host ip on which lamport will run")

	flag.Parse()
	server.Run(*ip, *port)
}
