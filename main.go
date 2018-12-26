//go:generate go run types/codegen/cleanup/main.go
//go:generate go run types/codegen/main.go

package main

import (
	"os"

	"github.com/rancher/alertmanager-webhook-adapter/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	VERSION = "v0.0.0-dev"
)

func main() {
	app := cli.NewApp()
	app.Name = "alertmanager-webhook-adapter"
	app.Version = VERSION
	app.Usage = "alertmanager-webhook-adapter needs help!"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "configpath",
			EnvVar: "CONFIGPATH",
			Value:  "/etc/alertmanager/config/webhooktemplates.yaml",
		},
		cli.IntFlag{
			Name:   "port",
			EnvVar: "PORT",
			Value:  7890,
		},
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	server := server.Server{
		Debug:      c.Bool("debug"),
		Port:       c.Int("port"),
		ConfigPath: c.String("configpath"),
	}
	server.Start()
	return nil
}
