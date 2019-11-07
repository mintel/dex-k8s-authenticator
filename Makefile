GOPATH=$(shell pwd)/vendor:$(shell pwd)
GOBIN=$(shell pwd)/bin
GOFILES=$(wildcard *.go)
GOPROXY ?= ""
GONAME=dex-k8s-authenticator
DOCKER_REPO=mintel/dex-k8s-authenticator
DOCKER_TAG=latest

all: build 

get:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go get -d .

build: get
	@echo "Building $(GOFILES) to ./bin"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build -o bin/$(GONAME) $(GOFILES)

container:
	@echo "Building container image"
	docker build -t ${DOCKER_REPO}:${DOCKER_TAG} --build-arg GOPROXY=${GOPROXY} .

build:
	GO111MODULE=on go build -o $(OUT_BIN) main.go

clean:
	rm -rf $(OUT_BIN)

clean:
	@echo "Cleaning"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean
	rm -rf ./bin
	rm -rf ./vendor

.PHONY: build get clean container
