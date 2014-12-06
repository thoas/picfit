BIN_DIR = $(ROOT_DIR)/bin
PICFIT_CMD_DIR = $(CMD_DIR)/picfit
PICFIT_BIN = $(BIN_DIR)/picfit
VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${ROOT_DIR}/application/constants.go)

test: unit

unit:
	@(go list ./... | xargs -n1 go test -v)

deps:
	godep save ./...

all: picfit
	@(mkdir -p $(BIN_DIR))

picfit:
	@(go get github.com/tools/godep)
	@(echo "-> Compiling picfit binary")
	@(mkdir -p $(BIN_DIR))
	@(cd $(PICFIT_CMD_DIR) && godep go build -o $(PICFIT_BIN)) 
	@(echo "-> picfit binary created: $(PICFIT_BIN)")

format:
	@(go fmt ./...)
	@(go vet ./...)

.PNONY: all test format
