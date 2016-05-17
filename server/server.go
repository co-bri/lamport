package server

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var acl = zk.WorldACL(zk.PermAll)

func init() {
	log.SetOutput(os.Stdout)
}

func Run(ip string, host string) {
	log.Print("Initializing lamport...")
	ch := make(chan string)

	connectZk(ch, host, ip)
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

func connectZk(ch chan<- string, host string, port string) {
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	createParentZNodes(conn)

	data := []uint8(host + ":" + port)
	path, err := conn.CreateProtectedEphemeralSequential("/lamport/nodes/", data, acl)
	if err != nil {
		panic(err)
	}
	log.Printf("Created znode %s", path)

	nodes, _, _, err := conn.ChildrenW("/lamport/nodes")
	if err != nil {
		panic(err)
	}
	log.Print(nodes)
	// TODO - setup watch on node adjacent node
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
