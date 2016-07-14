package node

import (
	"os"
	"testing"
	"time"

	"github.com/Distributed-Computing-Denver/lamport/config"
)

func TestStart(t *testing.T) {
	ch := make(chan bool)
	rnr := func(sigCh chan bool) {
		<-sigCh
		ch <- true
	}
	go Start(rnr)
	time.Sleep(1 * time.Second)
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
	p.Signal(os.Interrupt)
	<-ch
}

func TestNewNode(t *testing.T) {
	c := config.Config{
		Zookeeper: []string{"127.0.0.1"},
	}

	n := newNode(c)
	if n.lw == nil {
		t.Fatal("Expected node.lw to be non-nil")
	}

	n = newNode(config.Config{})
	if n.lw != nil {
		t.Fatal("Expected node.lw to be nil")
	}
}

type ldrWtch struct{}

func (l ldrWtch) leaderWatch(ch chan bool) chan bool {
	go func() {
		<-ch
		ch <- true
	}()
	return make(chan bool)
}

func TestRun(t *testing.T) {
	ch := make(chan bool)
	n := node{}
	n.lw = ldrWtch{}

	go n.run(ch)
	time.Sleep(1 * time.Second)
	ch <- true
	<-ch
}
