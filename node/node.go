package node

import (
	"log"
	"os"
	"os/signal"

	"github.com/Distributed-Computing-Denver/lamport/config"
)

// Runner runs until signalled to stop by sigCh
type Runner interface {
	Run(sigCh chan bool)
}

type node struct {
	conf config.Config
}

// Start starts a new lamport node using the supplied
// Runner
func Start(r Runner) error {
	sigCh := make(chan bool)
	go r.Run(sigCh)

	// handle SIGINT, notify node, wait for confirm to exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("Received SIGINT, terminating lamport")
	sigCh <- true
	<-sigCh

	return nil
}

// New creates a Runner that can be used to
// start a Lamport node
func New(conf config.Config) Runner {
	return node{conf: conf}
}

func (n node) Run(sigCh chan bool) {
	sig := <-sigCh
	if sig {
		sigCh <- true
		return
	}
}
