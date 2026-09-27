package main

import (
	"bytes"
	"errors"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	easyssh "github.com/appleboy/easyssh-proxy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func TestRunDebugRedactsCredentials(t *testing.T) {
	output, err := os.CreateTemp(t.TempDir(), "debug-output")
	require.NoError(t, err)
	original := os.Stdout
	os.Stdout = output
	t.Cleanup(func() {
		os.Stdout = original
		_ = output.Close()
	})

	flags := flag.NewFlagSet("test", flag.ContinueOnError)
	flags.Bool("debug", true, "")
	flags.String("user", "debug-user", "")
	flags.Int("port", 2222, "")
	flags.String("proxy.host", "bastion.example.invalid", "")
	for _, name := range []string{
		"ssh-key", "password", "ssh-passphrase",
		"proxy.ssh-key", "proxy.password", "proxy.ssh-passphrase",
	} {
		flags.String(name, "secret-"+name, "")
	}
	ctx := cli.NewContext(cli.NewApp(), flags, nil)
	require.ErrorIs(t, run(ctx), errMissingHost)
	_, err = output.Seek(0, io.SeekStart)
	require.NoError(t, err)
	logged, err := io.ReadAll(output)
	require.NoError(t, err)
	assert.NotContains(t, string(logged), "secret-")
	assert.Equal(t, 6, strings.Count(string(logged), "[REDACTED]"))
	assert.Contains(t, string(logged), "debug-user")
	assert.Contains(t, string(logged), "2222")
	assert.Contains(t, string(logged), "bastion.example.invalid")
}

func TestConfigRedactedPreservesOriginal(t *testing.T) {
	original := Config{
		Key: "private-key", Password: "password", Passphrase: "passphrase",
		Host: []string{"example.invalid"}, Port: 2222, Username: "deploy",
		KeyPath: "/keys/deploy", Timeout: time.Second, Debug: true,
		Envs: []string{"TOKEN"}, Script: []string{"whoami"},
		Proxy: easyssh.DefaultConfig{
			Key: "proxy-key", Password: "proxy-password", Passphrase: "proxy-passphrase",
			Server: "bastion.example.invalid", User: "jump", Port: "2223",
		},
	}
	config := original
	got := config.redacted()
	assert.Equal(t, original, config, "debug rendering must not change connection credentials")
	expected := original
	expected.Key, expected.Password, expected.Passphrase = "[REDACTED]", "[REDACTED]", "[REDACTED]"
	expected.Proxy.Key, expected.Proxy.Password, expected.Proxy.Passphrase = "[REDACTED]", "[REDACTED]", "[REDACTED]"
	assert.Equal(t, expected, got, "non-sensitive settings must remain available for debugging")
	assert.Equal(
		t,
		Config{},
		(Config{}).redacted(),
		"unset credentials must remain distinguishable",
	)
}

func TestDebugRedactsEnvironment(t *testing.T) {
	t.Setenv("DRONE_SSH_TEST_SECRET", "sensitive-value'with-quotes")
	var output bytes.Buffer
	p := Plugin{
		Config: Config{
			Host:           []string{"127.0.0.1"},
			Password:       "test",
			Port:           0,
			Timeout:        100 * time.Millisecond,
			Debug:          true,
			Envs:           []string{"DRONE_SSH_TEST_SECRET"},
			EnvsFormat:     envsFormat,
			CommandTimeout: time.Second,
		},
		Writer: &output,
	}
	// Debug output is produced before this intentionally refused connection.
	require.Error(t, p.exec("127.0.0.1"))
	assert.Contains(t, output.String(), "DRONE_SSH_TEST_SECRET=[REDACTED]")
	assert.NotContains(t, output.String(), "sensitive-value")
}

