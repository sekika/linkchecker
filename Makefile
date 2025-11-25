.PHONY: all build install clean deps

# Name of the application binary
BIN_NAME := linkchecker
# Main package path for building and installing
MAIN_PACKAGE := github.com/sekika/linkchecker/cmd/$(BIN_NAME)

all: build

# Build the binary locally in the project root
build:
	@echo "Building $(BIN_NAME)..."
	go build -o $(BIN_NAME) $(MAIN_PACKAGE)

# Install using the Go toolchain (Recommended)
# This uses the package path to ensure the binary is correctly placed in $GOBIN ($GOPATH/bin).
install:
	@echo "Installing $(BIN_NAME) to $(GOPATH)/bin or equivalent..."
	go install $(MAIN_PACKAGE)

# Clean up locally built binary
clean:
	@echo "Cleaning up local binary..."
	rm -f $(BIN_NAME)

deps:
	@echo "Fetching dependencies..."
	go mod download