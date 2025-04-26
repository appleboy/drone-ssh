# drone-ssh

> **English** | [繁體中文](./README.zh-tw.md) | [简体中文](./README.zh-cn.md)

![sshlog](images/ssh.png)

[![GitHub tag](https://img.shields.io/github/tag/appleboy/drone-ssh.svg)](https://github.com/appleboy/drone-ssh/releases)
[![GoDoc](https://godoc.org/github.com/appleboy/drone-ssh?status.svg)](https://godoc.org/github.com/appleboy/drone-ssh)
[![Lint and Testing](https://github.com/appleboy/drone-ssh/actions/workflows/testing.yml/badge.svg?branch=master)](https://github.com/appleboy/drone-ssh/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/appleboy/drone-ssh/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-ssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-ssh)](https://goreportcard.com/report/github.com/appleboy/drone-ssh)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-ssh.svg)](https://hub.docker.com/r/appleboy/drone-ssh/)

A Drone plugin for executing commands on remote hosts via SSH. For usage instructions and a list of available options, please refer to [the documentation](http://plugins.drone.io/appleboy/drone-ssh/).

**Note: Please update your Drone image config path to `appleboy/drone-ssh`. The `plugins/ssh` image is no longer maintained.**

![demo](./images/demo2017.05.10.gif)

## Table of Contents

- [drone-ssh](#drone-ssh)
  - [Table of Contents](#table-of-contents)
  - [Breaking Changes](#breaking-changes)
  - [Build or Download a Binary](#build-or-download-a-binary)
  - [Docker](#docker)
  - [Usage](#usage)
  - [Mount Key from File Path](#mount-key-from-file-path)
  - [Configuration](#configuration)

## Breaking Changes

As of `v1.5.0`, the command timeout flag has changed to use the `Duration` format. See the following example:

```diff
pipeline:
  scp:
    image: ghcr.io/appleboy/drone-ssh
    settings:
      host:
        - example1.com
        - example2.com
      username: ubuntu
      password:
        from_secret: ssh_password
      port: 22
-     command_timeout: 120
+     command_timeout: 2m
      script:
        - echo "Hello World"
```

## Build or Download a Binary

Pre-compiled binaries are available on the [releases page](https://github.com/appleboy/drone-ssh/releases), supporting the following operating systems:

- Windows amd64/386
- Linux arm/amd64/386
- macOS (Darwin) amd64/386

If you have `Go` installed:

```sh
go install github.com/appleboy/drone-ssh@latest
```

Or build the binary manually with the following commands:

```sh
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go test -cover ./...

go build -v -a -tags netgo -o release/linux/amd64/drone-ssh .
```

## Docker

Build the Docker image with the following command:

```sh
make docker
```

## Usage

Run from your working directory:

```sh
docker run --rm \
  -e PLUGIN_HOST=foo.com \
  -e PLUGIN_USERNAME=root \
  -e PLUGIN_KEY="$(cat ${HOME}/.ssh/id_rsa)" \
  -e PLUGIN_SCRIPT=whoami \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  ghcr.io/appleboy/drone-ssh
```

## Mount Key from File Path

Make sure to enable `trusted` mode in your project settings (for [Drone 0.8 version](https://0-8-0.docs.drone.io/)).

![trusted mode](./images/trust.png)

Mount the private key in the `volumes` section of your `.drone.yml` config:

```diff
pipeline:
  ssh:
    image: ghcr.io/appleboy/drone-ssh
    host: xxxxx.com
    username: deploy
+   volumes:
+     - /root/drone_rsa:/root/ssh/drone_rsa
    key_path: /root/ssh/drone_rsa
    script:
      - echo "test ssh"
```

See details in [this issue comment](https://github.com/appleboy/drone-ssh/issues/51#issuecomment-336732928).

## Configuration

See [DOCS.md](./DOCS.md) for examples and full configuration options.

Configuration options are loaded from multiple sources:

0. Hardcoded drone-ssh defaults. See [main.go CLI Flags](https://github.com/appleboy/drone-ssh/blob/6d9d6acc6aef1f9166118c6ba8bd214d3a582bdb/main.go#L39) for more information.
1. From a dotenv file at a path specified by the `PLUGIN_ENV_FILE` environment variable.
2. From your `.drone.yml` Drone configuration.

Later sources override earlier ones. For example, if `PORT` is set in an `.env` file committed in the repository or created by previous test steps, it will override the default set in `main.go`.
