# Docker image for the Drone build runner
#
#     CGO_ENABLED=0 go build -a -tags netgo
#     docker build --rm=true -t plugins/drone-ssh .

FROM gliderlabs/alpine:3.1
RUN apk add --update \
	python \
	py-pip \
	&& pip install awscli
ADD drone-ssh /bin/
ENTRYPOINT ["/bin/drone-ssh"]