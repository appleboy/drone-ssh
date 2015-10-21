# Docker image for the Drone build runner
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-ssh .

FROM gliderlabs/alpine:3.1
ADD drone-ssh /bin/
ENTRYPOINT ["/bin/drone-ssh"]