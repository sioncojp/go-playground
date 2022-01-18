REVISION := $(shell git describe --always)
LDFLAGS	 := -ldflags="-X \"main.Revision=$(REVISION)\""

.PHONY: help go/* dist test clean build build build/* run

name		:= slack-reaction-hook
linux_name	:= $(name)-linux-amd64
darwin_name	:= $(name)-darwin-amd64
darwin_arm_name	:= $(name)-darwin-arm64

go_version  := $(shell cat $(realpath .go-version))
bindir      := $(realpath bin/)
goroot      := $(bindir)/go
go          := $(bindir)/go/bin/go
arch        := $(shell arch)
go_mod      := $(realpath go.mod)
go_mod_name := sioncojp/playgroud/go/$(name)

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-22s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

### go install
go/install: file         = go.tar.gz
go/install: download_url = https://golang.org/dl/go$(go_version).darwin-$(arch).tar.gz
go/install:
# If you have a different version, delete it.
	@if [ -f $(go) ]; then \
		$(go) version | grep -q "$(go_version)" || rm -f $(bindir)/go; \
	fi

# If the file is not there, download it.
	@if [ ! -f $(go) ]; then \
		curl -L -fsS --retry 2 -o $(file) $(download_url) && \
		tar zxvf $(file) -C $(bindir) && rm -f $(file); \
	fi

dist: build-docker ## create .tar.gz linux & darwin to /bin
	cd bin && tar zcvf $(linux_name).tar.gz $(linux_name) && rm -f $(linux_name)
	cd bin && tar zcvf $(darwin_name).tar.gz $(darwin_name) && rm -f $(darwin_name)
	cd bin && tar zcvf $(darwin_arch_name).tar.gz $(darwin_arch_name) && rm -f $(darwin_arch_name)

gomod/init:
	@if [ ! -f $(go_mod) ]; then \
		$GOROOT=$(goroot) $(go) mod init $(go_mod_name); \
	fi

test: ## go test
	GOROOT=$(goroot) $(go) test -v $$(go list ./... | grep -v /vendor/)

clean: ## remove bin/*
	rm -f bin/*

build: go/install gomod/init ## build
	GOROOT=$(goroot) $(go) get -u
	GOROOT=$(goroot) $(go) build -o bin/$(name) *.go

build/cross: go/install ## create to build for linux & darwin to bin/
	GOOS=linux GOARCH=amd64 $(go) build -o bin/$(linux_name) $(LDFLAGS) cmd/$(name)/*.go
	GOOS=darwin GOARCH=amd64 $(go) build -o bin/$(darwin_name) $(LDFLAGS) cmd/$(name)/*.go
	GOOS=darwin GOARCH=arm64 $(go) build -o bin/$(darwin_arm_name) $(LDFLAGS) cmd/$(name)/*.go

get:
	GOROOT=$(goroot) $(go) get $(MODULE)