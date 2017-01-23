package main

import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type (
	// Config for the plugin.
	Config struct {
		Key      string
		User     string
		Password string
		Host     []string
		Port     int
		Sleep    int
		Timeout  time.Duration
		Script   []string
	}

	// Plugin structure
	Plugin struct {
		Config Config
	}
)

// Exec executes the plugin.
func (p Plugin) Exec() error {
	if p.Config.Key == "" && p.Config.Password == "" {
		return errors.New("Error: can't connect without a private SSH key or password")
	}

	for i, host := range p.Config.Host {
		addr := net.JoinHostPort(
			host,
			strconv.Itoa(p.Config.Port),
		)

		// auths holds the detected ssh auth methods
		auths := []ssh.AuthMethod{}

		if p.Config.Key != "" {
			signer, err := ssh.ParsePrivateKey([]byte(p.Config.Key))

			if err != nil {
				return errors.New("Error: Failed to parse private key. %s", err)
			}

			auths = append(auths, ssh.PublicKeys(signer))
		}

		// figure out what auths are requested, what is supported
		if p.Config.Password != "" {
			auths = append(auths, ssh.Password(p.Config.Password))
		}

		config := &ssh.ClientConfig{
			Timeout: p.Config.Timeout,
			User:    p.Config.User,
			Auth:    auths,
		}

		log.Printf("+ ssh %s@%s -p %d\n", p.Config.User, addr, p.Config.Port)
		client, err := ssh.Dial("tcp", addr, config)

		if err != nil {
			return errors.New("Error: Failed to dial to server. %s", err)
		}

		session, err := client.NewSession()

		if err != nil {
			return errors.New("Error: Failed to start a SSH session. %s", err)
		}

		defer session.Close()

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		if err := session.Run(strings.Join(p.Config.Script, "\n")); err != nil {
			return err
		}

		if p.Config.Sleep != 0 && i != len(p.Config.Host)-1 {
			log.Printf("+ sleep %d\n", p.Config.Sleep)
			time.Sleep(time.Duration(p.Config.Sleep) * time.Second)
		}
	}

	return nil
}
