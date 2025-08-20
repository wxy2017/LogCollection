.PHONY: all build debug tidy clean help

BIN_FILE=logCollection
MAIN_FILE=./main.go

SHELL := /bin/bash
BASEDIR = $(shell pwd)

#amd64 arm64
GOARCH=amd64

# build with verison infos
versionDir = "FDManagerGo/version"

# 默认执行
all: build

# 打包成二进制文件
build: tidy
	@CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH)  go build -v -trimpath -o $(BIN_FILE) $(MAIN_FILE)

debug: tidy
	@CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH)  go build -v -trimpath -gcflags "all=-N -m -l" -o $(BIN_FILE) $(MAIN_FILE)

tidy:
	@go mod tidy

clean:
	@go clean -x
	rm -f $(BIN_FILE)

help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  the default target is 'all'."
	@echo "  all          Build the project"
	@echo "  build        Build the project"
	@echo "  debug        Build the project with debug mode"
	@echo "  tidy         Tidy the go.mod file"
	@echo "  clean        Clean the project"
	@echo "  help         Print this help message"