FROM alpine:3.21

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="SSH Plugin" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

LABEL org.opencontainers.image.source=https://github.com/appleboy/drone-ssh
LABEL org.opencontainers.image.description="Execute commands on a remote host through SSH"
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache ca-certificates tzdata && \
  rm -rf /var/cache/apk/*

RUN addgroup \
  -S -g 1000 \
  deploy && \
  adduser \
  -S -H -D \
  -h /home/deploy \
  -s /bin/sh \
  -u 1000 \
  -G deploy \
  deploy

RUN mkdir -p /home/deploy && \
  chown deploy:deploy /home/deploy

# deploy:deploy
USER 1000:1000

COPY release/${TARGETOS}/${TARGETARCH}/drone-ssh /bin/

ENTRYPOINT ["/bin/drone-ssh"]
