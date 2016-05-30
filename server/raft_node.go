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

type RaftNode struct {
	raft     *raft.Raft
	raftAddr string
	raftDir  string
}

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

	r.raftAddr = net.JoinHostPort(host, port)
	log.Printf("Raft protocol listening on " + r.raftAddr)

	config := raft.DefaultConfig()
	config.EnableSingleNode = true
	config.DisableBootstrapAfterElect = false

	addr, err := net.ResolveTCPAddr("tcp", r.raftAddr)
	if err != nil {
		return fmt.Errorf("Error resolving tcp address: %s", err)
	}

	transport, err :=
		raft.NewTCPTransport(r.raftAddr, addr, maxTcpPoolSize, tcpTimeout, os.Stderr)
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

func (r *RaftNode) State() string {
	return r.raft.State().String()
}
