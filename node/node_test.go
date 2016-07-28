package node

import (
	"os"
	"testing"
	"time"

	"github.com/Distributed-Computing-Denver/lamport/config"
)

type testnode struct {
	ch chan bool
}

func (tn testnode) Run(sigCh chan bool) {
	<-sigCh
	sigCh <- true
	tn.ch <- true
}

func TestStart(t *testing.T) {
	ch := make(chan bool)
	tn := testnode{ch: ch}
	go Start(tn)
	time.Sleep(1 * time.Second)
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
	p.Signal(os.Interrupt)
	<-ch
}

func TestNewNode(t *testing.T) {
	h, p := "127.0.0.1", "5936"

	c := config.Config{
		Host: h,
		Port: p,
	}

	r := New(c)
	n, ok := r.(node)

	if !ok {
		t.Fatal("Unable to get `node` from `Runner`")
	}

	if n.conf.Host != h {
		t.Fatalf("Expected node 'Host' to be %s, but found %s", h, n.conf.Host)
	}

	if n.conf.Port != p {
		t.Fatalf("Expected node 'Port' to be %s, but found %s", p, n.conf.Port)
	}
}
