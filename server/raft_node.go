package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

const (
	maxTCPPoolSize      = 3
	raftDbName          = "raft.db"
	retainSnapshotCount = 2
	tcpTimeout          = 10 * time.Second
)

func init() {
	log.SetOutput(os.Stdout)
}

type fsm RaftNode

type fsmSnapshot struct {
	leaderAddr string
}

type command struct {
	Val string
}

// RaftNode wraps the raft struct, along with fields
// for storage and communication.
type RaftNode struct {
	raft        *raft.Raft
	raftAddr    string
	raftDir     string
	lamportAddr string

	leaderAddr string
	mu         sync.Mutex
}

// NewRaftNode creates a new Raft node that can be reached and communicates on
// a given host and port. Also takes in a raft directory to use for various raft
// storage functions.
func NewRaftNode(host, port, lamportPort, raftDir string) (*RaftNode, error) {
	r := &RaftNode{}

	err := r.init(host, port, lamportPort, raftDir)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *RaftNode) init(host, port, lamportPort, raftDir string) error {
	r.raftDir = raftDir
	log.Printf("Setting raft directory to " + r.raftDir)

	r.raftAddr = net.JoinHostPort(host, port)
	log.Printf("Raft protocol listening on " + r.raftAddr)

	r.lamportAddr = net.JoinHostPort(host, lamportPort)

	config := raft.DefaultConfig()
	config.EnableSingleNode = true

	addr, err := net.ResolveTCPAddr("tcp", r.raftAddr)
	if err != nil {
		return fmt.Errorf("Error resolving tcp address: %s", err)
	}

	transport, err :=
		raft.NewTCPTransport(r.raftAddr, addr, maxTCPPoolSize, tcpTimeout, os.Stderr)
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
		raft.NewRaft(config, (*fsm)(r), boltStore, boltStore, snapshots, peerStore, transport)
	if err != nil {
		return fmt.Errorf("Error creating raft struct: %s", err)
	}
	r.raft = raft

	go r.updateFsm()

	return nil
}

func (r *RaftNode) updateFsm() {
	ch := r.raft.LeaderCh()
	for {
		isLeader := <-ch
		if isLeader {
			r.Set(r.lamportAddr)
		}
	}
}

// Applies a command to the state machine
func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command

	if err := json.Unmarshal(l.Data, &c); err != nil {
		log.Printf("Couldn't unmarshal command - %s", err)
	}

	return f.apply(c.Val)
}

// Restores the state machine from a previous snapshot
func (f *fsm) Restore(rc io.ReadCloser) error {
	var str string
	if err := json.NewDecoder(rc).Decode(&str); err != nil {
		return err
	}
	f.leaderAddr = str

	return nil
}

// Takes a snapshot of the state machine
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return &fsmSnapshot{leaderAddr: f.leaderAddr}, nil
}

func (f *fsm) apply(value string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()
	log.Printf("Setting leader address to %s", value)
	f.leaderAddr = value
	return nil
}

// Saves a snapshot of the state machine
func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.leaderAddr)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		if err := sink.Close(); err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		sink.Cancel()
		return err
	}

	return nil
}

func (f *fsmSnapshot) Release() {}

// Join attempts to add another Raft node to the cluster by passing in its
// address. Will fail if not run on the leader.
func (r *RaftNode) Join(addr string) error {
	log.Printf("Attempting to add %s to the cluster...", addr)

	if future := r.raft.AddPeer(addr); future.Error() != nil {
		return fmt.Errorf("Error adding %s to the cluster: %s", addr, future.Error())
	}

	log.Printf("Successfully added %s to the cluster!", addr)

	err := r.Set(r.LeaderAddr())

	return err
}

// Leader returns the raft address of the leader of the cluster.
func (r *RaftNode) Leader() string {
	return r.raft.Leader()
}

// LamportAddr gets the lamport address of this node. If this node is the
// leader, it will be equal to LeaderAddr()
func (r *RaftNode) LamportAddr() string {
	return r.lamportAddr
}

// LeaderAddr returns the lamport address of the leader node.
func (r *RaftNode) LeaderAddr() string {
	return r.leaderAddr
}

// RaftAddr returns this node's raft address.
func (r *RaftNode) RaftAddr() string {
	return r.raftAddr
}

// Shutdown shuts down the raft cluster. A blocking operation.
func (r *RaftNode) Shutdown() error {
	future := r.raft.Shutdown()
	if future.Error() != nil {
		return future.Error()
	}
	return nil
}

// Set sets the leader address of the Raft node and applies
// it to all nodes in the cluster. Can only be run
// on the leader.
func (r *RaftNode) Set(leaderAddr string) error {
	if r.State() != "Leader" {
		return fmt.Errorf("Can't set value on non-leader node")
	}

	c := &command{Val: leaderAddr}
	log.Printf(c.Val)

	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	future := r.raft.Apply(bytes, tcpTimeout)
	if err := future.Error(); err != nil {
		return err
	}

	return nil
}

// State returns the state of the Raft node, which is one of Candidate,
// Leader, Follower, or Shutdown.
func (r *RaftNode) State() string {
	return r.raft.State().String()
}

// Stats returns the map of various internal raft stats.
func (r *RaftNode) Stats() map[string]string {
	return r.raft.Stats()
}
