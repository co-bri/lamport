package node

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Distributed-Computing-Denver/lamport/config"
)

// Raft is an interface that contains operations that should
// be implemented by a raft node
type rafter interface {
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

type raftLdrWatch struct {
	conf config.Config
}

func newRaftLdrWatch(conf config.Config) raftLdrWatch {
	return raftLdrWatch{conf: conf}
}

func (r raftLdrWatch) leaderWatch(sigCh chan bool) chan bool {
	ch := make(chan bool)
	go r.startLeaderWatch(ch, sigCh)
	return ch
}

func (r raftLdrWatch) startLeaderWatch(ch chan bool, sigCh chan bool) {
	rn, err := NewRaftNode(r.conf.Host, r.conf.RaftPort, r.conf.Port, r.conf.RaftDir)
	if err != nil {
		panic(fmt.Errorf("Error creating raft node: %s", err))
	}
	connCh := make(chan net.Conn)
	go r.listen(connCh)

	if r.conf.Bootstrap != "" {
		r.writeJoinMessage(rn.RaftAddr(), r.conf.Bootstrap)
	}

	for {
		select {
		case c := <-connCh:
			log.Printf("Incoming connection from: %s", c.RemoteAddr())
			r.handleJoinMessage(c, rn)
		}
	}
}

// listen on ports for connections
func (r raftLdrWatch) listen(ch chan net.Conn) {
	ln, err := net.Listen("tcp", r.conf.Host+":"+r.conf.Port)
	defer ln.Close()
	if err != nil {
		panic(err)
	}
	log.Printf("Lamport listening on " + r.conf.Host + ":" + r.conf.Port)
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
func (r raftLdrWatch) handleJoinMessage(conn net.Conn, rn rafter) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading message: %s", err)
	}
	node := strings.TrimSpace(msg)
	if rn.State() == "Leader" {
		rn.Join(node)
	} else {
		if rn.LeaderAddr() != rn.LamportAddr() {
			r.writeJoinMessage(msg, rn.LeaderAddr())
		} else {
			log.Printf("Can't add %s to the cluster as this Raft node doesn't know the correct leader address", node)
		}
	}
	conn.Close()
}

// Write a join message to a given server. The join message consists
// of the raft host name and port commnunicated over a tcp socket.
func (r raftLdrWatch) writeJoinMessage(message, server string) {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Printf("Error connecting to %s: %s", server, err)
	}
	fmt.Fprintf(conn, message)
	conn.Close()
}
