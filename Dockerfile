FROM appleboy/golang-testing AS build-env
ADD . /go/src/github.com/appleboy/drone-ssh
RUN cd /go/src/github.com/appleboy/drone-ssh && make static_build

FROM alpine:3.4

RUN apk update && \
  apk add \
    ca-certificates \
    openssh-client && \
  rm -rf /var/cache/apk/*

COPY --from=build-env /go/src/github.com/appleboy/drone-ssh/drone-ssh /bin
ENTRYPOINT ["/bin/drone-ssh"]
