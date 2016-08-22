package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "ssh plugin"
	app.Usage = "ssh plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ssh-key",
			Usage:  "private ssh key",
			EnvVar: "PLUGIN_SSH_KEY,SSH_KEY",
		},
		cli.StringFlag{
			Name:   "user",
			Usage:  "connect as user",
			EnvVar: "PLUGIN_USER,SSH_USER",
			Value:  "root",
		},
		cli.StringSliceFlag{
			Name:   "host",
			Usage:  "connect to host",
			EnvVar: "PLUGIN_HOST,SSH_HOST",
		},
		cli.IntFlag{
			Name:   "port",
			Usage:  "connect to port",
			EnvVar: "PLUGIN_PORT,SSH_PORT",
			Value:  22,
		},
		cli.IntFlag{
			Name:   "sleep",
			Usage:  "sleep between hosts",
			EnvVar: "PLUGIN_SLEEP,SSH_SLEEP",
		},
		cli.DurationFlag{
			Name:   "timeout",
			Usage:  "connection timeout",
			EnvVar: "PLUGIN_TIMEOUT,SSH_TIMEOUT",
		},
		cli.StringSliceFlag{
			Name:   "commands",
			Usage:  "execute commands",
			EnvVar: "PLUGIN_COMMANDS,SSH_COMMANDS",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Config: Config{
			Key:      c.String("ssh-key"),
			User:     c.String("user"),
			Host:     c.StringSlice("host"),
			Port:     c.Int("port"),
			Sleep:    c.Int("sleep"),
			Timeout:  c.Duration("timeout"),
			Commands: c.StringSlice("commands"),
		},
	}

	return plugin.Exec()
}
