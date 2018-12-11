test: unit

unit:
	@(go list ./... | xargs -n1 go test -v)

format:
	@(go fmt ./...)
	@(go vet ./...)

.PNONY: test
