# drone-ssh

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-ssh/status.svg)](http://beta.drone.io/drone-plugins/drone-ssh)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-ssh?status.svg)](http://godoc.org/github.com/drone-plugins/drone-ssh)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-ssh)](https://goreportcard.com/report/github.com/drone-plugins/drone-ssh)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)

Drone plugin to execute commands on a remote host through SSH. For the usage
information and a listing of the available options please take a look at
[the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/ssh .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-ssh' not found or does not exist..
```

## Usage

Execute from the working directory:

```sh
docker run --rm \
  -e PLUGIN_HOST=foo.com \
  -e PLUGIN_USER=root \
  -e PLUGIN_KEY="$(cat ${HOME}/.ssh/id_rsa)" \
  -e PLUGIN_SCRIPT=whoami \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/ssh
```
