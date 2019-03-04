FROM plugins/base:linux-arm64

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="Drone SSH" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

RUN apk add --no-cache ca-certificates && \
  rm -rf /var/cache/apk/*

ADD release/linux/arm64/drone-ssh /bin/
ENTRYPOINT ["/bin/drone-ssh"]
