# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Build info
BINARY_NAME=goutils-example
BINARY_PATH=./examples/

# Test parameters
TEST_ARGS=-v -race -coverprofile=coverage.out

.PHONY: all build clean test coverage lint fmt help deps example

all: fmt lint test build

# Build the example application
build:
	$(GOBUILD) -o $(BINARY_PATH)$(BINARY_NAME) $(BINARY_PATH)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)$(BINARY_NAME)
	rm -f coverage.out

# Run tests
test:
	$(GOTEST) $(TEST_ARGS) ./...

# Run tests with coverage
coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linting
lint:
	$(GOLINT) run ./...

# Format code
fmt:
	$(GOFMT) -s -w .

# Update dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run the example
example: build
	$(BINARY_PATH)$(BINARY_NAME)

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run benchmarks
bench:
	$(GOTEST) -bench=. -benchmem ./...

# Check for security vulnerabilities
security:
	@which govulncheck > /dev/null || $(GOGET) golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all . > docs.txt
	@echo "Documentation generated: docs.txt"

# Check code quality
quality: fmt lint test coverage

# Help
help:
	@echo "Available commands:"
	@echo "  all          - Run fmt, lint, test, and build"
	@echo "  build        - Build the example application"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage report"
	@echo "  lint         - Run linting"
	@echo "  fmt          - Format code"
	@echo "  deps         - Update dependencies"
	@echo "  example      - Build and run the example"
	@echo "  install-tools- Install development tools"
	@echo "  bench        - Run benchmarks"
	@echo "  security     - Check for security vulnerabilities"
	@echo "  docs         - Generate documentation"
	@echo "  quality      - Run all quality checks"
	@echo "  help         - Show this help message"
