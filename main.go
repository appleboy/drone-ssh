package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
	"golang.org/x/crypto/ssh"
)

var (
	buildDate string
)

func main() {
	fmt.Printf("Drone SSH Plugin built at %s\n", buildDate)

	workspace := drone.Workspace{}
	vargs := Params{}

	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	for i, host := range vargs.Host.Slice() {
		err := run(workspace.Keys, &vargs, host)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if vargs.Sleep != 0 && i != vargs.Host.Len()-1 {
			fmt.Printf("$ sleep %d\n", vargs.Sleep)
			time.Sleep(time.Duration(vargs.Sleep) * time.Second)
		}
	}
}

func run(key *drone.Key, params *Params, host string) error {
	if params.Login == "" {
		params.Login = "root"
	}

	if params.Port == 0 {
		params.Port = 22
	}

	addr := net.JoinHostPort(
		host,
		strconv.Itoa(params.Port),
	)

	fmt.Printf("$ ssh %s@%s -p %d\n", params.Login, addr, params.Port)
	signer, err := ssh.ParsePrivateKey([]byte(key.Private))

	if err != nil {
		return fmt.Errorf("Error: Failed to parse private key. %s", err)
	}

	config := &ssh.ClientConfig{
		User: params.Login,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

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

	return session.Run(strings.Join(params.Commands, "\n"))
}
