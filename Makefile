#!/usr/bin/make -f

# Project variables
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell git describe --tags --match 'v*' --abbrev=8 | sed 's/-g/-/' | sed 's/-[0-9]*-/-/')
GOFLAGS:=""
GOTAGS:=""

SHELL := /bin/bash # Use bash syntax

# currently installed Go version
GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)

# minimum supported Go version
GO_MINIMUM_MAJOR_VERSION = $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f2 | cut -d'.' -f1)
GO_MINIMUM_MINOR_VERSION = $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f2 | cut -d'.' -f2)

RED=\033[0;31m
GREEN=\033[0;32m
LGREEN=\033[1;32m
NOCOLOR=\033[0m
GO_CURR_VERSION=$(shell echo -e "Current Go version: $(LGREEN)$(GO_MAJOR_VERSION).$(GREEN)$(GO_MINOR_VERSION)$(NOCOLOR)")
GO_VERSION_ERR_MSG=$(shell echo -e '$(RED)❌ ERROR$(NOCOLOR): Go version $(LGREEN)$(GO_MINIMUM_MAJOR_VERSION).$(GREEN)$(GO_MINIMUM_MINOR_VERSION)$(NOCOLOR)+ is required')

GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)

BUILD_FOLDER = ./build

ldflags =	-X github.com/elys-network/post-upgrade-snapshot-generator/version.Name=post-upgrade-snapshot-generator \
			-X github.com/elys-network/post-upgrade-snapshot-generator/version.AppName=post-upgrade-snapshot-generator \
			-X github.com/elys-network/post-upgrade-snapshot-generator/version.Version=$(VERSION) \
			-X github.com/elys-network/post-upgrade-snapshot-generator/version.Commit=$(COMMIT)
build_flags = -ldflags '$(ldflags)' -tags '$(GOTAGS)'

## install: Install post-upgrade-snapshot-generator binary in $GOBIN
install: check-version go.sum
	@echo Installing Post upgrade snapshot generator binary...
	@GOFLAGS=$(GOFLAGS) go build $(build_flags) -o $(HOME)/go/bin/post-upgrade-snapshot-generator ./cmd
	@post-upgrade-snapshot-generator version

## build: Build post-upgrade-snapshot-generator binary
build: check-version go.sum
	@echo Building Post upgrade snapshot generator binary...
	@-mkdir -p $(BUILD_FOLDER) 2> /dev/null
	@GOFLAGS=$(GOFLAGS) go build $(build_flags) -o $(BUILD_FOLDER)/post-upgrade-snapshot-generator ./cmd

.PHONY: install build

## clean: Clean build files. Runs `go clean` internally.
clean:
	@echo Cleaning build cache...
	@rm -rf $(BUILD_FOLDER) 2> /dev/null
	@go clean ./...

.PHONY: clean

## go-mod-cache: Retrieve the go modules and store them in the local cache
go-mod-cache: go.sum
	@echo "--> Retrieve the go modules and store them in the local cache."
	@go mod download

## go.sum: Ensure dependencies have not been modified
go.sum: go.mod
	@echo "--> Make sure that the dependencies haven't been altered."
	@go mod verify

# Add check to make sure we are using the proper Go version before proceeding with anything
check-version:
	@echo '$(GO_CURR_VERSION)'
	@if [[ $(GO_MAJOR_VERSION) -eq $(GO_MINIMUM_MAJOR_VERSION) && $(GO_MINOR_VERSION) -ge $(GO_MINIMUM_MINOR_VERSION) ]]; then \
		exit 0; \
	elif [[ $(GO_MAJOR_VERSION) -lt $(GO_MINIMUM_MAJOR_VERSION) ]]; then \
		echo '$(GO_VERSION_ERR_MSG)'; \
		exit 1; \
	elif [[ $(GO_MINOR_VERSION) -lt $(GO_MINIMUM_MINOR_VERSION) ]]; then \
		echo '$(GO_VERSION_ERR_MSG)'; \
		exit 1; \
	fi

.PHONY: go-mod-cache go.sum check-version

help: Makefile
	@echo
	@echo " Choose a command to run, or just run 'make' for building all binaries."
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: help

.DEFAULT_GOAL := build

GORELEASER_IMAGE := ghcr.io/goreleaser/goreleaser-cross:v$(GO_VERSION)

## release: Build binaries for all platforms and generate checksums
ifdef GITHUB_TOKEN
release:
	docker run \
		--rm \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/post-upgrade-snapshot-generator \
		-w /go/src/post-upgrade-snapshot-generator \
		$(GORELEASER_IMAGE) \
		release \
		--clean
else
release:
	@echo "Error: GITHUB_TOKEN is not defined. Please define it before running 'make release'."
endif

## release-dry-run: Dry-run build process for all platforms and generate checksums
release-dry-run:
	docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/post-upgrade-snapshot-generator \
		-w /go/src/post-upgrade-snapshot-generator \
		$(GORELEASER_IMAGE) \
		release \
		--clean \
		--skip=publish

## release-snapshot: Build snapshots for all platforms and generate checksums
release-snapshot:
	docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/post-upgrade-snapshot-generator \
		-w /go/src/post-upgrade-snapshot-generator \
		$(GORELEASER_IMAGE) \
		release \
		--clean \
		--snapshot \
		--skip-validate \
		--skip=publish

.PHONY: release release-dry-run release-snapshot