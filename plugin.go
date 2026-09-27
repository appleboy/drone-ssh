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

	easyssh "github.com/appleboy/easyssh-proxy"
)

var (
	errMissingHost          = errors.New("error: missing server host")
	errMissingPasswordOrKey = errors.New(
		"error: can't connect without a private SSH key or password",
	)
	errCommandTimeOut = errors.New("error: command timeout")
	envsFormat        = "export {NAME}={VALUE}"
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
		RequireTty        bool
	}

	// Plugin structure
	Plugin struct {
		Config Config
		Writer io.Writer
	}
)

// redacted returns a display-only copy, preserving whether credentials are set.
func (c Config) redacted() Config {
	for _, credential := range []*string{
		&c.Key, &c.Password, &c.Passphrase,
		&c.Proxy.Key, &c.Proxy.Password, &c.Proxy.Passphrase,
	} {
		if *credential != "" {
			*credential = "[REDACTED]"
		}
	}
	return c
}

func escapeArg(arg string) string {
	return "'" + strings.ReplaceAll(arg, "'", `'\''`) + "'"
}

func (p Plugin) hostPort(host string) (string, string) {
	hosts := strings.Split(host, ":")
	port := strconv.Itoa(p.Config.Port)
	if len(hosts) > 1 &&
		(p.Config.Protocol == easyssh.PROTOCOL_TCP ||
			p.Config.Protocol == easyssh.PROTOCOL_TCP4) {
		host = hosts[0]
		port = hosts[1]
	}

	return host, port
}

func (p Plugin) exec(host string) error {
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
		RequestPty:        p.Config.RequireTty,
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

	if p.Config.Debug {
		p.log(host, "======CMD======")
		p.log(host, strings.Join(p.Config.Script, "\n"))
		p.log(host, "======END======")
	}

	env := []string{}
	envNames := []string{}
	if p.Config.AllEnvs {
		allenvs := findEnvs("DRONE_", "PLUGIN_", "INPUT_", "GITHUB_")
		p.Config.Envs = append(p.Config.Envs, allenvs...)
	}
	for _, key := range p.Config.Envs {
		key = strings.ToUpper(key)
		if val, found := os.LookupEnv(key); found {
			envNames = append(envNames, key+"=[REDACTED]")
			env = append(
				env,
				p.format(p.Config.EnvsFormat, "{NAME}", key, "{VALUE}", escapeArg(val)),
			)
		}
	}

	if p.Config.Debug && len(env) > 0 {
		p.log(host, "======ENV======")
		p.log(host, strings.Join(envNames, "\n"))
		p.log(host, "======END======")
	}

	env = append(env, p.scriptCommands()...)
	p.Config.Script = env

	stdoutChan, stderrChan, doneChan, errChan, err := ssh.Stream(
		strings.Join(p.Config.Script, "\n"),
		p.Config.CommandTimeout,
	)
	if err != nil {
		return err
	}
	return p.readStream(host, stdoutChan, stderrChan, doneChan, errChan)
}

func (p Plugin) readStream(
	host string,
	stdoutChan, stderrChan <-chan string,
	doneChan <-chan bool,
	errChan <-chan error,
) error {
	var err error
	// read from the output channel until the done signal is passed
	var isTimeout bool
loop:
	for {
		select {
		case isTimeout = <-doneChan:
			break loop
		case outline, ok := <-stdoutChan:
			if !ok {
				stdoutChan = nil
				continue
			}
			if outline != "" {
				p.log(host, outline)
			}
		case errline, ok := <-stderrChan:
			if !ok {
				stderrChan = nil
				continue
			}
			if errline != "" {
				p.log(host, errline)
			}
		case streamErr, ok := <-errChan:
			if !ok {
				errChan = nil
			} else if streamErr != nil {
				err = streamErr
			}
		}
	}

	// get exit code or command error.
	if err != nil {
		return err
	}

	// command time out
	if !isTimeout {
		return errCommandTimeOut
	}
	return nil
}

// format string
func (p Plugin) format(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}

func (p Plugin) getWriter() io.Writer {
	if p.Writer != nil {
		return p.Writer
	}
	return os.Stdout
}

// log output to console
func (p Plugin) log(host string, message ...any) {
	w := p.getWriter()
	if count := len(p.Config.Host); count == 1 {
		fmt.Fprintf(w, "%s", fmt.Sprintln(message...))
		return
	}

	fmt.Fprintf(w, "%s: %s", host, fmt.Sprintln(message...))
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

	if err := p.execHosts(p.exec); err != nil {
		return err
	}

	w := p.getWriter()
	fmt.Fprintln(w, "===============================================")
	fmt.Fprintln(w, "✅ Successfully executed commands to all hosts.")
	fmt.Fprintln(w, "===============================================")

	return nil
}

// execHosts waits for all started hosts before returning. In sync mode, a failure
// stops subsequent hosts; in parallel mode, already started hosts finish normally.
func (p Plugin) execHosts(run func(string) error) error {
	if p.Config.Sync {
		for _, host := range p.Config.Host {
			if err := run(host); err != nil {
				return fmt.Errorf("%s: %w", host, err)
			}
		}
		return nil
	}

	var wg sync.WaitGroup
	// Each host sends at most one error, so sends never block on a receiver.
	errChannel := make(chan error, len(p.Config.Host))
	for _, host := range p.Config.Host {
		wg.Go(func() {
			if err := run(host); err != nil {
				errChannel <- fmt.Errorf("%s: %w", host, err)
			}
		})
	}
	wg.Wait()
	close(errChannel)
	// Preserve the first reported error, after every worker has exited.
	return <-errChannel
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
			commands = append(
				commands,
				"DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;",
			)
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
