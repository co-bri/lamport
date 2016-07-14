package node

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	zkSrvr = "zkServer"
	zkStrt = "start"
	zkStop = "stop"
)

func TestWatchLeaderSingle(t *testing.T) {
	zkExec, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	z := newZkLdrWatch(getConf())
	defer func() {
		z.conn.Close()
		stopZk(zkExec)
	}()
	sigCh := make(chan bool)
	z.leaderWatch(sigCh)
	time.Sleep(10 * time.Second)

	nds, _, _, err := z.conn.ChildrenW(zkNodes)
	if err != nil {
		t.Fatalf("Error verifying candidate znode: %s", err)
	}

	if len(nds) != 1 {
		t.Fatalf("Expected 1 candidate node, but found %d", len(nds))
	}

	sigCh <- true
	<-sigCh

	z = newZkLdrWatch(getConf())
	nds, _, _, err = z.conn.ChildrenW(zkNodes)
	t.Log(nds)
	if err != nil {
		t.Fatalf("Error verifying candidate znode: %s", err)
	}

	if len(nds) != 0 {
		t.Fatalf("Expected 0 candidate node, but found %d", len(nds))
	}
}

func TestWatchCandidateNode(t *testing.T) {
	zkExec, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	z := newZkLdrWatch(getConf())
	defer func() {
		z.conn.Close()
		stopZk(zkExec)
	}()
	z.createParentZNodes() // Need to make sure lamport zk nodes are present
	pth := z.createZNode() // Need to create candidate node
	ch, ldr := z.watchNode(pth)

	if !ldr {
		t.Fatal("Expected single node to be leader")
	}

	// delete the candidate node to generate an event on the channel
	if err := z.conn.Delete(pth, 0); err != nil {
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
	if err := cleanZk(z.conn); err != nil {
		t.Fatalf("Error cleaning up zk nodes %s", err)
	}
}

func TestCreateCandidateZNode(t *testing.T) {
	zkExec, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	z := newZkLdrWatch(getConf())
	defer func() {
		z.conn.Close()
		stopZk(zkExec)
	}()
	z.createParentZNodes() // Need to make sure lamport zk nodes are present
	pth := z.createZNode()
	if exists, _, err := z.conn.Exists(pth); err != nil || !exists {
		t.Fatalf("Failed to create candidate node %s", pth)
	}
	if err := z.conn.Delete(pth, 0); err != nil {
		t.Fatalf("Error deleting zk root node %s: %s", pth, err)
	}
	if err := cleanZk(z.conn); err != nil {
		t.Fatalf("Error cleaning up zk nodes %s", err)
	}
}

func TestCreateParentZNodes(t *testing.T) {
	zkExec, err := startZk()
	if err != nil {
		t.Fatalf("Error starting zookeeper %s", err)
	}
	z := newZkLdrWatch(getConf())
	defer func() {
		z.conn.Close()
		stopZk(zkExec)
	}()
	z.createParentZNodes()
	if exists, _, err := z.conn.Exists(zkRoot); err != nil || !exists {
		t.Fatalf("Failed to create zk root node %s: %s", zkRoot, err)
	}

	if exists, _, err := z.conn.Exists(zkNodes); err != nil || !exists {
		t.Fatalf("Failed to create zk leader election node %s: %s", zkNodes, err)
	}

	if err := cleanZk(z.conn); err != nil {
		t.Fatalf("Error cleaning up zk nodes %s", err)
	}
}

func getConf() config.Config {
	return config.Config{Host: "127.0.0.1", Port: "5936", Zookeeper: []string{"127.0.0.1"}}
}

func startZk() (string, error) {
	zkExec := zkSrvr
	if _, err := exec.LookPath(zkExec); err != nil {
		// also check for zkServer.sh
		zkExec = zkSrvr + ".sh"
		if _, err := exec.LookPath(zkExec); err != nil {
			return zkExec, err
		}
	}
	cmd := exec.Command(zkExec, zkStrt)
	if err := cmd.Start(); err != nil {
		return zkExec, err
	}
	if err := cmd.Wait(); err != nil {
		return zkExec, err
	}
	defer func() {
		if r := recover(); r != nil {
			stopZk(zkExec)
		}
	}()

	// this sucks, would rather use channels to coordinate
	time.Sleep(25 * time.Second)

	return zkExec, nil
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

	if err = conn.Delete(zkNodes, 0); err != nil {
		return err
	}

	if err = conn.Delete(zkRoot, 0); err != nil {
		return err
	}

	if err = os.Remove("zookeeper.out"); err != nil {
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
