.PHONY: test
test:
	@go test -v ./... -cover -race

.PHONY: build
build:
	@rm -rf ./bin
	@mkdir -p ./bin
	@go build -ldflags="-s -w" -o ./bin/server ./cmd/server/*.go

.PHONY: lint
lint:
	@golangci-lint run ./... -v