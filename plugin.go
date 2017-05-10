package main

import (
	"fmt"
	"log"
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
		Proxy          easyssh.DefaultConfig
	}

	// Plugin structure
	Plugin struct {
		Config Config
	}
)

func (p Plugin) log(host string, message ...interface{}) {
	log.Printf("%s: %s", host, fmt.Sprintln(message...))
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
		go func(host string) {
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

			p.log(host, "commands: ", strings.Join(p.Config.Script, "\n"))
			stdoutChan, stderrChan, doneChan, errChan, err := ssh.Stream(strings.Join(p.Config.Script, "\n"), p.Config.CommandTimeout)
			if err != nil {
				errChannel <- err
			} else {
				// read from the output channel until the done signal is passed
				stillGoing := true
				isTimeout := true
				for stillGoing {
					select {
					case isTimeout = <-doneChan:
						stillGoing = false
					case outline := <-stdoutChan:
						p.log(host, "outputs:", outline)
					case errline := <-stderrChan:
						p.log(host, "errors:", errline)
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
		}(host)
	}

	go func() {
		wg.Wait()
		close(finished)
	}()

	select {
	case <-finished:
	case err := <-errChannel:
		if err != nil {
			log.Println("drone-ssh error: ", err)
			return err
		}
	}

	log.Println("Successfully executed commands to all host.")

	return nil
}
