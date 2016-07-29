package main

import (
	"testing"

	"github.com/urfave/cli"
)

const (
	name       = "lamport"
	usage      = "An academic exercise in building a distributed system"
	version    = "0.0.1"
	rName      = "run"
	rUsage     = "run a lamport node"
	rFlagName  = "config, c"
	rFlagValue = "lamport.toml"
	rFlagUsage = "lamport configuration `FILE`"
)

func TestGetApp(t *testing.T) {
	app := getApp()

	if app.Name != name {
		t.Fatalf("Expected %s for app Name, but found %s", name, app.Name)
	}

	if app.Usage != usage {
		t.Fatalf("Expected %s for app Usage, but found %s", usage, app.Usage)
	}

	if app.Version != version {
		t.Fatalf("Expected %s for app Version, but found %s", version, app.Version)
	}

	if len(app.Commands) != 1 {
		t.Fatalf("Expected 1 subcommand for app but found %d", len(app.Commands))
	}

	cmd := app.Commands[0]

	if cmd.Name != rName {
		t.Fatalf("Expected %s for subcommand Name, but found %s", rName, cmd.Name)
	}

	if cmd.Usage != rUsage {
		t.Fatalf("Expected %s for subcommand Name, but found %s", rUsage, cmd.Usage)
	}

	if len(cmd.Flags) != 1 {
		t.Fatalf("Expected single flag for subcommand %s, but found %d", cmd.Name, len(cmd.Flags))
	}

	strFlag := cmd.Flags[0].(cli.StringFlag)

	if strFlag.Name != rFlagName {
		t.Fatalf("Expected %s for subcommand %s flag name, but found %s", rFlagName, cmd.Name, strFlag.Name)
	}

	if strFlag.Value != rFlagValue {
		t.Fatalf("Expected %s for subcommand %s flag value, but found %s", rFlagValue, cmd.Name, strFlag.Value)
	}

	if strFlag.Usage != rFlagUsage {
		t.Fatalf("Expected %s for subcommand %s flag usage, but found %s", rFlagUsage, cmd.Name, strFlag.Usage)
	}
}
