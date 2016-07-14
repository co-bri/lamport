package node

import (
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Distributed-Computing-Denver/lamport/config"
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

type zkLdrWatch struct {
	conf config.Config
	conn *zk.Conn
	ch   <-chan zk.Event
}

func newZkLdrWatch(conf config.Config) zkLdrWatch {
	zkConn, sCh, err := zk.Connect(conf.Zookeeper, 250*time.Millisecond)
	if err != nil {
		log.Fatalf("Error connecting to zookeeper: %s", err)
	}
	z := zkLdrWatch{
		conf: conf,
		conn: zkConn,
		ch:   sCh,
	}
	return z
}

func (z zkLdrWatch) leaderWatch(sigCh chan bool) chan bool {
	ch := make(chan bool)
	go z.strtLeaderWatch(ch, sigCh)
	return ch
}

func (z zkLdrWatch) strtLeaderWatch(ch chan bool, sigCh chan bool) {
	// setup required znodes, set watch on candidate node
	z.createParentZNodes()
	zn := z.createZNode()
	nodeID := strings.Split(zn, "/")[3]
	wCh, ldr := z.watchNode(nodeID)

	if ldr {
		log.Print("No leader, entering leader mode")
	}

	for {
		select {
		// for zk session changes
		case se := <-z.ch:
			if se.State == zk.StateUnknown || se.State == zk.StateDisconnected {
				log.Printf("Zookeeper %s, state: %s, server: %s ", se.Type, se.State, se.Server)
			}
		// for potential leader changes
		case we := <-wCh:
			if we.Type == zk.EventNodeDeleted {
				log.Printf("Watched node deleted, resetting watch")
				if wCh, ldr = z.watchNode(nodeID); ldr {
					ch <- true
				} else {
					ch <- false
				}
			}
		case sig := <-sigCh:
			if sig {
				z.closeConn()
				sigCh <- true
				return
			}
		}
	}
}

// sets a watch on the candidate node having seq number prior to supplied node id
func (z zkLdrWatch) watchNode(nodeID string) (<-chan zk.Event, bool) {
	wchID := nodeID

	nds, _, ch, err := z.conn.ChildrenW(zkNodes)
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
		data, _, ch, err := z.conn.GetW(zkNodes + "/" + wchID)
		if err != nil {
			log.Fatalf("Error setting watch on candidate node: %s, %s", wchID, data)
		}
		log.Printf("Entering follower mode, watching candidate node: %s for changes", wchID)
		return ch, false
	}

	// leader, return child watch channel
	return ch, true
}

// creates a candidate znode for this lamport instance
func (z zkLdrWatch) createZNode() (path string) {
	data := []uint8(z.conf.Host + ":" + z.conf.Port)
	path, err := z.conn.Create(zkNodes+"/"+zkPrfx, data, zk.FlagEphemeral|zk.FlagSequence, acl)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created candidate znode %s", path)
	return path
}

// creates the required parent znodes if not present
func (z zkLdrWatch) createParentZNodes() {
	exists, _, err := z.conn.Exists(zkRoot)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Printf("Root znode not found, creating %s in zookeeper", zkRoot)
		if _, err := z.conn.Create(zkRoot, nil, 0, acl); err != nil {
			log.Fatal(err)
		}
	}

	exists, _, err = z.conn.Exists(zkNodes)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		log.Printf("Leader election parent node not found, Creating %s in zookeeper", zkNodes)
		if _, err := z.conn.Create(zkNodes, nil, 0, acl); err != nil {
			log.Fatal(err)
		}
	}
}

func (z zkLdrWatch) closeConn() {
	if z.conn == nil {
		return
	}

	if s := z.conn.State(); s == zk.StateHasSession || s == zk.StateConnected {
		z.conn.Close()
	}
}
