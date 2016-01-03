# Docker image for the Drone Swift plugin
#
#     cd $GOPATH/src/github.com/drone-plugins/drone-ssh
#     make deps build docker

FROM alpine:3.3

RUN apk update && \
  apk add \
    ca-certificates && \
  rm -rf /var/cache/apk/*

ADD drone-ssh /bin/
ENTRYPOINT ["/bin/drone-ssh"]
