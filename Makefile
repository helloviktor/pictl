.PHONY: build run clean test help

# Variables
OUTPUT_DIR=out
DASHBOARD_BINARY=$(OUTPUT_DIR)/dashboard
PICTL_BINARY=$(OUTPUT_DIR)/pictl
GO=go
GOFLAGS=-v
GOARCH?=$(shell $(GO) env GOARCH)

help:
	@echo "Available targets:"
	@echo "  make build       - Build the application"
	@echo "                    Set GOARCH to choose the binary architecture (for example, GOARCH=arm64)"
	@echo "  make run         - Run the application"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Lint code"

build:
	mkdir -p $(OUTPUT_DIR)
	GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o $(DASHBOARD_BINARY) ./cmd/dashboard
	GOARCH=$(GOARCH) $(GO) build $(GOFLAGS) -o $(PICTL_BINARY) ./cmd/pictl

run: build
	./$(DASHBOARD_BINARY)

clean:
	$(GO) clean
	rm -rf $(OUTPUT_DIR)

test:
	$(GO) test -v ./...

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run ./...

.DEFAULT_GOAL := help
