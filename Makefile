GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GONAME=dex-k8s-authenticator
TAG=0.0.1

all: build 

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -d .

build: get
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

container:
	@echo "Building container image"
	docker build -t ${GONAME}:${TAG} .

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	rm -rf ./bin
	rm -rf ./vendor

.PHONY: build get clean container
