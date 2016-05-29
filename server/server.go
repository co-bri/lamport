// Package server provides methods for running lamport
package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func init() {
	log.SetOutput(os.Stdout)
}

// Run starts lamport on the given ip and hostname
func Run(ip string, host string, r *RaftNode) {
	log.Print("Initializing lamport...")
	connCh := make(chan net.Conn)
	go listen(ip, host, connCh)

	for {
		select {
		case c := <-connCh:
			log.Printf("Incoming connection from: %s", c.RemoteAddr())
		}
	}
}

func listen(ip string, host string, ch chan net.Conn) {
	ln, err := net.Listen("tcp", host+":"+ip)
	if err != nil {
		panic(err)
	}
	log.Printf("Lamport listening on " + host + ":" + ip)
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
