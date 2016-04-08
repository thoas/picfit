ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${ROOT_DIR}/application/constants.go)

BIN_DIR = $(ROOT_DIR)/bin
SSL_DIR = $(ROOT_DIR)/ssl
APP_DIR = /go/src/github.com/thoas/picfit

test: unit

unit:
	@(go list ./... | grep -v "vendor/" | xargs -n1 go test -v -cover)

all: picfit
	@(mkdir -p $(BIN_DIR))

build:
	@(echo "-> Compiling picfit binary")
	@(mkdir -p $(BIN_DIR))
	@(go build -o $(BIN_DIR)/picfit)
	@(echo "-> picfit binary created")

format:
	@(go fmt ./...)
	@(go vet ./...)

build-static:
	@(echo "-> Creating statically linked binary...")
	mkdir -p $(BIN_DIR)
	@(CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BIN_DIR)/picfit)

docker-build-static: build-static
	cp -r /etc/ssl/* $(APP_DIR)/ssl/


.PNONY: all test format

docker-build:
	@(echo "-> Preparing builder...")
	@(docker build -t picfit-builder -f Dockerfile.build .)
	@(mkdir -p $(BIN_DIR))
	@(docker run --rm -v $(BIN_DIR):$(APP_DIR)/bin -v $(SSL_DIR):$(APP_DIR)/ssl picfit-builder)
