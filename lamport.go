package main

import (
	"fmt"
	"os"

	"github.com/Distributed-Computing-Denver/lamport/config"
	"github.com/Distributed-Computing-Denver/lamport/node"
	"github.com/urfave/cli"
)

func main() {
	app := getApp()
	app.Run(os.Args)
}

func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = "lamport"
	app.Usage = "An academic exercise in building a distributed system"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "run a lamport node",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config, c",
					Value: "lamport.toml",
					Usage: "lamport configuration `FILE`",
				},
			},
			Action: getAction(),
		},
	}
	return app
}

func getAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		if c == nil {
			return fmt.Errorf("Invalid context supplied to 'Action'")
		}
		cf := c.String("config")
		config, err := config.ReadConfig(cf)
		if err != nil {
			return fmt.Errorf("Error processing config file: %s", err)
		}
		node.Start(node.New(config))
		return nil
	}
}
