# drone-ssh

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-ssh/status.svg)](http://beta.drone.io/drone-plugins/drone-ssh)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-ssh/coverage.svg)](https://aircover.co/drone-plugins/drone-ssh)
[![](https://badge.imagelayers.io/plugins/drone-ssh:latest.svg)](https://imagelayers.io/?images=plugins/drone-ssh:latest 'Get your own badge on imagelayers.io')

Drone plugin to execute commands on a remote host through SSH. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
export GO15VENDOREXPERIMENT=1
go build
go test
```

## Docker

Build the docker image with the following commands:

```
export GO15VENDOREXPERIMENT=1
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo
```

Please note incorrectly building the image for the correct x64 linux and with GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-ssh' not found or does not exist..
```

## Usage

Execute a single remote command

```sh
docker run --rm \
  -e PLUGIN_HOST=foo.com \
  -e PLUGIN_USER=root \
  -e PLUGIN_KEY="$(cat ${HOME}/.ssh/id_rsa)" \
  -e PLUGIN_COMMANDS=whoami \
  -v $(pwd)/$(pwd) \
  -w $(pwd) \
  plugins/ssh
```
