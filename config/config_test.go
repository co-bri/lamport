package config

import "testing"

const testFile = "test.toml"

func TestReadConfig(t *testing.T) {
	c, err := ReadConfig("test.toml")

	if err != nil {
		t.Fatalf("Encountered error reading config file %s", testFile)
	}

	if c.Bootstrap != "127.0.0.1:5936" {
		t.Fatalf("Expected 127.0.0.1:5936 for 'Bootstrap', found %s", c.Bootstrap)
	}
	if len(c.Zookeeper) != 3 {
		t.Fatalf("Expected 3 items in 'Zookeeper', found %d", len(c.Zookeeper))
	}
	if c.Zookeeper[0] != "127.0.0.1:2190" {
		t.Fatalf("Expected 127.0.0.1:2190 for 'Zookeeper[0]', found %s", c.Zookeeper[0])
	}
	if c.Zookeeper[1] != "127.0.0.2:2190" {
		t.Fatalf("Expected 127.0.0.2:2190 for 'Zookeeper[1]', found %s", c.Zookeeper[1])
	}
	if c.Zookeeper[2] != "127.0.0.3:2190" {
		t.Fatalf("Expected 127.0.0.3:2190 for 'Zookeeper[2]', found %s", c.Zookeeper[2])
	}
	if c.Host != "127.0.0.1" {
		t.Fatalf("Expected 127.0.0.1 for 'IP', found %s", c.Host)
	}
	if c.Port != "5936" {
		t.Fatalf("Expected 5936 for 'Port', found %s", c.Port)
	}
	if c.RaftDir != ".raft" {
		t.Fatalf("Expected .raft for 'RaftDir', found %s", c.RaftDir)
	}
	if c.RaftPort != "8500" {
		t.Fatalf("Expected 8500 for 'Bootstrap', found %s", c.RaftPort)
	}
}
