package main

import (
	"fmt"
	"os"

	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/Distributed-Computing-Denver/lamport/node"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "lamport"
	app.Usage = "An academic exercise in distributed systems"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run lamport",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Value: "lamport.toml",
					Usage: "lamport configuration `FILE`",
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

	app.Run(os.Args)
}
