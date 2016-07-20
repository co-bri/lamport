package node

import (
	"os"
	"strconv"
	"testing"
	"time"
)

const (
	localhost = "127.0.0.1"
	port      = "4567"
	raftDir   = ".testRaftDir"
)

// Create a single node and verify that it becomes leader
func TestSingleNodeCluster(t *testing.T) {
	raftNode, err := NewRaftNode(localhost, port, port, raftDir)
	if err != nil {
		t.Error(err)
	}

	// give cluster time to elect a leader
	time.Sleep(3 * time.Second)

	if state := raftNode.State(); state != "Leader" {
		t.Errorf("Expected state - Leader, Actual state - %s", state)
	}

	if raftNode.Leader() != localhost+":"+port {
		t.Errorf("Expected leader - %s, Actual leader - %s", localhost+":"+port, raftNode.Leader())
	}

	shutdownRaftNode(raftNode, t)
	deleteDir(raftDir, t)
}

// Create three nodes and join them together to form a cluster.
// Each node should have 2 peers and the cluster should have a
// leader.
func TestThreeNodeCluster(t *testing.T) {
	portNum, err := strconv.Atoi(port)
	if err != nil {
		t.Error(err)
	}

	raftNode1, err := NewRaftNode(localhost, strconv.Itoa(portNum+1), port, raftDir+"1")
	if err != nil {
		t.Error(err)
	}

	raftNode2, err := NewRaftNode(localhost, strconv.Itoa(portNum+2), port, raftDir+"2")
	if err != nil {
		t.Error(err)
	}

	// wait for a leader
	time.Sleep(3 * time.Second)

	if raftNode2.State() == "Leader" {
		err = raftNode2.Join(raftNode1.RaftAddr())
		if err != nil {
			t.Error(err)
		}
	} else if raftNode1.State() == "Leader" {
		err = raftNode1.Join(raftNode2.RaftAddr())
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Errorf("No leader in cluster 5 seconds after node 1 + 2 creation. Node 1 state: %s, Node 2 state: %s", raftNode1.State(), raftNode2.State())
	}

	// wait for a leader
	time.Sleep(3 * time.Second)

	raftNode3, err := NewRaftNode(localhost, strconv.Itoa(portNum+3), port, raftDir+"3")
	if err != nil {
		t.Error(err)
	}

	if raftNode1.State() == "Leader" {
		err = raftNode1.Join(raftNode3.RaftAddr())
		if err != nil {
			t.Error(err)
		}
	} else if raftNode2.State() == "Leader" {
		err = raftNode2.Join(raftNode3.RaftAddr())
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Error("No leader in cluster 5 seconds after node 2 joined node 1.")
	}

	// wait for a leader again
	time.Sleep(3 * time.Second)

	clusterHasLeader := raftNode1.State() == "Leader" || raftNode2.State() == "Leader" ||
		raftNode3.State() == "Leader"

	if !clusterHasLeader {
		t.Error("No leader in cluster 5 seconds after node 3 joined non leader node")
	}

	nodesHaveThreePeers := getPeersCount(raftNode1.Stats()) == 2 &&
		getPeersCount(raftNode2.Stats()) == 2 &&
		getPeersCount(raftNode3.Stats()) == 2

	if !nodesHaveThreePeers {
		t.Error("All nodes do not have three peers")
	}

	shutdownRaftNode(raftNode1, t)
	shutdownRaftNode(raftNode2, t)
	shutdownRaftNode(raftNode3, t)
	deleteDir(raftDir+"1", t)
	deleteDir(raftDir+"2", t)
	deleteDir(raftDir+"3", t)
}

func deleteDir(dir string, t *testing.T) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Logf("Couldn't delete directory %s - %s", dir, err)
	}
}

func shutdownRaftNode(r *RaftNode, t *testing.T) {
	t.Logf("Shutting down raft node %v", r.raftAddr)

	err := r.Shutdown()
	if err != nil {
		t.Errorf("Error shutting down raft node: %s", err)
	}
}

func getPeersCount(stats map[string]string) int {
	peersCount, _ := strconv.Atoi(stats["num_peers"])
	return peersCount
}
