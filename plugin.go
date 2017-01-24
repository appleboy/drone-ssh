package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	missingHostOrUser    = "Error: missing server host or user"
	missingPasswordOrKey = "Error: can't connect without a private SSH key or password"
	unableConnectServer  = "Error: Failed to start a SSH session"
	failParsePrivateKey  = "Error: Failed to parse private key"
	sshKeyNotFound       = "ssh: no key found"
)

type (
	// Config for the plugin.
	Config struct {
		Key      string
		KeyPath  string
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

// returns ssh.Signer from user you running app home path + cutted key path.
// (ex. pubkey,err := getKeyFile("/.ssh/id_rsa") )
func getKeyFile(keypath string) (ssh.Signer, error) {
	buf, err := ioutil.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	pubkey, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

// Exec executes the plugin.
func (p Plugin) Exec() error {
	if len(p.Config.Host) == 0 && p.Config.User == "" {
		return fmt.Errorf(missingHostOrUser)
	}

	if p.Config.Key == "" && p.Config.Password == "" && p.Config.KeyPath == "" {
		return fmt.Errorf(missingPasswordOrKey)
	}

	for i, host := range p.Config.Host {
		addr := net.JoinHostPort(
			host,
			strconv.Itoa(p.Config.Port),
		)

		// auths holds the detected ssh auth methods
		auths := []ssh.AuthMethod{}

		if p.Config.KeyPath != "" {
			pubkey, err := getKeyFile(p.Config.KeyPath)

			if err != nil {
				return err
			}

			auths = append(auths, ssh.PublicKeys(pubkey))
		}

		if p.Config.Key != "" {
			signer, err := ssh.ParsePrivateKey([]byte(p.Config.Key))

			if err != nil {
				return fmt.Errorf(failParsePrivateKey)
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
			return fmt.Errorf(unableConnectServer)
		}

		session, _ := client.NewSession()
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
