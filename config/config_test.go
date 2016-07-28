package config

import "testing"

const testFile = "test.toml"

func TestReadConfig(t *testing.T) {
	c, err := ReadConfig("test.toml")

	if err != nil {
		t.Fatalf("Encountered error reading config file %s", testFile)
	}

	if c.Host != "127.0.0.1" {
		t.Fatalf("Expected 127.0.0.1 for 'IP', found %s", c.Host)
	}
	if c.Port != "5936" {
		t.Fatalf("Expected 5936 for 'Port', found %s", c.Port)
	}
}

func TestReadConfigError(t *testing.T) {
	_, err := ReadConfig("foo.toml")
	if err == nil {
		t.Fatal("Expected error when reading non-existent config")
	}
}
