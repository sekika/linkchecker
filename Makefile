.PHONY: all build install clean deps

# Name of the application binary
BIN_NAME := linkchecker
# Main entry point file (Updated path)
MAIN_FILE := ./cmd/$(BIN_NAME)/main.go

all: build

# Build the binary
build:
	@echo "Building $(BIN_NAME)..."
	go build -o $(BIN_NAME) $(MAIN_FILE)

# Install using the Go toolchain (Recommended)
# This achieves the same effect as `go install github.com/sekika/linkchecker/cmd/linkchecker@latest` locally.
install: build
	@echo "Installing $(BIN_NAME) to $(GOPATH)/bin or equivalent..."
	go install $(MAIN_FILE)

# Download and tidy dependencies
deps:
	@echo "Downloading and tidying dependencies..."
	go mod tidy

# Clean up built files
clean:
	@echo "Cleaning up..."
	rm -f $(BIN_NAME)