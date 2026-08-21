.PHONY: build run clean test help

# Variables
BINARY_NAME=pictl
GO=go
GOFLAGS=-v

help:
	@echo "Available targets:"
	@echo "  make build       - Build the application"
	@echo "  make run         - Run the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Lint code"

build:
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

run: build
	./$(BINARY_NAME)

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

test:
	$(GO) test -v ./...

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run ./...

.DEFAULT_GOAL := help
