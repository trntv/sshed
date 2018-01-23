BUILD := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags HEAD)

LDFLAGS=-ldflags "-X=main.version=$(VERSION) -X=main.build=$(BUILD)"

.PHONY: compile build build_all fmt lint test itest vet bootstrap

SOURCE_FOLDER := .

GOARCH ?= amd64

default: build

bootstrap:
	dep ensure

build_all: vet fmt
	for GOOS in darwin linux; do \
		$(MAKE) compile GOOS=$$GOOS GOARCH=$(GOARCH) BINARY=build/sshme-$(VERSION)-$$GOOS-amd64; \
	done

compile:
	CGO_ENABLED=0 go build -v $(LDFLAGS) -o $(BINARY) $(SOURCE_FOLDER)/cmd

build: vet fmt
	$(MAKE) compile BINARY=build/sshme

fmt:
	go fmt ./cmd ./commands ./db

vet:
	go fmt ./cmd ./commands ./db

lint:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 golint

test:
	go test ./cmd ./commands ./db

checksum:
	for GOOS in darwin linux; do \
		BINARY=build/sshme-$(VERSION)-$$GOOS-$(GOARCH); \
		openssl sha -sha256 $$BINARY > $$BINARY.sha256 ; \
	done

