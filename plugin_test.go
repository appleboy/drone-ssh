package main

import (
	"os"
	"testing"

	"bytes"
	"github.com/appleboy/easyssh-proxy"
	"github.com/stretchr/testify/assert"
	"strings"
)

func TestMissingHostOrUser(t *testing.T) {
	plugin := Plugin{}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, missingHostOrUser, err.Error())
}

func TestMissingKeyOrPassword(t *testing.T) {
	plugin := Plugin{
		Config{
			Host:     []string{"localhost"},
			UserName: "ubuntu",
		},
		os.Stdout,
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, missingPasswordOrKey, err.Error())
}

func TestSetPasswordAndKey(t *testing.T) {
	plugin := Plugin{
		Config{
			Host:     []string{"localhost"},
			UserName: "ubuntu",
			Password: "1234",
			Key:      "1234",
		},
		os.Stdout,
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, setPasswordandKey, err.Error())
}

func TestIncorrectPassword(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			UserName:       "drone-scp",
			Port:           22,
			Password:       "123456",
			Script:         []string{"whoami"},
			CommandTimeout: 60,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestSSHScriptFromRawKey(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			UserName:       "drone-scp",
			Port:           22,
			CommandTimeout: 60,
			Key: `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA4e2D/qPN08pzTac+a8ZmlP1ziJOXk45CynMPtva0rtK/RB26
VbfAF0hIJji7ltvnYnqCU9oFfvEM33cTn7T96+od8ib/Vz25YU8ZbstqtIskPuwC
bv3K0mAHgsviJyRD7yM+QKTbBQEgbGuW6gtbMKhiYfiIB4Dyj7AdS/fk3v26wDgz
7SHI5OBqu9bv1KhxQYdFEnU3PAtAqeccgzNpbH3eYLyGzuUxEIJlhpZ/uU2G9ppj
/cSrONVPiI8Ahi4RrlZjmP5l57/sq1ClGulyLpFcMw68kP5FikyqHpHJHRBNgU57
1y0Ph33SjBbs0haCIAcmreWEhGe+/OXnJe6VUQIDAQABAoIBAH97emORIm9DaVSD
7mD6DqA7c5m5Tmpgd6eszU08YC/Vkz9oVuBPUwDQNIX8tT0m0KVs42VVPIyoj874
bgZMJoucC1G8V5Bur9AMxhkShx9g9A7dNXJTmsKilRpk2TOk7wBdLp9jZoKoZBdJ
jlp6FfaazQjjKD6zsCsMATwAoRCBpBNsmT6QDN0n0bIgY0tE6YGQaDdka0dAv68G
R0VZrcJ9voT6+f+rgJLoojn2DAu6iXaM99Gv8FK91YCymbQlXXgrk6CyS0IHexN7
V7a3k767KnRbrkqd3o6JyNun/CrUjQwHs1IQH34tvkWScbseRaFehcAm6mLT93RP
muauvMECgYEA9AXGtfDMse0FhvDPZx4mx8x+vcfsLvDHcDLkf/lbyPpu97C27b/z
ia07bu5TAXesUZrWZtKA5KeRE5doQSdTOv1N28BEr8ZwzDJwfn0DPUYUOxsN2iIy
MheO5A45Ko7bjKJVkZ61Mb1UxtqCTF9mqu9R3PBdJGthWOd+HUvF460CgYEA7QRf
Z8+vpGA+eSuu29e0xgRKnRzed5zXYpcI4aERc3JzBgO4Z0er9G8l66OWVGdMfpe6
CBajC5ToIiT8zqoYxXwqJgN+glir4gJe3mm8J703QfArZiQrdk0NTi5bY7+vLLG/
knTrtpdsKih6r3kjhuPPaAsIwmMxIydFvATKjLUCgYEAh/y4EihRSk5WKC8GxeZt
oiZ58vT4z+fqnMIfyJmD5up48JuQNcokw/LADj/ODiFM7GUnWkGxBrvDA3H67WQm
49bJjs8E+BfUQFdTjYnJRlpJZ+7Zt1gbNQMf5ENw5CCchTDqEq6pN0DVf8PBnSIF
KvkXW9KvdV5J76uCAn15mDkCgYA1y8dHzbjlCz9Cy2pt1aDfTPwOew33gi7U3skS
RTerx29aDyAcuQTLfyrROBkX4TZYiWGdEl5Bc7PYhCKpWawzrsH2TNa7CRtCOh2E
R+V/84+GNNf04ALJYCXD9/ugQVKmR1XfDRCvKeFQFE38Y/dvV2etCswbKt5tRy2p
xkCe/QKBgQCkLqafD4S20YHf6WTp3jp/4H/qEy2X2a8gdVVBi1uKkGDXr0n+AoVU
ib4KbP5ovZlrjL++akMQ7V2fHzuQIFWnCkDA5c2ZAqzlM+ZN+HRG7gWur7Bt4XH1
7XC9wlRna4b3Ln8ew3q1ZcBjXwD4ppbTlmwAfQIaZTGJUgQbdsO9YA==
-----END RSA PRIVATE KEY-----
`,
			Script: []string{"whoami"},
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func TestSSHScriptFromKeyFile(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost", "127.0.0.1"},
			UserName:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami", "ls -al"},
			CommandTimeout: 60,
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func TestStreamFromSSHCommand(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost", "127.0.0.1"},
			UserName:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami", "for i in {1..5}; do echo ${i}; sleep 1; done", "echo 'done'"},
			CommandTimeout: 60,
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func TestSSHScriptWithError(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost", "127.0.0.1"},
			UserName:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"exit 1"},
			CommandTimeout: 60,
		},
	}

	err := plugin.Exec()
	// Process exited with status 1
	assert.NotNil(t, err)
}

func TestSSHCommandTimeOut(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			UserName:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"sleep 5"},
			CommandTimeout: 1,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestProxyCommand(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			UserName:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami"},
			CommandTimeout: 1,
			Proxy: easyssh.DefaultConfig{
				Server:  "localhost",
				User:    "drone-scp",
				Port:    "22",
				KeyPath: "./tests/.ssh/id_rsa",
			},
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func TestSSHCommandError(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			UserName:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"mkdir a", "mkdir a"},
			CommandTimeout: 60,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestSSHCommandExitCodeError(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			UserName: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"set -e",
				"echo 1",
				"mkdir a",
				"mkdir a",
				"echo 2",
			},
			CommandTimeout: 60,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestSetENV(t *testing.T) {
	os.Setenv("FOO", `'  1)  '`)
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			UserName:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Secrets:        []string{"FOO"},
			Envs:           []string{"foo"},
			Debug:          true,
			Script:         []string{"whoami; echo $FOO"},
			CommandTimeout: 1,
			Proxy: easyssh.DefaultConfig{
				Server:  "localhost",
				User:    "drone-scp",
				Port:    "22",
				KeyPath: "./tests/.ssh/id_rsa",
			},
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func TestSyncMode(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost", "127.0.0.1"},
			UserName:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami", "for i in {1..3}; do echo ${i}; sleep 1; done", "echo 'done'"},
			CommandTimeout: 60,
			Sync:           true,
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func Test_escapeArg(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "escape nothing",
			args: args{
				arg: "Hi I am appleboy",
			},
			want: `'Hi I am appleboy'`,
		},
		{
			name: "escape single quote",
			args: args{
				arg: "Hi I am 'appleboy'",
			},
			want: `'Hi I am '\''appleboy'\'''`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeArg(tt.args.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommandOutput(t *testing.T) {
	var (
		buffer   bytes.Buffer
		expected = `
			localhost: ======CMD======
			localhost: pwd
			whoami
			uname
			localhost: ======END======
			localhost: out: /home/drone-scp
			localhost: out: drone-scp
			localhost: out: Linux
			127.0.0.1: ======CMD======
			127.0.0.1: pwd
			whoami
			uname
			127.0.0.1: ======END======
			127.0.0.1: out: /home/drone-scp
			127.0.0.1: out: drone-scp
			127.0.0.1: out: Linux
		`
	)

	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost", "127.0.0.1"},
			UserName: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"pwd",
				"whoami",
				"uname",
			},
			CommandTimeout: 60,
			Sync:           true,
		},
		Writer: &buffer,
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	assert.Equal(t, unindent(expected), unindent(buffer.String()))
}

func TestEnvOutput(t *testing.T) {
	var (
		buffer   bytes.Buffer
		expected = `
			======CMD======
			echo "[${ENV_1}]"
			echo "[${ENV_2}]"
			echo "[${ENV_3}]"
			echo "[${ENV_4}]"
			echo "[${ENV_5}]"
			echo "[${ENV_6}]"
			echo "[${ENV_7}]"
			======END======
			======ENV======
			ENV_1='test'
			ENV_2='test test'
			ENV_3='test '
			ENV_4='  test  test  '
			ENV_5='test'\'''
			ENV_6='test"'
			ENV_7='test,!#;?.@$~'\''"'
			======END======
			out: [test]
			out: [test test]
			out: [test ]
			out: [  test  test  ]
			out: [test']
			out: [test"]
			out: [test,!#;?.@$~'"]
		`
	)

	os.Setenv("ENV_1", `test`)
	os.Setenv("ENV_2", `test test`)
	os.Setenv("ENV_3", `test `)
	os.Setenv("ENV_4", `  test  test  `)
	os.Setenv("ENV_5", `test'`)
	os.Setenv("ENV_6", `test"`)
	os.Setenv("ENV_7", `test,!#;?.@$~'"`)

	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			UserName: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Envs:     []string{"env_1", "env_2", "env_3", "env_4", "env_5", "env_6", "env_7"},
			Debug:    true,
			Script: []string{
				`echo "[${ENV_1}]"`,
				`echo "[${ENV_2}]"`,
				`echo "[${ENV_3}]"`,
				`echo "[${ENV_4}]"`,
				`echo "[${ENV_5}]"`,
				`echo "[${ENV_6}]"`,
				`echo "[${ENV_7}]"`,
			},
			CommandTimeout: 10,
			Proxy: easyssh.DefaultConfig{
				Server:  "localhost",
				User:    "drone-scp",
				Port:    "22",
				KeyPath: "./tests/.ssh/id_rsa",
			},
		},
		Writer: &buffer,
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	assert.Equal(t, unindent(expected), unindent(buffer.String()))
}

func unindent(text string) string {
	return strings.TrimSpace(strings.Replace(text, "\t", "", -1))
}
