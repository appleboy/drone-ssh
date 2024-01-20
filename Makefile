DIST := dist
EXECUTABLE := drone-ssh
GOFMT ?= gofumpt -l
DIST := dist
DIST_DIRS := $(DIST)/binaries $(DIST)/release
GO ?= go
SHASUM ?= shasum -a 256
GOFILES := $(shell find . -name "*.go" -type f)
HAS_GO = $(shell hash $(GO) > /dev/null 2>&1 && echo "GO" || echo "NOGO" )
XGO_PACKAGE ?= src.techknowlogick.com/xgo@latest
XGO_VERSION := go-1.19.x
GXZ_PAGAGE ?= github.com/ulikunitz/xz/cmd/gxz@v0.5.11

LINUX_ARCHS ?= linux/amd64,linux/arm64
DARWIN_ARCHS ?= darwin-10.12/amd64,darwin-10.12/arm64
WINDOWS_ARCHS ?= windows/*

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

ifeq ($(HAS_GO), GO)
	GOPATH ?= $(shell $(GO) env GOPATH)
	export PATH := $(GOPATH)/bin:$(PATH)

	CGO_EXTRA_CFLAGS := -DSQLITE_MAX_VARIABLE_NUMBER=32766
	CGO_CFLAGS ?= $(shell $(GO) env CGO_CFLAGS) $(CGO_EXTRA_CFLAGS)
endif

ifeq ($(OS), Windows_NT)
	GOFLAGS := -v -buildmode=exe
	EXECUTABLE ?= $(EXECUTABLE).exe
else ifeq ($(OS), Windows)
	GOFLAGS := -v -buildmode=exe
	EXECUTABLE ?= $(EXECUTABLE).exe
else
	GOFLAGS := -v
	EXECUTABLE ?= $(EXECUTABLE)
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(DRONE_TAG)
else
	VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
endif

TAGS ?=
LDFLAGS ?= -X 'main.Version=$(VERSION)'

all: build

fmt:
	@hash gofumpt > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install mvdan.cc/gofumpt; \
	fi
	$(GOFMT) -w $(GOFILES)

vet:
	$(GO) vet ./...

.PHONY: fmt-check
fmt-check:
	@hash gofumpt > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) install mvdan.cc/gofumpt; \
	fi
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

test:
	@$(GO) test -v -cover -coverprofile coverage.txt ./... && echo "\n==>\033[32m Ok\033[m\n" || exit 1

install: $(GOFILES)
	$(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

$(EXECUTABLE): $(GOFILES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o bin/$@

build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/amd64/$(DEPLOY_IMAGE)

build_linux_i386:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GO) build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/i386/$(DEPLOY_IMAGE)

build_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/arm64/$(DEPLOY_IMAGE)

build_linux_arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GO) build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/arm/$(DEPLOY_IMAGE)

ssh-server:
	adduser -h /home/drone-scp -s /bin/sh -D -S drone-scp
	echo drone-scp:1234 | chpasswd
	mkdir -p /home/drone-scp/.ssh
	chmod 700 /home/drone-scp/.ssh
	cat tests/.ssh/id_rsa.pub >> /home/drone-scp/.ssh/authorized_keys
	cat tests/.ssh/test.pub >> /home/drone-scp/.ssh/authorized_keys
	chmod 600 /home/drone-scp/.ssh/authorized_keys
	chown -R drone-scp /home/drone-scp/.ssh
	# add public key to root user
	mkdir -p /root/.ssh
	chmod 700 /root/.ssh
	cat tests/.ssh/id_rsa.pub >> /root/.ssh/authorized_keys
	cat tests/.ssh/test.pub >> /root/.ssh/authorized_keys
	chmod 600 /root/.ssh/authorized_keys
	# Append the following entry to run ALL command without a password for a user named drone-scp:
	cat tests/sudoers >> /etc/sudoers.d/sudoers
	# install ssh and start server
	apk add --update openssh openrc
	rm -rf /etc/ssh/ssh_host_rsa_key /etc/ssh/ssh_host_dsa_key
	sed -i 's/^#PubkeyAuthentication yes/PubkeyAuthentication yes/g' /etc/ssh/sshd_config
	sed -i 's/AllowTcpForwarding no/AllowTcpForwarding yes/g' /etc/ssh/sshd_config
	sed -i 's/^#ListenAddress 0.0.0.0/ListenAddress 0.0.0.0/g' /etc/ssh/sshd_config
	sed -i 's/^#ListenAddress ::/ListenAddress ::/g' /etc/ssh/sshd_config
	./tests/entrypoint.sh /usr/sbin/sshd -D &

coverage:
	sed -i '/main.go/d' coverage.txt

.PHONY: deps-backend
deps-backend:
	$(GO) mod download
	$(GO) install $(GXZ_PAGAGE)
	$(GO) install $(XGO_PACKAGE)

.PHONY: release
release: release-linux release-darwin release-windows release-copy release-compress release-check

$(DIST_DIRS):
	mkdir -p $(DIST_DIRS)

.PHONY: release-windows
release-windows: | $(DIST_DIRS)
	CGO_CFLAGS="$(CGO_CFLAGS)" $(GO) run $(XGO_PACKAGE) -go $(XGO_VERSION) -buildmode exe -dest $(DIST)/binaries -tags 'netgo osusergo $(TAGS)' -ldflags '-linkmode external -extldflags "-static" $(LDFLAGS)' -targets '$(WINDOWS_ARCHS)' -out $(EXECUTABLE)-$(VERSION) .
ifeq ($(CI),true)
	cp -r /build/* $(DIST)/binaries/
endif

.PHONY: release-linux
release-linux: | $(DIST_DIRS)
	CGO_CFLAGS="$(CGO_CFLAGS)" $(GO) run $(XGO_PACKAGE) -go $(XGO_VERSION) -dest $(DIST)/binaries -tags 'netgo osusergo $(TAGS)' -ldflags '-linkmode external -extldflags "-static" $(LDFLAGS)' -targets '$(LINUX_ARCHS)' -out $(EXECUTABLE)-$(VERSION) .
ifeq ($(CI),true)
	cp -r /build/* $(DIST)/binaries/
endif

.PHONY: release-darwin
release-darwin: | $(DIST_DIRS)
	CGO_CFLAGS="$(CGO_CFLAGS)" $(GO) run $(XGO_PACKAGE) -go $(XGO_VERSION) -dest $(DIST)/binaries -tags 'netgo osusergo $(TAGS)' -ldflags '$(LDFLAGS)' -targets '$(DARWIN_ARCHS)' -out $(EXECUTABLE)-$(VERSION) .
ifeq ($(CI),true)
	cp -r /build/* $(DIST)/binaries/
endif

.PHONY: release-copy
release-copy: | $(DIST_DIRS)
	cd $(DIST); for file in `find . -type f -name "*"`; do cp $${file} ./release/; done;

.PHONY: release-check
release-check: | $(DIST_DIRS)
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "checksumming $${file}" && $(SHASUM) `echo $${file} | sed 's/^..//'` > $${file}.sha256; done;

.PHONY: release-compress
release-compress: | $(DIST_DIRS)
	cd $(DIST)/release/; for file in `find . -type f -name "*"`; do echo "compressing $${file}" && $(GO) run $(GXZ_PAGAGE) -k -9 $${file}; done;

clean:
	$(GO) clean -x -i ./...
	rm -rf coverage.txt $(EXECUTABLE) $(DIST)

version:
	@echo $(VERSION)
