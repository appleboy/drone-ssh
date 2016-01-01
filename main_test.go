package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/drone/drone-go/drone"
)

var (
	host = os.Getenv("TEST_SSH_HOST")
	user = os.Getenv("TEST_SSH_USER")
	key  = os.Getenv("TEST_SSH_KEY")
)

func TestRun(t *testing.T) {
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
		Host: drone.NewStringSlice(
			[]string{
				host,
			},
		),
	}

	keys := &drone.Key{
		Private: string(out),
	}

	err = run(keys, params, host)

	if err != nil {
		t.Errorf("Unable to run SSH commands. %s.", err)
	}
}
