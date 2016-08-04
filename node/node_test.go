package node

import (
	"os"
	"testing"
	"time"

	"github.com/Distributed-Computing-Denver/lamport/config"
)

const (
	host = "127.0.0.1"
	port = "5936"
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

func TestRun(t *testing.T) {
	c := getConfig()
	r := New(c)
	n, ok := r.(node)

	if !ok {
		t.Fatal("Unable to get `node` from `Runner`")
	}

	if n.conf.Host != host {
		t.Fatalf("Expected node 'Host' to be %s, but found %s", host, n.conf.Host)
	}

	if n.conf.Port != port {
		t.Fatalf("Expected node 'Port' to be %s, but found %s", port, n.conf.Port)
	}

	ch := make(chan bool)
	go n.Run(ch)
	time.Sleep(1 * time.Second)
	ch <- true
	<-ch
}

func getConfig() config.Config {
	return config.Config{
		Host: host,
		Port: port,
	}
}
