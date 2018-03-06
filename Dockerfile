# Define build image which compiles binary
FROM golang:1.9.2-alpine as build

WORKDIR /go/src/github.com/appleboy/drone-ssh

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo \
 && cp drone-ssh /bin

# Define final image which consumes final artifact
FROM scratch

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/appleboy/drone-ssh.git"
LABEL org.label-schema.name="drone-ssh"
LABEL org.label-schema.vendor="Bo-Yi Wu"
LABEL org.label-schema.schema-version="1.0"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build /bin/drone-ssh /bin/

ENTRYPOINT ["/bin/drone-ssh"]
