.PHONY: build test clean run docker-build docker-run lint fmt vet

# Binary name
BINARY_NAME=meta_harness
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
GOFMT=gofmt
# Build flags
LDFLAGS=-ldflags "-s -w"
# Docker
DOCKER_IMAGE=meta_harness
DOCKER_TAG=latest

## build: Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/run

## test: Run unit tests
test:
	$(GOTEST) -v -race ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## clean: Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

## run: Build and run locally
run: build
	./$(BINARY_NAME) --backend=claude --cli-path=/usr/local/bin/claude --goal="test task"

## fmt: Format code
fmt:
	$(GOFMT) -s -w .

## vet: Run go vet
vet:
	$(GOVET) ./...

## lint: Run linting (requires golangci-lint)
lint:
	golangci-lint run ./...

## docker-build: Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

## docker-run: Run Docker container
docker-run:
	docker run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) --backend=claude --cli-path=/usr/local/bin/claude --goal="test task"

## docker-clean: Remove Docker image
docker-clean:
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG)

## mod-tidy: Update go.mod and go.sum
mod-tidy:
	$(GOCMD) mod tidy

## help: Show this help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
