GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=dex-k8s-authenticator
TAG=latest

all: build 

.PHONY: build
build:
	@echo "Building $(GOFILES) to ./bin"
	GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

.PHONY: container
container:
	@echo "Building container image"
	docker build -t ${GONAME}:${TAG} .
.PHONY: clean
clean:
	@echo "Cleaning"
	GOBIN=$(GOBIN) go clean
	rm -rf ./bin

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix: lint
	golangci-lint run --fix
