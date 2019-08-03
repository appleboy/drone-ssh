# drone-ssh

![sshlog](images/ssh.png)

[![GitHub tag](https://img.shields.io/github/tag/appleboy/drone-ssh.svg)](https://github.com/appleboy/drone-ssh/releases)
[![GoDoc](https://godoc.org/github.com/appleboy/drone-ssh?status.svg)](https://godoc.org/github.com/appleboy/drone-ssh)
[![Build Status](https://cloud.drone.io/api/badges/appleboy/drone-ssh/status.svg)](https://cloud.drone.io/appleboy/drone-ssh)
[![codecov](https://codecov.io/gh/appleboy/drone-ssh/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-ssh)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-ssh)](https://goreportcard.com/report/github.com/appleboy/drone-ssh)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-ssh.svg)](https://hub.docker.com/r/appleboy/drone-ssh/)
[![micro badger](https://images.microbadger.com/badges/image/appleboy/drone-ssh.svg)](https://microbadger.com/images/appleboy/drone-ssh "Get your own image badge on microbadger.com")

Drone plugin to execute commands on a remote host through SSH. For the usage
information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/appleboy/drone-ssh/).

**Note: Please update your image config path to `appleboy/drone-ssh` for drone. `plugins/ssh` is no longer maintained.**

![demo](./images/demo2017.05.10.gif)

## Breaking changes

`v1.5.0`: change command timeout flag to `Duration`. See the following setting:

```diff
pipeline:
  scp:
    image: appleboy/drone-scp
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

## Build or Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/drone-ssh/releases). Support the following OS type.

* Windows amd64/386
* Linux arm/amd64/386
* Darwin amd64/386

With `Go` installed

```sh
go get -u -v github.com/appleboy/drone-ssh
```

or build the binary with the following command:

```sh
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go test -cover ./...

go build -v -a -tags netgo -o release/linux/amd64/drone-ssh .
```

## Docker

Build the docker image with the following commands:

```sh
make docker
```

## Usage

Execute from the working directory:

```sh
docker run --rm \
  -e PLUGIN_HOST=foo.com \
  -e PLUGIN_USERNAME=root \
  -e PLUGIN_KEY="$(cat ${HOME}/.ssh/id_rsa)" \
  -e PLUGIN_SCRIPT=whoami \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-ssh
```

## Mount key from file path

Please make sure that enable the `trusted` mode in project setting for [drone 0.8 version](https://0-8-0.docs.drone.io/).

![trusted mode](./images/trust.png)

Mount private key in `volumes` setting of `.drone.yml` config

```diff
pipeline:
  ssh:
    image: appleboy/drone-ssh
    host: xxxxx.com
    username: deploy
+   volumes:
+     - /root/drone_rsa:/root/ssh/drone_rsa
    key_path: /root/ssh/drone_rsa
    script:
      - echo "test ssh"
```

See the detail of [issue comment](https://github.com/appleboy/drone-ssh/issues/51#issuecomment-336732928).
