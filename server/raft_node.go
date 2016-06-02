package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

const (
	maxTcpPoolSize      = 3
	raftDbName          = "raft.db"
	retainSnapshotCount = 2
	tcpTimeout          = 10 * time.Second
)

func init() {
	log.SetOutput(os.Stdout)
}

// A struct that wraps the raft struct, along with fields
// for storage and communication.
type RaftNode struct {
	raft     *raft.Raft
	RaftAddr string
	raftDir  string
}

// Create a new Raft node that can be reached and communicates on
// a given host and port. Also takes in a raft directory to
// use for various raft storage functions.
func NewRaftNode(host, port, raftDir string) (*RaftNode, error) {
	r := &RaftNode{}

	err := r.init(host, port, raftDir)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RaftNode) init(host, port, raftDir string) error {
	r.raftDir = raftDir
	log.Printf("Setting raft directory to " + r.raftDir)

	r.RaftAddr = net.JoinHostPort(host, port)
	log.Printf("Raft protocol listening on " + r.RaftAddr)

	config := raft.DefaultConfig()
	config.EnableSingleNode = true

	addr, err := net.ResolveTCPAddr("tcp", r.RaftAddr)
	if err != nil {
		return fmt.Errorf("Error resolving tcp address: %s", err)
	}

	transport, err :=
		raft.NewTCPTransport(r.RaftAddr, addr, maxTcpPoolSize, tcpTimeout, os.Stderr)
	if err != nil {
		return fmt.Errorf("Error creating raft TCP transport: %s", err)
	}

	peerStore := raft.NewJSONPeers(r.raftDir, transport)

	snapshots, err := raft.NewFileSnapshotStore(r.raftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("Error creating file snap shot store: %s", err)
	}

	boltStore, err := raftboltdb.NewBoltStore(filepath.Join(r.raftDir, raftDbName))
	if err != nil {
		return fmt.Errorf("Error creating bolt db store: %s", err)
	}

	raft, err :=
		raft.NewRaft(config, nil, boltStore, boltStore, snapshots, peerStore, transport)
	if err != nil {
		return fmt.Errorf("Error creating raft struct: %s", err)
	}
	r.raft = raft

	return nil
}

// Attempts to add another Raft node to the cluster by passing in its
// address. Will fail if not run on the leader.
func (r *RaftNode) Join(addr string) error {
	log.Printf("Attempting to add %s to the cluster...", addr)

	if future := r.raft.AddPeer(addr); future.Error() != nil {
		return fmt.Errorf("Error adding %s to the cluster: %s", addr, future.Error())
	}

	log.Printf("Successfully added %s to the cluster!", addr)

	return nil
}

// Shuts down the raft cluster. A blocking operation.
func (r *RaftNode) Shutdown() error {
	future := r.raft.Shutdown()
	if future.Error() != nil {
		return future.Error()
	} else {
		return nil
	}
}

// Returns the state of the Raft node, which is one of Candidate,
// Leader, Follower, or Shutdown.
func (r *RaftNode) State() string {
	return r.raft.State().String()
}

// Returns the map of various internal raft stats.
func (r *RaftNode) Stats() map[string]string {
	return r.raft.Stats()
}
