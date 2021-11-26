.PHONY: build-cross dist build clean run help go/* build

name       := fargate-non-exit-0-notify

go_version := $(shell cat $(realpath .go-version))
bindir     := $(realpath bin/)
go         := $(bindir)/go/bin/go
arch       := $(shell arch)

all: go/install ## build
	GOOS=linux GOARCH=amd64 $(go) build -o bin/$(name) *.go


### go install
go/install: file         = go.tar.gz
go/install: download_url = https://go.dev/dl/go$(go_version).darwin-$(arch).tar.gz
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

