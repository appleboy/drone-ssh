package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/drone/drone-plugin-go/plugin"
	"golang.org/x/crypto/ssh"
)

// Params stores the git clone parameters used to
// configure and customzie the git clone behavior.
type Params struct {
	Commands []string `json:"commands"`
	Login    string   `json:"user"`
	Port     int      `json:"port"`
	Host     string   `json:"host"`
}

func main() {
	v := new(Params)
	w := new(plugin.Workspace)
	plugin.Param("workspace", w)
	plugin.Param("vargs", &v)
	plugin.MustParse()

	err := run(w.Keys, v)
	if err != nil {
		os.Exit(1)
	}
}

func run(keys *plugin.Keypair, params *Params) error {

	// if no username is provided assume root
	if len(params.Login) == 0 {
		params.Login = "root"
	}

	// if no username is provided assume root
	if params.Port == 0 {
		params.Port = 22
	}

	// join the host and port if necessary
	host := net.JoinHostPort(
		params.Host,
		strconv.Itoa(params.Port),
	)

	signer, err := ssh.ParsePrivateKey([]byte(keys.Private))
	if err != nil {
		return fmt.Errorf("Error parsing private key. %s.", err)
	}

	config := &ssh.ClientConfig{
		User: params.Login,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return fmt.Errorf("Error dialing server. %s.", err)
	}

	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Error starting ssh session. %s.", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	return session.Run(strings.Join(params.Commands, "\n"))
}
