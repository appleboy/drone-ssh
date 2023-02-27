FROM alpine:3.17

ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="Bo-Yi Wu <appleboy.tw@gmail.com>" \
  org.label-schema.name="SSH Plugin" \
  org.label-schema.vendor="Bo-Yi Wu" \
  org.label-schema.schema-version="1.0"

RUN apk add --no-cache ca-certificates && \
  rm -rf /var/cache/apk/*

COPY release/${TARGETOS}/${TARGETARCH}/drone-ssh /bin/

ENTRYPOINT ["/bin/drone-ssh"]