func TestExecHostsWaitsForAllWorkers(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		p := Plugin{Config: Config{Host: []string{"fast", "slow", "success"}}}
		failure := errors.New("host failure")
		fastDone := make(chan struct{})
		release := make(chan struct{})
		result := make(chan error, 1)
		go func() {
			result <- p.execHosts(func(host string) error {
				switch host {
				case "fast":
					close(fastDone)
					return failure
				case "slow":
					<-release
					return errors.New("second failure")
				default:
					<-release
					return nil
				}
			})
		}()
		<-fastDone
		synctest.Wait()
		select {
		case <-result:
			t.Fatal("returned before all started hosts finished")
		default:
		}
		close(release)
		err := <-result
		require.ErrorIs(t, err, failure)
		assert.Contains(t, err.Error(), "fast")
		// synctest also requires every goroutine in the bubble to exit.
	})
}

func TestExecHostsSyncStopsOnFailure(t *testing.T) {
	p := Plugin{Config: Config{Host: []string{"first", "failed", "unstarted"}, Sync: true}}
	failure := errors.New("failure")
	var visited []string
	err := p.execHosts(func(host string) error {
		visited = append(visited, host)
		if host == "failed" {
			return failure
		}
		return nil
	})
	require.ErrorIs(t, err, failure)
	assert.Equal(t, []string{"first", "failed"}, visited)
}

func TestExecHostsSuccess(t *testing.T) {
	for _, syncMode := range []bool{false, true} {
		p := Plugin{Config: Config{Host: []string{"one", "two"}, Sync: syncMode}}
		visited := make(chan string, 2)
		require.NoError(t, p.execHosts(func(host string) error {
			visited <- host
			return nil
		}))
		close(visited)
		var hosts []string
		for host := range visited {
			hosts = append(hosts, host)
		}
		assert.ElementsMatch(t, p.Config.Host, hosts)
	}
}

func TestReadStreamClosedChannelsBlockUntilDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		stdout := make(chan string)
		stderr := make(chan string)
		streamErrors := make(chan error)
		done := make(chan bool)
		result := make(chan error, 1)
		var output bytes.Buffer
		p := Plugin{Config: Config{Host: []string{"host"}}, Writer: &output}
		go func() {
			result <- p.readStream("host", stdout, stderr, done, streamErrors)
		}()
		stdout <- "stdout line"
		stderr <- "stderr line"
		close(stdout)
		close(stderr)
		close(streamErrors)
		// Wait only returns when the reader blocks. A closed-channel busy loop
		// cannot become durably blocked and fails the test's timeout.
		synctest.Wait()
		select {
		case <-result:
			t.Fatal("returned before the completion signal")
		default:
		}
		done <- true
		require.NoError(t, <-result)
		assert.Equal(t, "stdout line\nstderr line\n", output.String())
	})
}

func TestReadStreamErrorsAndTimeout(t *testing.T) {
	failure := errors.New("remote command failed")
	for _, tc := range []struct {
		name      string
		completed bool
		streamErr error
		want      error
	}{
		{"exit error", true, failure, failure},
		{"timeout with error", false, failure, failure},
		{"timeout without error", false, nil, errCommandTimeOut},
	} {
		t.Run(tc.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				streamErrors := make(chan error)
				done := make(chan bool)
				result := make(chan error, 1)
				go func() {
					result <- (Plugin{}).readStream("host", nil, nil, done, streamErrors)
				}()
				streamErrors <- tc.streamErr
				close(streamErrors)
				synctest.Wait()
				done <- tc.completed
				require.ErrorIs(t, <-result, tc.want)
			})
		})
	}
}

func TestExecMultipleConnectionFailures(t *testing.T) {
	p := Plugin{Config: Config{
		Host:     []string{"127.0.0.1", "127.0.0.1"},
		Password: "test",
		Timeout:  100 * time.Millisecond,
	}, Writer: io.Discard}
	err := p.Exec()
	require.Error(t, err)
	assert.True(t, strings.HasPrefix(err.Error(), "127.0.0.1: "))
}
