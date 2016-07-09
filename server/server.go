package server

import (
	"log"

	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/Distributed-Computing-Denver/lamport/zk"
)

// Run starts a Lamport node using the supplied configuration
func Run(config config.Config, sigCh chan bool) {
	ch := make(chan bool)
	ldrCh := electLeader(config, ch)

	for {
		select {
		case ldr := <-ldrCh:
			if ldr {
				log.Print("Node entering leader mode")
			}
		case q := <-sigCh:
			if q {
				ch <- true
				<-ch
				return
			}
		}
	}
}

func electLeader(config config.Config, sigCh chan bool) chan bool {
	return zk.LeaderWatch(config.Host, config.Port, config.Zookeeper, sigCh)
}
