package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/appleboy/easyssh-proxy"
)

const (
	missingHostOrUser    = "Error: missing server host or user"
	missingPasswordOrKey = "Error: can't connect without a private SSH key or password"
	commandTimeOut       = "Error: command timeout"
	setPasswordandKey    = "can't set password and key at the same time"
)

type (
	// Config for the plugin.
	Config struct {
		Key            string
		KeyPath        string
		UserName       string
		Password       string
		Host           []string
		Port           int
		Timeout        time.Duration
		CommandTimeout int
		Script         []string
		Secrets        []string
		Envs           []string
		Proxy          easyssh.DefaultConfig
		Debug          bool
		Sync           bool
	}

	// Plugin structure
	Plugin struct {
		Config Config
	}
)

func (p Plugin) exec(host string, wg *sync.WaitGroup, errChannel chan error) {
	// Create MakeConfig instance with remote username, server address and path to private key.
	ssh := &easyssh.MakeConfig{
		Server:   host,
		User:     p.Config.UserName,
		Password: p.Config.Password,
		Port:     strconv.Itoa(p.Config.Port),
		Key:      p.Config.Key,
		KeyPath:  p.Config.KeyPath,
		Timeout:  p.Config.Timeout,
		Proxy: easyssh.DefaultConfig{
			Server:   p.Config.Proxy.Server,
			User:     p.Config.Proxy.User,
			Password: p.Config.Proxy.Password,
			Port:     p.Config.Proxy.Port,
			Key:      p.Config.Proxy.Key,
			KeyPath:  p.Config.Proxy.KeyPath,
			Timeout:  p.Config.Proxy.Timeout,
		},
	}

	p.log(host, "======CMD======")
	p.log(host, strings.Join(p.Config.Script, "\n"))
	p.log(host, "======END======")

	env := []string{}
	for _, key := range p.Config.Envs {
		key = strings.ToUpper(key)
		val := os.Getenv(key)
		val = strings.Replace(val, " ", "", -1)
		env = append(env, key+"='"+val+"'")
	}

	p.Config.Script = append(env, p.Config.Script...)

	if p.Config.Debug {
		p.log(host, "======ENV======")
		p.log(host, strings.Join(env, "\n"))
		p.log(host, "======END======")
	}

	stdoutChan, stderrChan, doneChan, errChan, err := ssh.Stream(strings.Join(p.Config.Script, "\n"), p.Config.CommandTimeout)
	if err != nil {
		errChannel <- err
	} else {
		// read from the output channel until the done signal is passed
		isTimeout := true
	loop:
		for {
			select {
			case isTimeout = <-doneChan:
				break loop
			case outline := <-stdoutChan:
				p.log(host, "out:", outline)
			case errline := <-stderrChan:
				p.log(host, "err:", errline)
			case err = <-errChan:
			}
		}

		// get exit code or command error.
		if err != nil {
			errChannel <- err
		}

		// command time out
		if !isTimeout {
			errChannel <- fmt.Errorf(commandTimeOut)
		}
	}

	wg.Done()
}

func (p Plugin) log(host string, message ...interface{}) {
	if count := len(p.Config.Host); count == 1 {
		fmt.Printf("%s", fmt.Sprintln(message...))
	} else {
		fmt.Printf("%s: %s", host, fmt.Sprintln(message...))
	}
}

// Exec executes the plugin.
func (p Plugin) Exec() error {
	if len(p.Config.Host) == 0 && len(p.Config.UserName) == 0 {
		return fmt.Errorf(missingHostOrUser)
	}

	if len(p.Config.Key) == 0 && len(p.Config.Password) == 0 && len(p.Config.KeyPath) == 0 {
		return fmt.Errorf(missingPasswordOrKey)
	}

	if len(p.Config.Key) != 0 && len(p.Config.Password) != 0 {
		return fmt.Errorf(setPasswordandKey)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(p.Config.Host))
	errChannel := make(chan error, 1)
	finished := make(chan bool, 1)
	for _, host := range p.Config.Host {
		if p.Config.Sync {
			p.exec(host, &wg, errChannel)
		} else {
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
			fmt.Println("drone-ssh error: ", err)
			return err
		}
	}

	fmt.Println("==========================================")
	fmt.Println("Successfully executed commands to all host.")
	fmt.Println("==========================================")

	return nil
}
