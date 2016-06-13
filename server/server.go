// Package server provides methods for running lamport
package server

import (
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	zkRoot  = "/lamport"
	zkNodes = zkRoot + "/nodes"
	zkPrfx  = "node"
)

var acl = zk.WorldACL(zk.PermAll)

func init() {
	log.SetOutput(os.Stdout)
}

// Run starts lamport on the given ip and port
func Run(ip string, port string) {
	log.Print("Starting lamport...")

	zkConn, sCh, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		log.Fatalf("Error connecting to zookeeper: %s", err)
	}

	// setup required znodes, set watch for leader election
	createParentZNodes(zkConn)
	zn := createZNode(zkConn, ip, port)
	wCh := watchNode(zkConn, strings.Split(zn, "/")[3])

	for {
		select {
		case se := <-sCh:
			if se.State == zk.StateUnknown || se.State == zk.StateDisconnected {
				log.Printf("Zookeeper %s, state: %s, server: %s ", se.Type, se.State, se.Server)
			}
		case we := <-wCh:
			log.Printf("Zookeeper watch event: %s", we)
			_, _, wCh, err = zkConn.ChildrenW(zkNodes)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// sets a watch on the candidate node having seq number prior to supplied node id
func watchNode(conn *zk.Conn, nodeID string) <-chan zk.Event {
	wchID := nodeID

	nds, _, ch, err := conn.ChildrenW(zkNodes)
	if err != nil {
		log.Fatal(err)
	}

	sort.Strings(nds)
	log.Printf("Current leader id: %s", nds[0])

	for _, id := range nds {
		if id >= nodeID {
			break
		} else if id < wchID {
			wchID = id
		}
	}

	if wchID == nodeID {
		log.Printf("Entering leader mode")
	} else {
		log.Printf("Entering follower mode, watching candidate node: %s for changes", wchID)
	}
	return ch
}

// creates a candidate znode for this lamport instance
func createZNode(conn *zk.Conn, host string, port string) (path string) {
	data := []uint8(host + ":" + port)
	path, err := conn.Create(zkNodes+"/"+zkPrfx, data, zk.FlagEphemeral|zk.FlagSequence, acl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created candidate znode %s", path)
	return path
}

// creates the required parent znodes if not present
func createParentZNodes(conn *zk.Conn) {
	exists, _, err := conn.Exists(zkRoot)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Printf("Root znode not found, creating %s in zookeeper", zkRoot)
		if _, err := conn.Create(zkRoot, nil, 0, acl); err != nil {
			log.Fatal(err)
		}
	}

	exists, _, err = conn.Exists(zkNodes)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Printf("Leader election parent node not found, Creating %s in zookeeper", zkNodes)
		if _, err := conn.Create(zkNodes, nil, 0, acl); err != nil {
			log.Fatal(err)
		}
	}
}
