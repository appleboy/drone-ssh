package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/appleboy/drone-ssh/easyssh"
)

var wg sync.WaitGroup

const (
	missingHostOrUser    = "Error: missing server host or user"
	missingPasswordOrKey = "Error: can't connect without a private SSH key or password"
)

type (
	// Config for the plugin.
	Config struct {
		Key      string
		KeyPath  string
		UserName string
		Password string
		Host     []string
		Port     int
		Timeout  time.Duration
		Script   []string
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
	if len(p.Config.Host) == 0 && p.Config.UserName == "" {
		return fmt.Errorf(missingHostOrUser)
	}

	if p.Config.Key == "" && p.Config.Password == "" && p.Config.KeyPath == "" {
		return fmt.Errorf(missingPasswordOrKey)
	}

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
			}

			p.log(host, "commands: ", strings.Join(p.Config.Script, "\n"))
			response, err := ssh.Run(strings.Join(p.Config.Script, "\n"))
			p.log(host, "outputs:", response)

			if err != nil {
				errChannel <- err
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

	log.Println("Successfully executed commnads to all host.")

	return nil
}
