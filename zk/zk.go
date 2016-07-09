package zk

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

// WatchLeader uses zookeeper to elect a leader, and returns
// a channel that recieves leader change events
func WatchLeader(ip string, port string, zkh []string, sigCh <-chan bool) chan bool {
	ch := make(chan bool)
	go watchLeader(ip, port, zkh, ch, sigCh)
	return ch
}

func watchLeader(ip string, port string, zkh []string, ch chan bool, sigCh <-chan bool) {
	zkConn, sCh, err := zk.Connect(zkh, time.Second)
	if err != nil {
		log.Fatalf("Error connecting to zookeeper: %s", err)
	}
	defer func() {
		zkConn.Close()
	}()

	// setup required znodes, set watch on candidate node
	createParentZNodes(zkConn)
	zn := createZNode(zkConn, ip, port)
	nodeID := strings.Split(zn, "/")[3]
	wCh, ldr := watchNode(zkConn, nodeID)

	for {
		select {
		// for zk session changes
		case se := <-sCh:
			if se.State == zk.StateUnknown || se.State == zk.StateDisconnected {
				log.Printf("Zookeeper %s, state: %s, server: %s ", se.Type, se.State, se.Server)
			}
		// for potential leader changes
		case we := <-wCh:
			if we.Type == zk.EventNodeDeleted {
				log.Printf("Watched node deleted, resetting watch")
				if wCh, ldr = watchNode(zkConn, nodeID); ldr {
					ch <- true
				} else {
					ch <- false
				}
			}
		case sig := <-sigCh:
			if sig {
				return
			}
		}
	}
}

// sets a watch on the candidate node having seq number prior to supplied node id
func watchNode(conn *zk.Conn, nodeID string) (<-chan zk.Event, bool) {
	wchID := nodeID

	nds, _, ch, err := conn.ChildrenW(zkNodes)
	if err != nil {
		log.Fatal(err)
	}

	sort.Strings(nds)

	// set the watch id to the preceding node to prevent herd effect,
	// see "http://zookeeper.apache.org/doc/trunk/recipes.html#sc_leaderElection"
	for _, id := range nds {
		if id >= nodeID {
			break
		} else if id < wchID {
			wchID = id
		}
	}

	// not leader, set watch for potential leader changes
	if wchID != nodeID {
		data, _, ch, err := conn.GetW(zkNodes + "/" + wchID)
		if err != nil {
			log.Fatalf("Error setting watch on candidate node: %s, %s", wchID, data)
		}
		log.Printf("Entering follower mode, watching candidate node: %s for changes", wchID)
		return ch, false
	}

	// leader, return child watch channel
	log.Printf("Entering leader mode")
	return ch, true
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
