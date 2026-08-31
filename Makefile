.PHONY: build run clean test help deploy

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
	@echo "  make deploy      - Build arm64 binaries and rsync them to the host in .env"

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

deploy:
	./scripts/deploy.sh

.DEFAULT_GOAL := help
