# Makefile for GeoRoute
# Binary output name
BINARY=georoute

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
GOLINT=golangci-lint
GOLINTCMD=$(GOLINT) run

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY) -v cmd/main.go

# Run the application
.PHONY: run
run:
	$(GORUN) cmd/main.go

# Test the application (if tests are added later)
.PHONY: test
test:
	$(GOTEST) -v ./...

# Format the code
.PHONY: fmt
fmt:
	$(GOFMT) ./...

# Lint the code (requires golangci-lint to be installed)
.PHONY: lint
lint:
	$(GOLINTCMD) ./...

# Install the binary to $GOPATH/bin
.PHONY: install
install:
	$(GOCMD) install -v

# Clean up generated files
.PHONY: clean
clean:
	rm -f $(BINARY)

# Build and run in one step
.PHONY: build-run
build-run: build
	./$(BINARY)

# Help command to display available targets
.PHONY: help
help:
	@echo "Makefile for GeoRoute"
	@echo "Usage:"
	@echo "  make          - Build the binary (default)"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Run the application directly"
	@echo "  make test     - Run tests"
	@echo "  make fmt      - Format the code"
	@echo "  make lint     - Lint the code (requires golangci-lint)"
	@echo "  make install  - Install the binary to $$GOPATH/bin"
	@echo "  make clean    - Remove the binary"
	@echo "  make build-run- Build and run the binary"
	@echo "  make help     - Show this help message"