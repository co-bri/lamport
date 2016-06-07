// Unit tests for the 'server' package
package server

import (
	"os/exec"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	zkSrvr = "zkServer.sh"
	zkStrt = "start"
	zkStop = "stop"
)

func TestWatchCandidateNode(t *testing.T) {
	conn, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	defer func() {
		conn.Close()
		stopZk()
	}()
	createParentZNodes(conn)                             // Need to make sure lamport zk nodes are present
	pth := createCandidateZNode(conn, "0.0.0.0", "1234") // Need to create candidate node
	ch := watchCandidateNode(conn, pth)

	// delete the candidate node to generate an event on the channel
	if err := conn.Delete(pth, 0); err != nil {
		t.Fatalf("Error deleting %s: %s", pth, err)
	}
	// server_test.go:35: {EventNodeChildrenChanged Unknown /lamport/nodes <nil> }
	// receive the delete event
	e := <-ch
	if e.Type != zk.EventNodeChildrenChanged {
		t.Fatalf("Expected EventType of 'EventNodeChildrenChanged', but found %s", e.Type)
	}
	if e.Path != zkNodes {
		t.Fatalf("Expected a Path of %s, but found %s", zkNodes, e.Path)
	}
	if err := cleanZk(conn); err != nil {
		t.Fatalf("Error cleaning up zk nodes %s", err)
	}
}

func TestCreateCandidateZNode(t *testing.T) {
	conn, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	defer func() {
		conn.Close()
		stopZk()
	}()
	createParentZNodes(conn) // Need to make sure lamport zk nodes are present
	pth := createCandidateZNode(conn, "0.0.0.0", "1234")
	if exists, _, err := conn.Exists(pth); err != nil || !exists {
		t.Fatalf("Failed to create candidate node %s", pth)
	}
	if err := conn.Delete(pth, 0); err != nil {
		t.Fatalf("Error deleting zk root node %s: %s", pth, err)
	}
	if err := cleanZk(conn); err != nil {
		t.Fatalf("Error cleaning up zk nodes %s", err)
	}
}

func TestCreateParentZNodes(t *testing.T) {
	conn, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	defer func() {
		conn.Close()
		stopZk()
	}()
	createParentZNodes(conn)
	if exists, _, err := conn.Exists(zkRoot); err != nil || !exists {
		t.Fatalf("Failed to create zk root node %s: %s", zkRoot, err)
	}

	if exists, _, err := conn.Exists(zkNodes); err != nil || !exists {
		t.Fatalf("Failed to create zk leader election node %s: %s", zkNodes, err)
	}

	if err := cleanZk(conn); err != nil {
		t.Fatalf("Error cleaning up zk nodes %s", err)
	}
}

func TestNodeSeq(t *testing.T) {
	s := "_c_4f4bfdb5805df1a619e1c8e8b26d86e6-0000000002"
	seq := getNodeSeq(s)

	// happy path
	if seq != 2 {
		t.Fatalf("Expected seq 2 for node id %s ", s)
	}

	// test for panic with bad node id
	s = "ABCDEFGHIJK"
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic for seq for node id %s", s)
		}
	}()
	getNodeSeq(s)
}

func startZk() (*zk.Conn, error) {
	if _, err := exec.LookPath(zkSrvr); err != nil {
		return nil, err
	}
	cmd := exec.Command(zkSrvr, zkStrt)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func cleanZk(conn *zk.Conn) error {
	if err := conn.Delete(zkNodes, 0); err != nil {
		return err
	}

	if err := conn.Delete(zkRoot, 0); err != nil {
		return err
	}
	return nil
}

func stopZk() error {
	if _, err := exec.LookPath(zkSrvr); err != nil {
		return err
	}
	cmd := exec.Command(zkSrvr, zkStop)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
