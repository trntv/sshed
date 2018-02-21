BUILD := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags HEAD)

LDFLAGS=-ldflags "-X=main.version=$(VERSION) -X=main.build=$(BUILD)"

.PHONY: compile build build_all fmt bootstrap

SOURCE_FOLDER := .

GOARCH ?= amd64

default: build

bootstrap:
	dep ensure

build_all: vet fmt
	for GOOS in darwin linux; do \
		$(MAKE) compile GOOS=$$GOOS GOARCH=$(GOARCH) BINARY=build/sshed-$(VERSION)-$$GOOS-amd64; \
	done

compile:
	CGO_ENABLED=0 go build -v $(LDFLAGS) -o $(BINARY) $(SOURCE_FOLDER)/cmd

build: vet fmt
	$(MAKE) compile BINARY=build/sshed

vet:
	go vet $(SOURCE_FOLDER)/...

fmt:
	go fmt $(SOURCE_FOLDER)/...

test:
	go test $(SOURCE_FOLDER)/...

checksum:
	for GOOS in darwin linux; do \
		BINARY=build/sshed-$(VERSION)-$$GOOS-$(GOARCH); \
		openssl sha -sha256 $$BINARY > $$BINARY.sha256 ; \
	done

