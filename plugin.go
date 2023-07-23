package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/appleboy/easyssh-proxy"
)

var (
	errMissingHost          = errors.New("Error: missing server host")
	errMissingPasswordOrKey = errors.New("Error: can't connect without a private SSH key or password")
	errCommandTimeOut       = errors.New("Error: command timeout")
	envsFormat              = "export {NAME}={VALUE}"
)

type (
	// Config for the plugin.
	Config struct {
		Key               string
		Passphrase        string
		KeyPath           string
		Username          string
		Password          string
		Host              []string
		Port              int
		Protocol          easyssh.Protocol
		Fingerprint       string
		Timeout           time.Duration
		CommandTimeout    time.Duration
		Script            []string
		ScriptStop        bool
		Envs              []string
		Proxy             easyssh.DefaultConfig
		Debug             bool
		Sync              bool
		Ciphers           []string
		UseInsecureCipher bool
		EnvsFormat        string
		AllEnvs           bool
	}

	// Plugin structure
	Plugin struct {
		Config Config
		Writer io.Writer
	}
)

func escapeArg(arg string) string {
	return "'" + strings.Replace(arg, "'", `'\''`, -1) + "'"
}

func (p Plugin) hostPort(host string) (string, string) {
	hosts := strings.Split(host, ":")
	port := strconv.Itoa(p.Config.Port)
	if len(hosts) > 1 {
		host = hosts[0]
		port = hosts[1]
	}

	return host, port
}

func (p Plugin) exec(host string, wg *sync.WaitGroup, errChannel chan error) {
	defer wg.Done()
	host, port := p.hostPort(host)
	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		Server:            host,
		User:              p.Config.Username,
		Password:          p.Config.Password,
		Port:              port,
		Protocol:          p.Config.Protocol,
		Key:               p.Config.Key,
		KeyPath:           p.Config.KeyPath,
		Passphrase:        p.Config.Passphrase,
		Timeout:           p.Config.Timeout,
		Ciphers:           p.Config.Ciphers,
		Fingerprint:       p.Config.Fingerprint,
		UseInsecureCipher: p.Config.UseInsecureCipher,
		Proxy: easyssh.DefaultConfig{
			Server:            p.Config.Proxy.Server,
			User:              p.Config.Proxy.User,
			Password:          p.Config.Proxy.Password,
			Port:              p.Config.Proxy.Port,
			Protocol:          p.Config.Proxy.Protocol,
			Key:               p.Config.Proxy.Key,
			KeyPath:           p.Config.Proxy.KeyPath,
			Passphrase:        p.Config.Proxy.Passphrase,
			Timeout:           p.Config.Proxy.Timeout,
			Ciphers:           p.Config.Proxy.Ciphers,
			Fingerprint:       p.Config.Proxy.Fingerprint,
			UseInsecureCipher: p.Config.Proxy.UseInsecureCipher,
		},
	}

	p.log(host, "======CMD======")
	p.log(host, strings.Join(p.Config.Script, "\n"))
	p.log(host, "======END======")

	env := []string{}
	if p.Config.AllEnvs {
		allenvs := findEnvs("DRONE_", "PLUGIN_", "INPUT_", "GITHUB_")
		p.Config.Envs = append(p.Config.Envs, allenvs...)
	}
	for _, key := range p.Config.Envs {
		key = strings.ToUpper(key)
		if val, found := os.LookupEnv(key); found {
			env = append(env, p.format(p.Config.EnvsFormat, "{NAME}", key, "{VALUE}", escapeArg(val)))
		}
	}

	p.Config.Script = append(env, p.scriptCommands()...)

	if p.Config.Debug {
		p.log(host, "======ENV======")
		p.log(host, strings.Join(env, "\n"))
		p.log(host, "======END======")
	}

	stdoutChan, stderrChan, doneChan, errChan, err := ssh.Stream(strings.Join(p.Config.Script, "\n"), p.Config.CommandTimeout)
	if err != nil {
		errChannel <- err
		return
	}
	// read from the output channel until the done signal is passed
	isTimeout := true
loop:
	for {
		select {
		case isTimeout = <-doneChan:
			break loop
		case outline := <-stdoutChan:
			if outline != "" {
				p.log(host, "out:", outline)
			}
		case errline := <-stderrChan:
			if errline != "" {
				p.log(host, "err:", errline)
			}
		case err = <-errChan:
		}
	}

	// get exit code or command error.
	if err != nil {
		errChannel <- err
	}

	// command time out
	if !isTimeout {
		errChannel <- errCommandTimeOut
	}
}

// format string
func (p Plugin) format(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

// log output to console
func (p Plugin) log(host string, message ...interface{}) {
	if p.Writer == nil {
		p.Writer = os.Stdout
	}
	if count := len(p.Config.Host); count == 1 {
		fmt.Fprintf(p.Writer, "%s", fmt.Sprintln(message...))
		return
	}

	fmt.Fprintf(p.Writer, "%s: %s", host, fmt.Sprintln(message...))
}

// Exec executes the plugin.
func (p Plugin) Exec() error {
	p.Config.Host = trimValues(p.Config.Host)

	if len(p.Config.Host) == 0 {
		return errMissingHost
	}

	if len(p.Config.Key) == 0 && len(p.Config.Password) == 0 && len(p.Config.KeyPath) == 0 {
		return errMissingPasswordOrKey
	}

	if p.Config.EnvsFormat == "" {
		p.Config.EnvsFormat = envsFormat
	}

	wg := sync.WaitGroup{}
	wg.Add(len(p.Config.Host))
	errChannel := make(chan error)
	finished := make(chan struct{})
	if p.Config.Sync {
		go func() {
			for _, host := range p.Config.Host {
				p.exec(host, &wg, errChannel)
			}
		}()
	} else {
		for _, host := range p.Config.Host {
			go p.exec(host, &wg, errChannel)
		}
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChannel:
		if err != nil {
			return err
		}
	}

	fmt.Println("==============================================")
	fmt.Println("✅ Successfully executed commands to all host.")
	fmt.Println("==============================================")

	return nil
}

func (p Plugin) scriptCommands() []string {
	scripts := []string{}

	for _, cmd := range p.Config.Script {
		if p.Config.ScriptStop {
			scripts = append(scripts, strings.Split(cmd, "\n")...)
		} else {
			scripts = append(scripts, cmd)
		}
	}

	commands := make([]string, 0)

	for _, cmd := range scripts {
		cmd = strings.TrimSpace(cmd)
		if strings.TrimSpace(cmd) == "" {
			continue
		}
		commands = append(commands, cmd)
		if p.Config.ScriptStop && cmd[(len(cmd)-1):] != "\\" {
			commands = append(commands, "DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;")
		}
	}

	return commands
}

func trimValues(keys []string) []string {
	var newKeys []string

	for _, value := range keys {
		value = strings.TrimSpace(value)
		if len(value) == 0 {
			continue
		}

		newKeys = append(newKeys, value)
	}

	return newKeys
}

// Find all envs from specified prefix
func findEnvs(prefix ...string) []string {
	envs := []string{}
	for _, e := range os.Environ() {
		for _, p := range prefix {
			if strings.HasPrefix(e, p) {
				e = strings.Split(e, "=")[0]
				envs = append(envs, e)
				break
			}
		}
	}
	return envs
}
