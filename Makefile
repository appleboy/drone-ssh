.PHONY: test drone-ssh build fmt vet errcheck lint install update release-dirs release-build release-copy release-check release coverage

DIST := dist
EXECUTABLE := drone-ssh

# for dockerhub
DEPLOY_ACCOUNT := appleboy
DEPLOY_IMAGE := $(EXECUTABLE)
GOFMT ?= gofmt "-s"

TARGETS ?= linux darwin windows
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
SOURCES ?= $(shell find . -name "*.go" -type f)
TAGS ?=
LDFLAGS ?= -X 'main.Version=$(VERSION)'

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(DRONE_TAG)
else
	VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
endif

all: build

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@files=$$($(GOFILES) | xargs $(GOFMT) -l); if [ -n "$$files" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${files}"; \
		exit 1; \
		fi;

fmt:
	$(GOFMT) -w $(GOFILES)

vet:
	go vet $(PACKAGES)

errcheck:
	@hash errcheck > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kisielk/errcheck; \
	fi
	errcheck $(PACKAGES)

lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

unconvert:
	@hash unconvert > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/mdempsky/unconvert; \
	fi
	for PKG in $(PACKAGES); do unconvert -v $$PKG || exit 1; done;

test: fmt-check
	for PKG in $(PACKAGES); do go test -v -cover -coverprofile $$GOPATH/src/$$PKG/coverage.txt $$PKG || exit 1; done;

html:
	go tool cover -html=coverage.txt

install: $(SOURCES)
	go install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	go build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o $@

release: release-dirs release-build release-copy release-check

release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

release-build:
	@hash gox > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/mitchellh/gox; \
	fi
	gox -os="$(TARGETS)" -arch="amd64 386" -tags="$(TAGS)" -ldflags="-s -w $(LDFLAGS)" -output="$(DIST)/binaries/$(EXECUTABLE)-$(VERSION)-{{.OS}}-{{.Arch}}"

release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

# for docker.
static_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o $(DEPLOY_IMAGE)

docker_image:
	docker build -t $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE) .

docker: static_build docker_image

docker_deploy:
ifeq ($(tag),)
	@echo "Usage: make $@ tag=<tag>"
	@exit 1
endif
	# deploy image
	docker tag $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):latest $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):$(tag)
	docker push $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE):$(tag)

coverage:
	sed -i '/main.go/d' .cover/coverage.txt
	curl -s https://codecov.io/bash > .codecov && \
	chmod +x .codecov && \
	./.codecov -f .cover/coverage.txt

clean:
	go clean -x -i ./...
	rm -rf coverage.txt $(EXECUTABLE) $(DIST) vendor

ssh-server:
	adduser -h /home/drone-scp -s /bin/bash -D -S drone-scp
	passwd -d drone-scp
	mkdir -p /home/drone-scp/.ssh
	chmod 700 /home/drone-scp/.ssh
	cp tests/.ssh/id_rsa.pub /home/drone-scp/.ssh/authorized_keys
	chown -R drone-scp /home/drone-scp/.ssh
	# install ssh and start server
	apk update && apk add openssh openrc
	rm -rf /etc/ssh/ssh_host_rsa_key /etc/ssh/ssh_host_dsa_key
	./tests/entrypoint.sh /usr/sbin/sshd -D &

version:
	@echo $(VERSION)
