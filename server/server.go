// Package server provides methods for running lamport
package server

import (
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	zkRoot  = "/lamport"
	zkNodes = zkRoot + "/nodes"
)

var acl = zk.WorldACL(zk.PermAll)

func init() {
	log.SetOutput(os.Stdout)
}

// Run starts lamport on the given ip and hostname
func Run(ip string, host string) {
	log.Print("Initializing lamport...")
	connCh := make(chan net.Conn)
	zkCh, zkConn := connectZk(host, ip)
	go listen(ip, host, connCh)

	for {
		select {
		case e := <-zkCh:
			if e.Err != nil {
				panic(e.Err)
			} else {
				var err error
				log.Print("Zookeeper watch event", e)
				_, _, zkCh, err = zkConn.ChildrenW(zkNodes)
				if err != nil {
					panic(err)
				}
			}
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

func connectZk(host string, port string) (<-chan zk.Event, *zk.Conn) {
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		panic(err)
	}

	createParentZNodes(conn)
	node := createZNode(conn, host, port)
	return leaderWatch(conn, node), conn
}

func leaderWatch(conn *zk.Conn, nodeId string) <-chan zk.Event {
	id := getNodeId(nodeId)
	watchId := id

	nodes, _, ch, err := conn.ChildrenW(zkNodes)
	if err != nil {
		panic(err)
	}

	seq := make([]int, len(nodes))
	for i, v := range nodes {
		nodeId := getNodeId(v)
		seq[i] = nodeId
	}
	sort.Ints(seq)
	log.Printf("Current leader is sequence number: %d", seq[0])

	for _, v := range seq {
		if v >= id {
			break
		} else if v < watchId {
			watchId = v
		}
	}
	if watchId == id {
		log.Printf("Entering leader mode")
	} else {
		log.Printf("Entering follower mode, watching id: %d for changes", watchId)
	}
	return ch
}

func createZNode(conn *zk.Conn, host string, port string) (path string) {
	data := []uint8(host + ":" + port)
	path, err := conn.CreateProtectedEphemeralSequential(zkNodes+"/", data, acl)
	if err != nil {
		panic(err)
	}
	log.Printf("Created znode %s", path)
	return path
}

func createParentZNodes(conn *zk.Conn) {
	exists, _, err := conn.Exists(zkRoot)
	if err != nil {
		panic(err)
	}

	if !exists {
		log.Print("Creating parent znode 'lamport' in zookeeper")
		_, err := conn.Create(zkRoot, nil, 0, acl)
		if err != nil {
			panic(err)
		}
	}

	exists, _, err = conn.Exists(zkNodes)
	if err != nil {
		panic(err)
	}

	if !exists {
		log.Print("Creating parent znode 'nodes' in zookeeper")
		_, err := conn.Create(zkNodes, nil, 0, acl)
		if err != nil {
			panic(err)
		}
	}
}

func getNodeId(nodeId string) (id int) {
	id, err := strconv.Atoi(strings.Split(nodeId, "-")[1])
	if err != nil {
		panic(err)
	}
	return id
}
