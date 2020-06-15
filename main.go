package main

import (
	"log"
	"os"
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version string

func main() {
	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}
	app := cli.NewApp()
	app.Name = "Drone SSH"
	app.Usage = "Executing remote ssh commands"
	app.Copyright = "Copyright (c) 2019 Bo-Yi Wu"
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
			EnvVar: "PLUGIN_SSH_KEY,PLUGIN_KEY,SSH_KEY,KEY,INPUT_KEY",
		},
		cli.StringFlag{
			Name:   "ssh-passphrase",
			Usage:  "The purpose of the passphrase is usually to encrypt the private key.",
			EnvVar: "PLUGIN_SSH_PASSPHRASE,PLUGIN_PASSPHRASE,SSH_PASSPHRASE,PASSPHRASE,INPUT_PASSPHRASE",
		},
		cli.StringFlag{
			Name:   "key-path,i",
			Usage:  "ssh private key path",
			EnvVar: "PLUGIN_KEY_PATH,SSH_KEY_PATH,INPUT_KEY_PATH",
		},
		cli.StringFlag{
			Name:   "username,user,u",
			Usage:  "connect as user",
			EnvVar: "PLUGIN_USERNAME,PLUGIN_USER,SSH_USERNAME,USERNAME,INPUT_USERNAME",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "password,P",
			Usage:  "user password",
			EnvVar: "PLUGIN_PASSWORD,SSH_PASSWORD,PASSWORD,INPUT_PASSWORD",
		},
		cli.StringSliceFlag{
			Name:   "ciphers",
			Usage:  "The allowed cipher algorithms. If unspecified then a sensible",
			EnvVar: "PLUGIN_CIPHERS,SSH_CIPHERS,CIPHERS,INPUT_CIPHERS",
		},
		cli.BoolFlag{
			Name:   "useInsecureCipher",
			Usage:  "include more ciphers with use_insecure_cipher",
			EnvVar: "PLUGIN_USE_INSECURE_CIPHER,SSH_USE_INSECURE_CIPHER,USE_INSECURE_CIPHER,INPUT_USE_INSECURE_CIPHER",
		},
		cli.StringFlag{
			Name:   "fingerprint",
			Usage:  "fingerprint SHA256 of the host public key, default is to skip verification",
			EnvVar: "PLUGIN_FINGERPRINT,SSH_FINGERPRINT,FINGERPRINT,INPUT_FINGERPRINT",
		},
		cli.StringSliceFlag{
			Name:     "host,H",
			Usage:    "connect to host",
			EnvVar:   "PLUGIN_HOST,SSH_HOST,HOST,INPUT_HOST",
			FilePath: ".host",
		},
		cli.IntFlag{
			Name:   "port,p",
			Usage:  "connect to port",
			EnvVar: "PLUGIN_PORT,SSH_PORT,PORT,INPUT_PORT",
			Value:  22,
		},
		cli.BoolFlag{
			Name:   "sync",
			Usage:  "sync mode",
			EnvVar: "PLUGIN_SYNC,SYNC,INPUT_SYNC",
		},
		cli.DurationFlag{
			Name:   "timeout,t",
			Usage:  "connection timeout",
			EnvVar: "PLUGIN_TIMEOUT,SSH_TIMEOUT,TIMEOUT,INPUT_TIMEOUT",
			Value:  30 * time.Second,
		},
		cli.DurationFlag{
			Name:   "command.timeout,T",
			Usage:  "command timeout",
			EnvVar: "PLUGIN_COMMAND_TIMEOUT,SSH_COMMAND_TIMEOUT,COMMAND_TIMEOUT,INPUT_COMMAND_TIMEOUT",
			Value:  10 * time.Minute,
		},
		cli.StringSliceFlag{
			Name:   "script,s",
			Usage:  "execute commands",
			EnvVar: "PLUGIN_SCRIPT,SSH_SCRIPT,SCRIPT",
		},
		cli.StringFlag{
			Name:   "script.string",
			Usage:  "execute single commands for github action",
			EnvVar: "INPUT_SCRIPT",
		},
		cli.BoolFlag{
			Name:   "script.stop",
			Usage:  "stop script after first failure",
			EnvVar: "PLUGIN_SCRIPT_STOP,STOP,INPUT_SCRIPT_STOP",
		},
		cli.StringFlag{
			Name:   "proxy.ssh-key",
			Usage:  "private ssh key of proxy",
			EnvVar: "PLUGIN_PROXY_SSH_KEY,PLUGIN_PROXY_KEY,PROXY_SSH_KEY,INPUT_PROXY_KEY",
		},
		cli.StringFlag{
			Name:   "proxy.ssh-passphrase",
			Usage:  "The purpose of the passphrase is usually to encrypt the private key.",
			EnvVar: "PLUGIN_PROXY_SSH_PASSPHRASE,PLUGIN_PROXY_PASSPHRASE,PROXY_SSH_PASSPHRASE,PROXY_PASSPHRASE,INPUT_PROXY_PASSPHRASE",
		},
		cli.StringFlag{
			Name:   "proxy.key-path",
			Usage:  "ssh private key path of proxy",
			EnvVar: "PLUGIN_PROXY_KEY_PATH,PROXY_SSH_KEY_PATH,INPUT_PROXY_KEY_PATH",
		},
		cli.StringFlag{
			Name:   "proxy.username",
			Usage:  "connect as user of proxy",
			EnvVar: "PLUGIN_PROXY_USERNAME,PLUGIN_PROXY_USER,PROXY_SSH_USERNAME,INPUT_PROXY_USERNAME",
			Value:  "root",
		},
		cli.StringFlag{
			Name:   "proxy.password",
			Usage:  "user password of proxy",
			EnvVar: "PLUGIN_PROXY_PASSWORD,PROXY_SSH_PASSWORD,INPUT_PROXY_PASSWORD",
		},
		cli.StringFlag{
			Name:   "proxy.host",
			Usage:  "connect to host of proxy",
			EnvVar: "PLUGIN_PROXY_HOST,PROXY_SSH_HOST,INPUT_PROXY_HOST",
		},
		cli.StringFlag{
			Name:   "proxy.port",
			Usage:  "connect to port of proxy",
			EnvVar: "PLUGIN_PROXY_PORT,PROXY_SSH_PORT,INPUT_PROXY_PORT",
			Value:  "22",
		},
		cli.DurationFlag{
			Name:   "proxy.timeout",
			Usage:  "proxy connection timeout",
			EnvVar: "PLUGIN_PROXY_TIMEOUT,PROXY_SSH_TIMEOUT,INPUT_PROXY_TIMEOUT",
		},
		cli.StringSliceFlag{
			Name:   "proxy.ciphers",
			Usage:  "The allowed cipher algorithms. If unspecified then a sensible",
			EnvVar: "PLUGIN_PROXY_CIPHERS,SSH_PROXY_CIPHERS,PROXY_CIPHERS,INPUT_PROXY_CIPHERS",
		},
		cli.BoolFlag{
			Name:   "proxy.useInsecureCipher",
			Usage:  "include more ciphers with use_insecure_cipher",
			EnvVar: "PLUGIN_PROXY_USE_INSECURE_CIPHER,SSH_PROXY_USE_INSECURE_CIPHER,PROXY_USE_INSECURE_CIPHER,INPUT_PROXY_USE_INSECURE_CIPHER",
		},
		cli.StringFlag{
			Name:   "proxy.fingerprint",
			Usage:  "fingerprint SHA256 of the host public key, default is to skip verification",
			EnvVar: "PLUGIN_PROXY_FINGERPRINT,SSH_PROXY_FINGERPRINT,PROXY_FINGERPRINT,INPUT_PROXY_FINGERPRINT",
		},
		cli.StringSliceFlag{
			Name:   "envs",
			Usage:  "pass environment variable to shell script",
			EnvVar: "PLUGIN_ENVS,INPUT_ENVS",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug mode",
			EnvVar: "PLUGIN_DEBUG,DEBUG,INPUT_DEBUG",
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
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	scripts := c.StringSlice("script")
	if s := c.String("script.string"); s != "" {
		scripts = append(scripts, s)
	}
	plugin := Plugin{
		Config: Config{
			Key:               c.String("ssh-key"),
			KeyPath:           c.String("key-path"),
			Username:          c.String("user"),
			Password:          c.String("password"),
			Passphrase:        c.String("ssh-passphrase"),
			Fingerprint:       c.String("fingerprint"),
			Host:              c.StringSlice("host"),
			Port:              c.Int("port"),
			Timeout:           c.Duration("timeout"),
			CommandTimeout:    c.Duration("command.timeout"),
			Script:            scripts,
			ScriptStop:        c.Bool("script.stop"),
			Envs:              c.StringSlice("envs"),
			Debug:             c.Bool("debug"),
			Sync:              c.Bool("sync"),
			Ciphers:           c.StringSlice("ciphers"),
			UseInsecureCipher: c.Bool("useInsecureCipher"),
			Proxy: easyssh.DefaultConfig{
				Key:               c.String("proxy.ssh-key"),
				KeyPath:           c.String("proxy.key-path"),
				User:              c.String("proxy.username"),
				Password:          c.String("proxy.password"),
				Passphrase:        c.String("proxy.ssh-passphrase"),
				Fingerprint:       c.String("proxy.fingerprint"),
				Server:            c.String("proxy.host"),
				Port:              c.String("proxy.port"),
				Timeout:           c.Duration("proxy.timeout"),
				Ciphers:           c.StringSlice("proxy.ciphers"),
				UseInsecureCipher: c.Bool("proxy.useInsecureCipher"),
			},
		},
		Writer: os.Stdout,
	}

	return plugin.Exec()
}
