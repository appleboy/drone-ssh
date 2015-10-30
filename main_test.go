package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/drone/drone-plugin-go/plugin"
)

var (
	host = os.Getenv("TEST_SSH_HOST")
	user = os.Getenv("TEST_SSH_USER")
	key  = os.Getenv("TEST_SSH_KEY")
)

func TestRun(t *testing.T) {

	// only runs the test if a host server is provided
	if len(host) == 0 {
		t.Skipf("TEST_SSH_HOST not provided")
		return
	}

	out, err := ioutil.ReadFile(key)
	if err != nil {
		t.Errorf("Unable to read or find a test privte key. %s", err)
	}

	params := &Params{
		Commands: []string{"whoami", "time", "ps -ax"},
		Login:    user,
		Host:     StrSlice{[]string{host}},
	}

	keys := &plugin.Keypair{
		Private: string(out),
	}

	err = run(keys, params, host)
	if err != nil {
		t.Errorf("Unable to run SSH commands. %s.", err)
	}
}
