package server

import (
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func Run(ip int, host string) {
	log.Print("Initializing lamport...")
	log.Printf("Running lamport on host: %s, port: %d", host, ip)
}
