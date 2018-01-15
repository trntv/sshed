BUILD := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags HEAD)

LDFLAGS=-ldflags "-X=main.version=$(VERSION) -X=main.build=$(BUILD)"

.PHONY: compile build build_all fmt lint test itest vet bootstrap

SOURCE_FOLDER := .

BINARY_PATH ?= $(GOPATH)/bin/sshed

GOARCH ?= amd64

ifdef GOOS
BINARY_PATH :=$(BINARY_PATH).$(GOOS)-$(GOARCH)
endif

default: build

build_all: vet fmt
	for GOOS in darwin linux windows; do \
		$(MAKE) compile GOOS=$$GOOS GOARCH=amd64 ; \
	done

compile:
	CGO_ENABLED=0 go build -i -v $(LDFLAGS) -o $(BINARY_PATH) $(SOURCE_FOLDER)/cmd

build: vet fmt compile

fmt:
	go fmt ./cmd ./commands ./db

vet:
	go fmt ./cmd ./commands ./db

lint:
	go list $(SOURCE_FOLDER)/... | grep -v /vendor/ | xargs -L1 golint

test:
	go test ./cmd ./commands ./db

itest:
	$(MAKE) compile GOOS=linux GOARCH=amd64
	bats $(SPECS)

bootstrap:
	dep ensure
