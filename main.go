package main

import (
	"fmt"
	"os"

	"github.com/appleboy/easyssh-proxy"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// build number set at compile-time
var build = "0"

// Version set at compile-time
var Version string

func main() {
	if Version == "" {
		Version = fmt.Sprintf("1.3.1+%s", build)
	}

	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

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
			Name:   "key-path,i",
			Usage:  "ssh private key path",
			EnvVar: "PLUGIN_KEY_PATH,SSH_KEY_PATH",
		},
		cli.StringFlag{
			Name:   "username,user,u",
			Usage:  "connect as user",
			EnvVar: "PLUGIN_USERNAME,PLUGIN_USER,SSH_USERNAME",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "password,P",
			Usage:  "user password",
			EnvVar: "PLUGIN_PASSWORD,SSH_PASSWORD",
		},
		cli.StringSliceFlag{
			Name:   "host,H",
			Usage:  "connect to host",
			EnvVar: "PLUGIN_HOST,SSH_HOST",
		},
		cli.IntFlag{
			Name:   "port,p",
			Usage:  "connect to port",
			EnvVar: "PLUGIN_PORT,SSH_PORT",
			Value:  22,
		},
		cli.BoolFlag{
			Name:   "sync",
			Usage:  "sync mode",
			EnvVar: "PLUGIN_SYNC",
		},
		cli.DurationFlag{
			Name:   "timeout,t",
			Usage:  "connection timeout",
			EnvVar: "PLUGIN_TIMEOUT,SSH_TIMEOUT",
		},
		cli.IntFlag{
			Name:   "command.timeout,T",
			Usage:  "command timeout",
			EnvVar: "PLUGIN_COMMAND_TIMEOUT,SSH_COMMAND_TIMEOUT",
			Value:  60,
		},
		cli.StringSliceFlag{
			Name:   "script,s",
			Usage:  "execute commands",
			EnvVar: "PLUGIN_SCRIPT,SSH_SCRIPT",
		},
		cli.BoolFlag{
			Name:   "script.stop",
			Usage:  "stop script after first failure",
			EnvVar: "PLUGIN_SCRIPT_STOP",
		},
		cli.StringFlag{
			Name:   "proxy.ssh-key",
			Usage:  "private ssh key of proxy",
			EnvVar: "PLUGIN_PROXY_SSH_KEY,PLUGIN_PROXY_KEY,PROXY_SSH_KEY",
		},
		cli.StringFlag{
			Name:   "proxy.key-path",
			Usage:  "ssh private key path of proxy",
			EnvVar: "PLUGIN_PROXY_KEY_PATH,PROXY_SSH_KEY_PATH",
		},
		cli.StringFlag{
			Name:   "proxy.username",
			Usage:  "connect as user of proxy",
			EnvVar: "PLUGIN_PROXY_USERNAME,PLUGIN_PROXY_USER,PROXY_SSH_USERNAME",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "proxy.password",
			Usage:  "user password of proxy",
			EnvVar: "PLUGIN_PROXY_PASSWORD,PROXY_SSH_PASSWORD",
		},
		cli.StringFlag{
			Name:   "proxy.host",
			Usage:  "connect to host of proxy",
			EnvVar: "PLUGIN_PROXY_HOST,PROXY_SSH_HOST",
		},
		cli.StringFlag{
			Name:   "proxy.port",
			Usage:  "connect to port of proxy",
			EnvVar: "PLUGIN_PROXY_PORT,PROXY_SSH_PORT",
			Value:  "22",
		},
		cli.DurationFlag{
			Name:   "proxy.timeout",
			Usage:  "proxy connection timeout",
			EnvVar: "PLUGIN_PROXY_TIMEOUT,PROXY_SSH_TIMEOUT",
		},
		cli.StringSliceFlag{
			Name:   "secrets",
			Usage:  "plugin secret",
			EnvVar: "PLUGIN_SECRETS",
		},
		cli.StringSliceFlag{
			Name:   "envs",
			Usage:  "Pass envs",
			EnvVar: "PLUGIN_ENVS",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "PLUGIN_DEBUG",
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

	if err := app.Run(os.Args); err != nil {
		fmt.Println("drone-ssh error: ", err)
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Config: Config{
			Key:            c.String("ssh-key"),
			KeyPath:        c.String("key-path"),
			Username:       c.String("user"),
			Password:       c.String("password"),
			Host:           c.StringSlice("host"),
			Port:           c.Int("port"),
			Timeout:        c.Duration("timeout"),
			CommandTimeout: c.Int("command.timeout"),
			Script:         c.StringSlice("script"),
			ScriptStop:     c.Bool("script.stop"),
			Secrets:        c.StringSlice("secrets"),
			Envs:           c.StringSlice("envs"),
			Debug:          c.Bool("debug"),
			Sync:           c.Bool("sync"),
			Proxy: easyssh.DefaultConfig{
				Key:      c.String("proxy.ssh-key"),
				KeyPath:  c.String("proxy.key-path"),
				User:     c.String("proxy.username"),
				Password: c.String("proxy.password"),
				Server:   c.String("proxy.host"),
				Port:     c.String("proxy.port"),
				Timeout:  c.Duration("proxy.timeout"),
			},
		},
		Writer: os.Stdout,
	}

	return plugin.Exec()
}
