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
	zkCand  = "cand"
)

var acl = zk.WorldACL(zk.PermAll)

func init() {
	log.SetOutput(os.Stdout)
}

// Run starts lamport on the given ip and port
func Run(ip string, port string) {
	log.Print("Starting lamport...")

	log.Print("Connecting to Zookeeper on 127.0.0.1") // TODO make ZK host/port configurable
	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		panic(err)
	}

	createParentZNodes(conn)
	zn := createCandidateZNode(conn, ip, port)
	ch := watchCandidateNode(conn, strings.Split(zn, "/")[3])

	for {
		if e := <-ch; e.Err != nil {
			log.Fatal(e.Err)
		} else {
			log.Printf("Zookeeper watch event: %v", e)
			_, _, ch, err = conn.ChildrenW(zkNodes)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// sets a watch on the candidate node having seq number prior to supplied node id
func watchCandidateNode(conn *zk.Conn, nodeID string) <-chan zk.Event {
	wchID := nodeID

	nds, _, ch, err := conn.ChildrenW(zkNodes)
	if err != nil {
		panic(err)
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
func createCandidateZNode(conn *zk.Conn, host string, port string) (path string) {
	data := []uint8(host + ":" + port)
	path, err := conn.Create(zkNodes+"/"+zkCand, data, zk.FlagEphemeral|zk.FlagSequence, acl)
	if err != nil {
		panic(err)
	}
	log.Printf("Created candidate znode %s", path)
	return path
}

// creates the required zknodes if not present
func createParentZNodes(conn *zk.Conn) {
	exists, _, err := conn.Exists(zkRoot)
	if err != nil {
		panic(err)
	}

	if !exists {
		log.Printf("Root znode not found, creating %s in zookeeper", zkRoot)
		if _, err := conn.Create(zkRoot, nil, 0, acl); err != nil {
			panic(err)
		}
	}

	exists, _, err = conn.Exists(zkNodes)
	if err != nil {
		panic(err)
	}

	if !exists {
		log.Printf("Leader election parent node not found, Creating %s in zookeeper", zkNodes)
		if _, err := conn.Create(zkNodes, nil, 0, acl); err != nil {
			panic(err)
		}
	}
}
