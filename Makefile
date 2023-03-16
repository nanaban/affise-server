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

.PHONY: docker-build
docker:
	@DOCKER_SCAN_SUGGEST=false docker build -t affise-server:latest .

.PHONY: docker-run
docker-run:
	@docker run -p 8080:8080 --name affise-server affise-server:latest