// Package server provides methods for running lamport
package server

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

type Raft interface {
	Join(addr string) error
	Leader() string
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
			handleMessage(c, r)
		}
	}
}

func listen(host string, port string, ch chan net.Conn) {
	ln, err := net.Listen("tcp", host+":"+port)
	defer ln.Close()
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

func handleMessage(conn net.Conn, r Raft) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading message: %s", err)
	}
	node := strings.TrimSpace(msg)
	r.Join(node)
}
