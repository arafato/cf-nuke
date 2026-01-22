# Build metadata
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go build flags
MODULE  := github.com/arafato/cf-nuke
LDFLAGS := -ldflags "-s -w \
	-X $(MODULE)/version.version=$(VERSION) \
	-X $(MODULE)/version.commit=$(COMMIT) \
	-X $(MODULE)/version.date=$(DATE)"

# Output directory
BIN_DIR := bin

.PHONY: build build-all clean test version help

## build: Build for current platform
build:
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/cf-nuke .

## build-all: Cross-compile for all platforms
build-all: clean
	@echo "Building for Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/cf-nuke-linux-amd64 .
	@echo "Building for Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/cf-nuke-windows-amd64.exe .
	@echo "Building for macOS (Apple Silicon)..."
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BIN_DIR)/cf-nuke-macos-arm64 .
	@echo "Building for macOS (Intel)..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BIN_DIR)/cf-nuke-macos-amd64 .
	@echo "Build complete."

## clean: Remove build artifacts
clean:
	@rm -rf $(BIN_DIR)
	@mkdir -p $(BIN_DIR)

## test: Run tests
test:
	go test -v ./...

## version: Print build version info
version:
	@echo "Version: $(VERSION)"
	@echo "Commit:  $(COMMIT)"
	@echo "Date:    $(DATE)"

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
