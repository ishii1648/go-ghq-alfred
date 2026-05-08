SHELL := /bin/bash
BIN_NAME := dist/go-ghq-alfred
BIN_NAME_AMD64 := build/go-ghq-alfred-amd64
BIN_NAME_ARM64 := build/go-ghq-alfred-arm64
WF_NAME := ghq-alfred.alfredworkflow
ASSETS := $(shell find -f ./resources)
TESTDIR := testdir

$(TESTDIR):
	mkdir -p $(TESTDIR)/{data,cache}

test: $(TESTDIR)
	env alfred_workflow_bundleid=testid \
		alfred_workflow_cache=$(TESTDIR)/cache \
		alfred_workflow_data=$(TESTDIR)/data \
		go test ./... -v

$(BIN_NAME): main.go
	go build -o $(BIN_NAME) .

build: $(BIN_NAME)

$(BIN_NAME_AMD64): main.go
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_NAME_AMD64) .

$(BIN_NAME_ARM64): main.go
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_NAME_ARM64) .

# Produce a fat Mach-O binary that runs on both Intel and Apple Silicon Macs.
# Overwrites $(BIN_NAME) with the universal output.
build-universal: $(BIN_NAME_AMD64) $(BIN_NAME_ARM64)
	mkdir -p dist
	lipo -create -output $(BIN_NAME) $(BIN_NAME_AMD64) $(BIN_NAME_ARM64)

$(WF_NAME): build-universal $(ASSETS)
	cp -r resources/* dist/
	cd dist && zip -r ../$(WF_NAME) ./*

dist: $(WF_NAME)

clean:
	rm -rf dist build $(WF_NAME)

.PHONY: build build-universal dist test clean
