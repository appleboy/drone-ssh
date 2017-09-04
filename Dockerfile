FROM alpine:3.4

RUN apk update && \
  apk add -U --no-cache \
  ca-certificates \
  openssh-client && \
  rm -rf /var/cache/apk/*

LABEL org.label-schema.version=latest
LABEL org.label-schema.vcs-url="https://github.com/appleboy/drone-ssh.git"
LABEL org.label-schema.name="drone-ssh"
LABEL org.label-schema.vendor="Bo-Yi Wu"
LABEL org.label-schema.schema-version="1.0"

ADD drone-ssh /bin/
ENTRYPOINT ["/bin/drone-ssh"]
