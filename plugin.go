package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type (
	Config struct {
		Key      string        `json:"key"`
		User     string        `json:"user"`
		Host     []string      `json:"host"`
		Port     int           `json:"port"`
		Sleep    int           `json:"sleep"`
		Timeout  time.Duration `json:"timeout"`
		Script   []string      `json:"script"`
	}

	Plugin struct {
		Config Config
	}
)

func (p Plugin) Exec() error {
	if p.Config.Key == "" {
		return fmt.Errorf("Error: Can't connect without a private SSH key.")
	}

	for i, host := range p.Config.Host {
		addr := net.JoinHostPort(
			host,
			strconv.Itoa(p.Config.Port),
		)

		signer, err := ssh.ParsePrivateKey([]byte(p.Config.Key))

		if err != nil {
			return fmt.Errorf("Error: Failed to parse private key. %s", err)
		}

		config := &ssh.ClientConfig{
			Timeout: p.Config.Timeout,
			User:    p.Config.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
		}

		fmt.Printf("+ ssh %s@%s -p %d\n", p.Config.User, addr, p.Config.Port)
		client, err := ssh.Dial("tcp", addr, config)

		if err != nil {
			return fmt.Errorf("Error: Failed to dial to server. %s", err)
		}

		session, err := client.NewSession()

		if err != nil {
			return fmt.Errorf("Error: Failed to start a SSH session. %s", err)
		}

		defer session.Close()

		session.Stdout = os.Stdout
		session.Stderr = os.Stderr

		if err := session.Run(strings.Join(p.Config.Script, "\n")); err != nil {
			return err
		}

		if p.Config.Sleep != 0 && i != len(p.Config.Host)-1 {
			fmt.Printf("+ sleep %d\n", p.Config.Sleep)
			time.Sleep(time.Duration(p.Config.Sleep) * time.Second)
		}
	}

	return nil
}
