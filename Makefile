# MCP Manager CLI - Development Makefile

.PHONY: build test clean run dev install deps

# Build the application
build:
	@mkdir -p bin
	go build -o bin/mcp-manager ./cmd/mcp-manager

# Run the application
run: build
	./bin/mcp-manager

# Development mode with auto-rebuild (requires entr or similar)
dev:
	@echo "Starting development mode..."
	@echo "Press Ctrl+C to stop"
	find . -name "*.go" | entr -r make run

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install the binary to $GOPATH/bin
install: build
	go install ./cmd/mcp-manager

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# Check terminal compatibility
term-test:
	@echo "Testing terminal compatibility..."
	@echo "Terminal: $$TERM"
	@echo "Columns: $$(tput cols)"
	@echo "Lines: $$(tput lines)"
	@echo "Colors: $$(tput colors)"
	@./bin/mcp-manager --version 2>/dev/null || echo "Binary not built yet, run 'make build' first"

# Help
help:
	@echo "MCP Manager CLI - Development Commands"
	@echo ""
	@echo "  build         Build the application binary"
	@echo "  run           Build and run the application"
	@echo "  dev           Development mode with auto-rebuild"
	@echo "  test          Run all tests"
	@echo "  test-coverage Run tests with coverage report"
	@echo "  deps          Install and tidy dependencies"
	@echo "  clean         Clean build artifacts"
	@echo "  install       Install binary to GOPATH/bin"
	@echo "  fmt           Format Go code"
	@echo "  lint          Lint Go code (requires golangci-lint)"
	@echo "  term-test     Test terminal compatibility"
	@echo "  help          Show this help message"