.PHONY: build run clean test help

# Variables
DASHBOARD_BINARY=dashboard
PICTL_BINARY=pictl
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
	$(GO) build $(GOFLAGS) -o $(DASHBOARD_BINARY) ./cmd/dashboard
	$(GO) build $(GOFLAGS) -o $(PICTL_BINARY) ./cmd/pictl

run: build
	./$(DASHBOARD_BINARY)

clean:
	$(GO) clean
	rm -f $(DASHBOARD_BINARY) $(PICTL_BINARY)

test:
	$(GO) test -v ./...

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run ./...

.DEFAULT_GOAL := help
