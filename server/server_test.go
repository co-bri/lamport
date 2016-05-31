// Unit tests for the 'server' package
package server

import (
	"net"
	"testing"
)

const (
	ip   = "127.0.0.1"
	port = "5936"
)

func TestListen(t *testing.T) {
	ch := make(chan bool)
	go listen(ip, port, ch)

	select {
	case <-ch:
		_, err := net.Dial("tcp", ip+":"+port)
		if err != nil {
			t.Error(err)
		}
	}
}
