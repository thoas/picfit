VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${ROOT_DIR}/application/constants.go)

test: unit

unit:
	@(go list ./... | xargs -n1 go test -v)

deps:
	godep save ./...

format:
	@(go fmt ./...)
	@(go vet ./...)

.PNONY: all test format
