package main

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestMissingHostOrUser(t *testing.T) {
	plugin := Plugin{}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, errMissingHost, err)
}

func TestMissingKeyOrPassword(t *testing.T) {
	plugin := Plugin{
		Config{
			Host:     []string{"localhost"},
			Username: "ubuntu",
		},
		os.Stdout,
	}

	err := plugin.Exec()

	assert.NotNil(t, err)
	assert.Equal(t, errMissingPasswordOrKey, err)
}

func TestIncorrectPassword(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           22,
			Password:       "123456",
			Script:         []string{"whoami"},
			CommandTimeout: 60 * time.Second,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestSSHScriptFromRawKey(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           22,
			CommandTimeout: 60 * time.Second,
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
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami", "ls -al"},
			CommandTimeout: 60 * time.Second,
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func TestSSHIPv4Only(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost", "127.0.0.1"},
			Username:       "drone-scp",
			Port:           22,
			Protocol:       easyssh.PROTOCOL_TCP4,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami", "ls -al"},
			CommandTimeout: 60 * time.Second,
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func TestSSHIPv6OnlyError(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"127.0.0.1"},
			Username:       "drone-scp",
			Port:           22,
			Protocol:       easyssh.PROTOCOL_TCP6,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami", "ls -al"},
			CommandTimeout: 60 * time.Second,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestStreamFromSSHCommand(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost", "127.0.0.1"},
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami", "for i in {1..5}; do echo ${i}; sleep 1; done", "echo 'done'"},
			CommandTimeout: 60 * time.Second,
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)
}

func TestSSHScriptWithError(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost", "127.0.0.1"},
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"exit 1"},
			CommandTimeout: 60 * time.Second,
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
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"sleep 5"},
			CommandTimeout: 1 * time.Second,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestProxyCommand(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami"},
			CommandTimeout: 1 * time.Second,
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
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"mkdir a", "mkdir a"},
			CommandTimeout: 60 * time.Second,
		},
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func TestSSHCommandExitCodeError(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"set -e",
				"echo 1",
				"mkdir a",
				"mkdir a",
				"echo 2",
			},
			CommandTimeout: 60 * time.Second,
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
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Envs:           []string{"foo"},
			Debug:          true,
			Script:         []string{"whoami; echo $FOO"},
			CommandTimeout: 1 * time.Second,
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

func TestSetExistingENV(t *testing.T) {
	os.Setenv("FOO", "Value for foo")
	os.Setenv("BAR", "")
	plugin := Plugin{
		Config: Config{
			Host:           []string{"localhost"},
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Envs:           []string{"foo", "bar", "baz"},
			Debug:          true,
			Script:         []string{"export FOO", "export BAR", "export BAZ", "env | grep -q '^FOO=Value for foo$'", "env | grep -q '^BAR=$'", "if env | grep -q BAZ; then false; else true; fi"},
			CommandTimeout: 1 * time.Second,
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
			Username:       "drone-scp",
			Port:           22,
			KeyPath:        "./tests/.ssh/id_rsa",
			Script:         []string{"whoami", "for i in {1..3}; do echo ${i}; sleep 1; done", "echo 'done'"},
			CommandTimeout: 60 * time.Second,
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
			Username: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"pwd",
				"whoami",
				"uname",
			},
			CommandTimeout: 60 * time.Second,
			Sync:           true,
		},
		Writer: &buffer,
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	assert.Equal(t, unindent(expected), unindent(buffer.String()))
}

func TestWrongFingerprint(t *testing.T) {
	var buffer bytes.Buffer

	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"whoami",
			},
			Fingerprint: "wrong",
		},
		Writer: &buffer,
	}

	err := plugin.Exec()
	assert.NotNil(t, err)
}

func getHostPublicKeyFile(keypath string) (ssh.PublicKey, error) {
	var pubkey ssh.PublicKey
	var err error
	buf, err := os.ReadFile(keypath)
	if err != nil {
		return nil, err
	}

	pubkey, _, _, _, err = ssh.ParseAuthorizedKey(buf)

	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

func TestFingerprint(t *testing.T) {
	var (
		buffer   bytes.Buffer
		expected = `
			======CMD======
			whoami
			======END======
			out: drone-scp
		`
	)

	hostKey, err := getHostPublicKeyFile("/etc/ssh/ssh_host_rsa_key.pub")
	assert.NoError(t, err)

	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"whoami",
			},
			Fingerprint:    ssh.FingerprintSHA256(hostKey),
			CommandTimeout: 10 * time.Second,
		},
		Writer: &buffer,
	}

	err = plugin.Exec()
	assert.Nil(t, err)
	assert.Equal(t, unindent(expected), unindent(buffer.String()))
}

