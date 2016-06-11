// Package server provides methods for running lamport
package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Raft interface {
	Join(addr string) error
	LamportAddr() string
	Leader() string
	LeaderAddr() string
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
	if r.State() == "Leader" {
		r.Join(node)
	} else {
		if r.LeaderAddr() != r.LamportAddr() {
			writeMessage(msg, r.LeaderAddr())
		} else {
			log.Printf("Can't add %s to the cluster as this Raft node doesn't know the correct leader address", node)
		}
	}
	conn.Close()
}

func writeMessage(message, server string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("Error connecting to %s: %s", server, err)
	}
	fmt.Fprintf(conn, message)
	conn.Close()
}
