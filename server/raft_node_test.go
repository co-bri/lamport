package server_test

import (
	"testing"
	"time"

	"github.com/Distributed-Computing-Denver/lamport/server"
)

const (
	localhost = "127.0.0.1"
	port      = "4567"
	raftDir   = ".testRaftDir"
)

func TestSingleNodeCluster(t *testing.T) {
	raftNode, err := server.NewRaftNode(port, localhost, raftDir)
	if err != nil {
		t.Error(err)
	}

	// give cluster time to elect a leader
	time.Sleep(5 * time.Second)

	if state := raftNode.State(); state != "Leader" {
		t.Errorf("Expected state - Leader, Actual state - %s", state)
	}
}