func TestScriptStopWithMultipleHostAndSyncMode(t *testing.T) {
	var (
		buffer   bytes.Buffer
		expected = `
			======CMD======
			mkdir a/b/c
			mkdir d/e/f
			======END======
			err: mkdir: can't create directory 'a/b/c': No such file or directory
		`
	)

	plugin := Plugin{
		Config: Config{
			Host:     []string{"", "localhost"},
			Username: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"mkdir a/b/c",
				"mkdir d/e/f",
			},
			CommandTimeout: 10 * time.Second,
			ScriptStop:     true,
			Sync:           true,
		},
		Writer: &buffer,
	}

	err := plugin.Exec()
	assert.NotNil(t, err)

	assert.Equal(t, unindent(expected), unindent(buffer.String()))
}

func TestScriptStop(t *testing.T) {
	var (
		buffer   bytes.Buffer
		expected = `
			======CMD======
			mkdir a/b/c
			mkdir d/e/f
			======END======
			err: mkdir: can't create directory 'a/b/c': No such file or directory
		`
	)

	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"mkdir a/b/c",
				"mkdir d/e/f",
			},
			CommandTimeout: 10 * time.Second,
			ScriptStop:     true,
		},
		Writer: &buffer,
	}

	err := plugin.Exec()
	assert.NotNil(t, err)

	assert.Equal(t, unindent(expected), unindent(buffer.String()))
}

