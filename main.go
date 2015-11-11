package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/drone/drone-plugin-go/plugin"
	"golang.org/x/crypto/ssh"
)

// Params stores the git clone parameters used to
// configure and customzie the git clone behavior.
type Params struct {
	Commands []string `json:"commands"`
	Login    string   `json:"user"`
	Port     int      `json:"port"`
	Host     StrSlice `json:"host"`
	Sleep    int      `json:"sleep"`
}

func main() {
	v := new(Params)
	w := new(plugin.Workspace)
	plugin.Param("workspace", w)
	plugin.Param("vargs", &v)
	plugin.MustParse()

	for i, host := range v.Host.Slice() {
		err := run(w.Keys, v, host)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if v.Sleep != 0 && i != v.Host.Len()-1 {
			fmt.Printf("$ sleep %d\n", v.Sleep)
			time.Sleep(time.Duration(v.Sleep) * time.Second)
		}
	}
}

func run(keys *plugin.Keypair, params *Params, host string) error {

	// if no username is provided assume root
	if len(params.Login) == 0 {
		params.Login = "root"
	}

	// if no port is provided use default
	if params.Port == 0 {
		params.Port = 22
	}

	// join the host and port if necessary
	addr := net.JoinHostPort(
		host,
		strconv.Itoa(params.Port),
	)

	// trace command used for debugging in the build logs
	fmt.Printf("$ ssh %s@%s -p %d\n", params.Login, addr, params.Port)

	signer, err := ssh.ParsePrivateKey([]byte(keys.Private))
	if err != nil {
		return fmt.Errorf("Error parsing private key. %s.", err)
	}

	config := &ssh.ClientConfig{
		User: params.Login,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}

	client, err := ssh.Dial("tcp", addr, config)
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
