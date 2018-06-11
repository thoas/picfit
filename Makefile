ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${ROOT_DIR}/application/constants.go)

branch = $(shell git rev-parse --abbrev-ref HEAD)
commit = $(shell git log --pretty=format:'%h' -n 1)
now = $(shell date "+%Y-%m-%d %T UTC%z")
compiler = $(shell go version)

BIN_DIR = $(ROOT_DIR)/bin
CONFIG=`pwd`/config.json
BIN = $(BIN_DIR)/picfit
SSL_DIR = $(ROOT_DIR)/ssl
APP_DIR = /go/src/github.com/thoas/picfit

test: unit

vendorize:
	find vendor/ -type f -not -path "*/.git*" -exec git add {} \;

run-server:
	@PICFIT_CONFIG_PATH=$(CONFIG) $(BIN)

serve:
	@modd

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
	CGO_ENABLED=0 go build -ldflags "\
		    -X github.com/thoas/picfit/constants.Branch=$(branch) \
		    -X github.com/thoas/picfit/constants.Revision=$(commit) \
		    -X 'github.com/thoas/picfit/constants.BuildTime=$(now)' \
		    -X 'github.com/thoas/picfit/constants.Compiler=$(compiler)'" -a -installsuffix cgo -o $(BIN_DIR)/picfit

docker-build-static: build-static


.PNONY: all test format

docker-build:
	@(echo "-> Preparing builder...")
	@(docker build -t picfit-builder -f Dockerfile.build .)
	@(mkdir -p $(BIN_DIR))
	@(docker run --rm -v $(BIN_DIR):$(APP_DIR)/bin picfit-builder)
