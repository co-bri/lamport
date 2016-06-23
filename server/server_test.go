// Unit tests for the 'server' package
package server

import (
	"os/exec"
	"testing"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

const (
	zkSrvr = "zkServer"
	zkStrt = "start"
	zkStop = "stop"
)

func TestRunSingleNode(t *testing.T) {
	conn, zkExec, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	defer func() {
		conn.Close()
		stopZk(zkExec)
	}()

	ch := make(chan bool)
	go Run("127.0.0.1", "5936", ch)
	time.Sleep(10 * time.Second)

	nds, _, _, err := conn.ChildrenW(zkNodes)
	if err != nil {
		t.Fatalf("Error verifying candidate znode: %s", err)
	}

	if len(nds) != 1 {
		t.Fatalf("Expected 1 candidate node, but found %d", len(nds))
	}

	ch <- true
	time.Sleep(5 * time.Second)

	nds, _, _, err = conn.ChildrenW(zkNodes)
	if err != nil {
		t.Fatalf("Error verifying candidate znode: %s", err)
	}

	if len(nds) != 0 {
		t.Fatalf("Expected 0 candidate node, but found %d", len(nds))
	}
}

func TestWatchCandidateNode(t *testing.T) {
	conn, zkExec, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	defer func() {
		conn.Close()
		stopZk(zkExec)
	}()
	createParentZNodes(conn)                    // Need to make sure lamport zk nodes are present
	pth := createZNode(conn, "0.0.0.0", "1234") // Need to create candidate node
	ch := watchNode(conn, pth)

	// delete the candidate node to generate an event on the channel
	if err := conn.Delete(pth, 0); err != nil {
		t.Fatalf("Error deleting %s: %s", pth, err)
	}
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
	conn, zkExec, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	defer func() {
		conn.Close()
		stopZk(zkExec)
	}()
	createParentZNodes(conn) // Need to make sure lamport zk nodes are present
	pth := createZNode(conn, "0.0.0.0", "1234")
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
	conn, zkExec, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	defer func() {
		conn.Close()
		stopZk(zkExec)
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

func startZk() (*zk.Conn, string, error) {
	zkExec := zkSrvr
	if _, err := exec.LookPath(zkExec); err != nil {
		// also check for zkServer.sh
		zkExec = zkSrvr + ".sh"
		if _, err := exec.LookPath(zkExec); err != nil {
			return nil, zkExec, err
		}
	}
	cmd := exec.Command(zkExec, zkStrt)
	if err := cmd.Start(); err != nil {
		return nil, zkExec, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, zkExec, err
	}
	defer func() {
		if r := recover(); r != nil {
			stopZk(zkExec)
		}
	}()

	// this sucks, but the zk library is not friendly to timeouts
	time.Sleep(25 * time.Second)

	conn, _, err := zk.Connect([]string{"127.0.0.1"}, time.Second)
	if err != nil {
		return nil, zkExec, err
	}

	return conn, zkExec, nil
}

func cleanZk(conn *zk.Conn) error {
	nds, _, _, err := conn.ChildrenW(zkNodes)
	if err != nil {
		return err
	}

	for _, id := range nds {
		if err := conn.Delete(zkNodes+"/"+id, 0); err != nil {
			return err
		}
	}

	if err := conn.Delete(zkNodes, 0); err != nil {
		return err
	}

	if err := conn.Delete(zkRoot, 0); err != nil {
		return err
	}
	return nil
}

func stopZk(zkExec string) error {
	if _, err := exec.LookPath(zkExec); err != nil {
		return err
	}
	cmd := exec.Command(zkExec, zkStop)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