func TestNoneScriptStop(t *testing.T) {
	var (
		buffer   bytes.Buffer
		expected = `
			======CMD======
			mkdir a/b/c
			mkdir d/e/f
			======END======
			err: mkdir: can't create directory 'a/b/c': No such file or directory
			err: mkdir: can't create directory 'd/e/f': No such file or directory
		`
	)

	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"mkdir a/b/c",
				"mkdir d/e/f",
			},
			CommandTimeout: 10 * time.Second,
		},
		Writer: &buffer,
	}

	err := plugin.Exec()
	assert.NotNil(t, err)

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
			export ENV_1='test'
			export ENV_2='test test'
			export ENV_3='test '
			export ENV_4='  test  test  '
			export ENV_5='test'\'''
			export ENV_6='test"'
			export ENV_7='test,!#;?.@$~'\''"'
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
			Host:       []string{"localhost"},
			Username:   "drone-scp",
			Port:       22,
			KeyPath:    "./tests/.ssh/test",
			Passphrase: "1234",
			Envs:       []string{"env_1", "env_2", "env_3", "env_4", "env_5", "env_6", "env_7"},
			Debug:      true,
			Script: []string{
				`echo "[${ENV_1}]"`,
				`echo "[${ENV_2}]"`,
				`echo "[${ENV_3}]"`,
				`echo "[${ENV_4}]"`,
				`echo "[${ENV_5}]"`,
				`echo "[${ENV_6}]"`,
				`echo "[${ENV_7}]"`,
			},
			CommandTimeout: 10 * time.Second,
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

func TestPlugin_scriptCommands(t *testing.T) {
	type fields struct {
		Config Config
		Writer io.Writer
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "normal testing",
			fields: fields{
				Config: Config{
					Script: []string{"mkdir a", "mkdir b"},
				},
			},
			want: []string{"mkdir a", "mkdir b"},
		},
		{
			name: "script stop",
			fields: fields{
				Config: Config{
					Script:     []string{"mkdir a", "mkdir b"},
					ScriptStop: true,
				},
			},
			want: []string{"mkdir a", "DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;", "mkdir b", "DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;"},
		},
		{
			name: "normal testing 2",
			fields: fields{
				Config: Config{
					Script:     []string{"mkdir a\nmkdir c", "mkdir b"},
					ScriptStop: true,
				},
			},
			want: []string{"mkdir a", "DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;", "mkdir c", "DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;", "mkdir b", "DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;"},
		},
		// See: https://github.com/appleboy/ssh-action/issues/75#issuecomment-668314271
		{
			name: "Multiline SSH commands interpreted as single lines",
			fields: fields{
				Config: Config{
					Script:     []string{"ls \\ ", "-lah", "mkdir a"},
					ScriptStop: true,
				},
			},
			want: []string{"ls \\", "-lah", "DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;", "mkdir a", "DRONE_SSH_PREV_COMMAND_EXIT_CODE=$? ; if [ $DRONE_SSH_PREV_COMMAND_EXIT_CODE -ne 0 ]; then exit $DRONE_SSH_PREV_COMMAND_EXIT_CODE; fi;"},
		},
		{
			name: "trim space",
			fields: fields{
				Config: Config{
					Script:     []string{"mkdir a", "mkdir b", "\t", " "},
					ScriptStop: false,
				},
			},
			want: []string{"mkdir a", "mkdir b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Plugin{
				Config: tt.fields.Config,
				Writer: tt.fields.Writer,
			}
			if got := p.scriptCommands(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Plugin.scriptCommands() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestUseInsecureCipher(t *testing.T) {
	var (
		buffer   bytes.Buffer
		expected = `
			======CMD======
			mkdir a/b/c
			mkdir d/e/f
			======END======
			err: mkdir: can't create directory 'a/b/c': No such file or directory
			err: mkdir: can't create directory 'd/e/f': No such file or directory
		`
	)

	plugin := Plugin{
		Config: Config{
			Host:     []string{"localhost"},
			Username: "drone-scp",
			Port:     22,
			KeyPath:  "./tests/.ssh/id_rsa",
			Script: []string{
				"mkdir a/b/c",
				"mkdir d/e/f",
			},
			CommandTimeout:    10 * time.Second,
			UseInsecureCipher: true,
		},
		Writer: &buffer,
	}

	err := plugin.Exec()
	assert.NotNil(t, err)

	assert.Equal(t, unindent(expected), unindent(buffer.String()))
}

func TestPlugin_hostPort(t *testing.T) {
	type fields struct {
		Config Config
		Writer io.Writer
	}
	type args struct {
		h string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantHost string
		wantPort string
	}{
		{
			name: "default host and port",
			fields: fields{
				Config: Config{
					Port: 22,
				},
			},
			args: args{
				h: "localhost",
			},
			wantHost: "localhost",
			wantPort: "22",
		},
		{
			name: "different port",
			fields: fields{
				Config: Config{
					Port: 22,
				},
			},
			args: args{
				h: "localhost:443",
			},
			wantHost: "localhost",
			wantPort: "443",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Plugin{
				Config: tt.fields.Config,
				Writer: tt.fields.Writer,
			}
			gotHost, gotPort := p.hostPort(tt.args.h)
			if gotHost != tt.wantHost {
				t.Errorf("Plugin.hostPort() gotHost = %v, want %v", gotHost, tt.wantHost)
			}
			if gotPort != tt.wantPort {
				t.Errorf("Plugin.hostPort() gotPort = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

func TestFindEnvs(t *testing.T) {
	testEnvs := []string{
		"INPUT_FOO",
		"INPUT_BAR",
		"NO_PREFIX",
		"INPUT_FOOBAR",
	}

	origEnviron := os.Environ()
	os.Clearenv()
	for _, env := range testEnvs {
		os.Setenv(env, "dummyValue")
	}

	defer func() {
		os.Clearenv()
		for _, env := range origEnviron {
			pair := strings.SplitN(env, "=", 2)
			os.Setenv(pair[0], pair[1])
		}
	}()

	t.Run("Find single prefix", func(t *testing.T) {
		expected := []string{"INPUT_FOO", "INPUT_BAR", "INPUT_FOOBAR"}
		result := findEnvs("INPUT_")
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	t.Run("Find multiple prefixes", func(t *testing.T) {
		expected := []string{"INPUT_FOO", "INPUT_BAR", "NO_PREFIX", "INPUT_FOOBAR"}
		result := findEnvs("INPUT_", "NO_PREFIX")
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})

	t.Run("Find non-existing prefix", func(t *testing.T) {
		expected := []string{}
		result := findEnvs("NON_EXISTING_")
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, but got %v", expected, result)
		}
	})
}

func TestAllEnvs(t *testing.T) {
	var (
		buffer   bytes.Buffer
		expected = `
			out: [foobar]
			out: [foobar]
			out: [foobar]
		`
	)

	os.Setenv("INPUT_1", `foobar`)
	os.Setenv("GITHUB_2", `foobar`)
	os.Setenv("PLUGIN_3", `foobar`)

	plugin := Plugin{
		Config: Config{
			Host:       []string{"localhost"},
			Username:   "drone-scp",
			Port:       22,
			KeyPath:    "./tests/.ssh/test",
			Passphrase: "1234",
			AllEnvs:    true,
			Script: []string{
				`echo "[${INPUT_1}]"`,
				`echo "[${GITHUB_2}]"`,
				`echo "[${PLUGIN_3}]"`,
			},
			CommandTimeout: 10 * time.Second,
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
