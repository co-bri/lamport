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

// Raft is an interface that contains operations that should
// be implemented by a raft node
type Raft interface {
	Join(addr string) error
	LamportAddr() string
	Leader() string
	LeaderAddr() string
	RaftAddr() string
	State() string
}

func init() {
	log.SetOutput(os.Stdout)
}

// RunRaftServer starts lamport on the given ip and hostname using raft
// for leader election
func RunRaftServer(host string, port string, r Raft, joinServer string) {
	log.Print("Initializing lamport...")
	connCh := make(chan net.Conn)
	go listen(host, port, connCh)

	if joinServer != "" {
		writeJoinMessage(r.RaftAddr(), joinServer)
	}

	for {
		select {
		case c := <-connCh:
			log.Printf("Incoming connection from: %s", c.RemoteAddr())
			handleJoinMessage(c, r)
		}
	}
}

// listen on ports for connections
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

// Given a join message, either add the raft node address in the message to
// the cluster, or forward the message on to the leader of the cluster.
func handleJoinMessage(conn net.Conn, r Raft) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading message: %s", err)
	}
	node := strings.TrimSpace(msg)
	if r.State() == "Leader" {
		r.Join(node)
	} else {
		if r.LeaderAddr() != r.LamportAddr() {
			writeJoinMessage(msg, r.LeaderAddr())
		} else {
			log.Printf("Can't add %s to the cluster as this Raft node doesn't know the correct leader address", node)
		}
	}
	conn.Close()
}

// Write a join message to a given server. The join message consists
// of the raft host name and port commnunicated over a tcp socket.
func writeJoinMessage(message, server string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("Error connecting to %s: %s", server, err)
	}
	fmt.Fprintf(conn, message)
	conn.Close()
}
