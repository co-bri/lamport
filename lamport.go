package main

import (
	"fmt"
	"os"

	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/Distributed-Computing-Denver/lamport/node"
	"github.com/urfave/cli"
)

const (
	name       = "lamport"
	usage      = "An academic exercise in distributed systems"
	version    = "0.0.1"
	rName      = "run"
	rUsage     = "run lamport"
	rFlagName  = "config, c"
	rFlagValue = "lamport.toml"
	rFlagUsage = "lamport configuration `FILE`"
)

func main() {
	app := getApp()
	app.Run(os.Args)
}

func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = name
	app.Usage = usage
	app.Version = version

	app.Commands = []cli.Command{
		{
			Name:    rName,
			Aliases: []string{"r"},
			Usage:   rUsage,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  rFlagName,
					Value: rFlagValue,
					Usage: rFlagUsage,
				},
			},
			Action: func(c *cli.Context) error {
				cf := c.String("config")
				config, err := config.ReadConfig(cf)
				if err != nil {
					panic(fmt.Errorf("Error reading config file: %s", err))
				}
				node.Start(node.LamportRunner(config))
				return nil
			},
		},
	}
	return app
}
