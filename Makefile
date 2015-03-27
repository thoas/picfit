ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${ROOT_DIR}/application/constants.go)

BIN_DIR = $(ROOT_DIR)/bin

test: unit

unit:
	@(go list ./... | xargs -n1 go test -v -cover)

deps:
	godep save ./...

all: picfit
	@(mkdir -p $(BIN_DIR))

build:
	@(go get github.com/tools/godep)
	@(godep restore)
	@(echo "-> Compiling picfit binary")
	@(mkdir -p $(BIN_DIR))
	@(godep go build -o $(BIN_DIR)/picfit)
	@(echo "-> picfit binary created")

format:
	@(go fmt ./...)
	@(go vet ./...)

.PNONY: all test format
