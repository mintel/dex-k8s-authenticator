GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=dex-k8s-authenticator
TAG=latest

all: build 

build:
	@echo "Building $(GOFILES) to ./bin"
	GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

container:
	@echo "Building container image"
	docker build -t ${GONAME}:${TAG} .

clean:
	@echo "Cleaning"
	GOBIN=$(GOBIN) go clean
	rm -rf ./bin

.PHONY: build clean container
