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
	@(godep go build) 
	@(echo "-> picfit binary created")

format:
	@(go fmt ./...)
	@(go vet ./...)

.PNONY: all test format
