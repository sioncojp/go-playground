.PHONY: go/* build run

NAME := s3-move-other-bucket

### go install
GO_VERSION := 1.16.4
BINDIR     := $(realpath bin/)
GO         := $(BINDIR)/go/bin/go
ARCH       := $(shell arch)

go/install: file         = go.tar.gz
go/install: download_url = https://golang.org/dl/go$(GO_VERSION).darwin-$(ARCH).tar.gz
go/install:
# バージョンが違ったら削除する
	@if [ -f $(GO) ]; then \
		$(GO) version | grep -q "$(GO_VERSION)" || rm -f $(BINDIR)/go; \
	fi

# ファイルがなかったらダウンロード
	@if [ ! -f $(GO) ]; then \
		curl -L -fsS --retry 2 -o $(file) $(download_url) && \
		tar zxvf $(file) -C $(BINDIR) && rm -f $(file); \
	fi

build: go/install
	$(GO) build -o $(BINDIR)/$(NAME) *.go

build/cross: go/install
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BINDIR)/$(NAME)_linux *.go
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BINDIR)/$(NAME)_darwin_amd64 *.go
