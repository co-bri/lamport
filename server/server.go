// Package server provides methods for running lamport
package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Raft interface {
	State() string
}

func init() {
	log.SetOutput(os.Stdout)
}

// Run starts lamport on the given ip and hostname
func Run(host string, port string, r Raft) {
	log.Print("Initializing lamport...")
	connCh := make(chan net.Conn)
	go listen(host, port, connCh)

	for {
		select {
		case c := <-connCh:
			log.Printf("Incoming connection from: %s", c.RemoteAddr())
		}
	}
}

func listen(host string, port string, ch chan net.Conn) {
	ln, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		panic(err)
	}
	log.Printf("Lamport listening on " + host + ":" + port)
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		ch <- conn
	}
}

func handleConnection(conn net.Conn) {
	log.Print("Incoming connection made...")
}

func getNodeId(nodeId string) (id int) {
	id, err := strconv.Atoi(strings.Split(nodeId, "-")[1])
	if err != nil {
		panic(err)
	}
	return id
}
