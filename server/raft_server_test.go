package server_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/Distributed-Computing-Denver/lamport/server"
)

func TestThreeServerCluster(t *testing.T) {
	localhost := "127.0.0.1"
	lamportPort1 := "3456"
	lamportPort2 := "3457"
	lamportPort3 := "3458"
	raftPort1 := "2345"
	raftPort2 := "2346"
	raftPort3 := "2347"
	raftDir := ".testRaftDir"

	// Create 3 raft nodes
	raftNode1, err := server.NewRaftNode(localhost, raftPort1, lamportPort1, raftDir+"1")
	if err != nil {
		t.Error(err)
	}

	raftNode2, err := server.NewRaftNode(localhost, raftPort2, lamportPort2, raftDir+"2")
	if err != nil {
		t.Error(err)
	}

	raftNode3, err := server.NewRaftNode(localhost, raftPort3, lamportPort3, raftDir+"3")
	if err != nil {
		t.Error(err)
	}

	// Spin up 2 servers and tell server 2 to join server 1 on startup
	go server.RunRaftServer(localhost, lamportPort1, raftNode1, "")

	time.Sleep(3 * time.Second)

	go server.RunRaftServer(localhost, lamportPort2, raftNode2, net.JoinHostPort(localhost, lamportPort1))

	// Wait for nodes to communicate
	time.Sleep(3 * time.Second)

	node1Peers := getPeersCount(raftNode1.Stats())
	node2Peers := getPeersCount(raftNode2.Stats())

	if node1Peers != 1 && node2Peers != 1 {
		t.Errorf("Both raft nodes should have one peer, node 1 has %d and node 2 has %d", node1Peers, node2Peers)
	}

	// Create a third server
	go server.RunRaftServer(localhost, lamportPort3, raftNode3, "")

	time.Sleep(3 * time.Second)

	// Tell server 3 to join the non-leader server
	var nonLeaderNode string
	if raftNode1.State() == "Leader" {
		nonLeaderNode = net.JoinHostPort(localhost, lamportPort2)
	} else if raftNode2.State() == "Leader" {
		nonLeaderNode = net.JoinHostPort(localhost, lamportPort1)
	} else {
		t.Errorf("node 2 stats: %s", raftNode2.Stats())
		t.Errorf("node 1 stats: %s", raftNode1.Stats())
		t.Fatalf("No leader in two-node cluster")
		cleanup(t)
	}

	conn, err := net.Dial("tcp", nonLeaderNode)
	if err != nil {
		t.Errorf("Error connecting to %s: %s", nonLeaderNode, err)
	}
	fmt.Fprintf(conn, net.JoinHostPort(localhost, raftPort3))
	conn.Close()

	// Wait for nodes to communicate
	time.Sleep(3 * time.Second)

	// Verify cluster has leader and has a size of 3
	clusterHasLeader := raftNode1.State() == "Leader" || raftNode2.State() == "Leader" ||
		raftNode3.State() == "Leader"

	if !clusterHasLeader {
		t.Error("No leader in cluster 5 seconds after node 3 joined non leader node")
	}

	nodesHaveTwoPeers := getPeersCount(raftNode1.Stats()) == 2 &&
		getPeersCount(raftNode2.Stats()) == 2 &&
		getPeersCount(raftNode3.Stats()) == 2

	if !nodesHaveTwoPeers {
		t.Error("All nodes do not have three peers")
	}

	shutdownRaftNode(raftNode1, t)
	shutdownRaftNode(raftNode2, t)
	shutdownRaftNode(raftNode3, t)
	cleanup(t)
}

func cleanup(t *testing.T) {
	deleteDir(raftDir+"1", t)
	deleteDir(raftDir+"2", t)
	deleteDir(raftDir+"3", t)
}
