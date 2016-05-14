package server

import (
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"os"
	"time"
)

func init() {
	log.SetOutput(os.Stdout)
}

func Run(ip int, host string) {
	log.Print("Initializing lamport...")
	connectZk()
	log.Printf("Running lamport on host: %s, port: %d", host, ip)
}

func connectZk() {
	log.Print("Connecting to Zookeper...")
	c, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		panic(err)
	}
	children, stat, _, err := c.ChildrenW("/")
	if err != nil {
		panic(err)
	}
	log.Printf("Connected to Zookeper: %+v %+v\n", children, stat)
}
