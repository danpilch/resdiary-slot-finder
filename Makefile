# Makefile for Go project

# Variables
GO = go
BUILD_DIR = build
SRC = ./...
BIN = ./
LINTER = golangci-lint

# Default target
all: build

# Build the project
build: fmt vet
	$(GO) build -o $(BIN)

# Run the project
run: build
	$(BIN)

# Format the code
fmt:
	$(GO) fmt $(SRC)

# Lint the project
lint:
	$(LINTER) run

# Test the project
test:
	$(GO) test -v $(SRC)

# Install dependencies
install:
	$(GO) mod tidy

# Vet the code (checks for suspicious constructs)
vet:
	$(GO) vet $(SRC)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# View Go version
version:
	$(GO) version

# Help
help:
	@echo "Makefile commands:"
	@echo "  make build     - Build the project"
	@echo "  make run       - Run the project"
	@echo "  make fmt       - Format the code"
	@echo "  make lint      - Lint the code"
	@echo "  make test      - Run tests"
	@echo "  make install   - Install dependencies"
	@echo "  make vet       - Run vet checks"
	@echo "  make clean     - Clean build artifacts"
	@echo "  make version   - Show Go version"
	@echo "  make help      - Show this help message"

.PHONY: all build run fmt lint test install vet clean version help
