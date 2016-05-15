package server

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"net"
	"os"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
}

func Run(ip string, host string) {
	log.Print("Initializing lamport...")
	connectZk()
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

func connectZk() {
	log.Print("Connecting to Zookeper...")
	_, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to Zookeper")
}
