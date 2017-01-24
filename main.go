package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version = "v1.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "Drone SSH"
	app.Usage = "Executing remote ssh commands"
	app.Copyright = "Copyright (c) 2017 Bo-Yi Wu"
	app.Authors = []cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "ssh-key",
			Usage:  "private ssh key",
			EnvVar: "PLUGIN_SSH_KEY,PLUGIN_KEY,SSH_KEY",
		},
		cli.StringFlag{
			Name:   "key-path",
			Usage:  "ssh private key path",
			EnvVar: "PLUGIN_KEY_PATH,SSH_KEY_PATH",
		},
		cli.StringFlag{
			Name:   "user",
			Usage:  "connect as user",
			EnvVar: "PLUGIN_USER,SSH_USER",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "user password",
			EnvVar: "PLUGIN_PASSWORD,SSH_PASSWORD",
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
			Name:   "script",
			Usage:  "execute commands",
			EnvVar: "PLUGIN_SCRIPT,SSH_SCRIPT",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	// Override a template
	cli.AppHelpTemplate = `
________                                         _________ _________ ___ ___
\______ \_______  ____   ____   ____            /   _____//   _____//   |   \
 |    |  \_  __ \/  _ \ /    \_/ __ \   ______  \_____  \ \_____  \/    ~    \
 |    |   \  | \(  <_> )   |  \  ___/  /_____/  /        \/        \    Y    /
/_______  /__|   \____/|___|  /\___  >         /_______  /_______  /\___|_  /
        \/                  \/     \/                  \/        \/       \/
                                                    version: {{.Version}}
NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
REPOSITORY:
    Github: https://github.com/appleboy/drone-ssh
`

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Config: Config{
			Key:      c.String("ssh-key"),
			KeyPath:  c.String("key-path"),
			User:     c.String("user"),
			Password: c.String("password"),
			Host:     c.StringSlice("host"),
			Port:     c.Int("port"),
			Sleep:    c.Int("sleep"),
			Timeout:  c.Duration("timeout"),
			Script:   c.StringSlice("script"),
		},
	}

	return plugin.Exec()
}
