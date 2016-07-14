package node

import (
	"log"
	"os"
	"os/signal"

	"github.com/Distributed-Computing-Denver/lamport/config"
)

type node struct {
	lw leaderWatcher
}

// Runner runs until signalled to stop by sigCh
type Runner func(sigCh chan bool)

type leaderWatcher interface {
	leaderWatch(ch chan bool) chan bool
}

// LamportRunner creates a Runner from the supplied config
func LamportRunner(conf config.Config) Runner {
	return func(sigCh chan bool) {
		newNode(conf).run(sigCh)
	}
}

// Start starts a new lamport node using the supplied
// Creator and Config
func Start(r Runner) {
	sigCh := make(chan bool)
	go r(sigCh)

	// handle SIGINT, notify node, wait for confirm to exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Print("Received SIGINT, terminating lamport")
	sigCh <- true
	<-sigCh
}

func newNode(conf config.Config) node {
	n := node{}

	// only support zk leader elec
	if len(conf.Zookeeper) > 0 {
		n.lw = newZkLdrWatch(conf)
	}

	return n
}

func (n node) run(sigCh chan bool) {
	ch := make(chan bool)
	ldrCh := n.lw.leaderWatch(ch)

	for {
		select {
		case ldr := <-ldrCh:
			if ldr {
				log.Print("Entering leader mode")
			}
		case q := <-sigCh:
			if q {
				// signal leader watcher
				ch <- true
				<-ch
				sigCh <- true
				return
			}
		}
	}
}
