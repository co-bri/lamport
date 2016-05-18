// Package server provides methods for running lamport
package server

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var acl = zk.WorldACL(zk.PermAll)

func init() {
	log.SetOutput(os.Stdout)
}

// Run starts lamport on the given ip and hostname
func Run(ip string, host string) {
	log.Print("Initializing lamport...")

	connectZk(host, ip)
	listen(ip, host)
}

func listen(ip string, host string) {
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	log.Print("Incoming connection made...")
}

func connectZk(host string, port string) {
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		panic(err)
	}

	createParentZNodes(conn)
	node := createZNode(conn, host, port)
	leaderWatch(conn, node)
}

func leaderWatch(conn *zk.Conn, node string) {
	id, err := strconv.Atoi(strings.Split(node, "-")[1])
	if err != nil {
		panic(err)
	}

	nodes, _, _, err := conn.ChildrenW("/lamport/nodes")
	if err != nil {
		panic(err)
	}
	watchId := 0
	for _, v := range nodes {
		nodeId, err := strconv.Atoi(strings.Split(v, "-")[1])
		if err != nil {
			panic(err)
		}
		if (nodeId < id) && (nodeId > watchId) {
			watchId = nodeId
		}
	}
	if watchId == 0 {
		log.Print("Running in single node cluster, leader by default")
	} else {
		log.Printf("Watching %s for leader changes")
	}
	//TODO - handle leader changes
}

func createZNode(conn *zk.Conn, host string, port string) (path string) {
	data := []uint8(host + ":" + port)
	path, err := conn.CreateProtectedEphemeralSequential("/lamport/nodes/", data, acl)
	if err != nil {
		panic(err)
	}
	log.Printf("Created znode %s", path)
	return path
}

func createParentZNodes(conn *zk.Conn) {
	// Check if the parent nodes exist
	exists, _, err := conn.Exists("/lamport")
	if err != nil {
		panic(err)
	}

	if !exists {
		log.Print("Creating parent znode 'lamport' in zookeeper")
		_, err := conn.Create("/lamport", nil, 0, acl)
		if err != nil {
			panic(err)
		}
	}

	exists, _, err = conn.Exists("/lamport/nodes")
	if err != nil {
		panic(err)
	}

	if !exists {
		log.Print("Creating parent znode 'nodes' in zookeeper")
		_, err := conn.Create("/lamport/nodes", nil, 0, acl)
		if err != nil {
			panic(err)
		}
	}
}
