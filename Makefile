.PHONY: clean deps test build docker

export GOOS ?= linux
export GOARCH ?= amd64
export CGO_ENABLED ?= 0

CI_BUILD_NUMBER ?= 0

LDFLAGS += -X "main.buildDate=$(shell date -u '+%Y-%m-%d %H:%M:%S %Z')"
LDFLAGS += -X "main.build=$(CI_BUILD_NUMBER)"

clean:
	go clean -i ./...

deps:
	go get -t ./...

test:
	go test -cover ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	go build -ldflags '-s -w $(LDFLAGS)'

docker:
	docker build --rm=true -t plugins/drone-ssh .
