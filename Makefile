.PHONY: all build test clean lint fmt vet cover bench install help

# Variables
BINARY_NAME=pdfsdk-example
VERSION=$(shell grep 'Version = ' pdfsdk.go | cut -d'"' -f2)
GO=go
GOFLAGS=-v
COVERAGE_FILE=coverage.out

# Default target
all: fmt vet test build

# Build the example binary
build:
	@echo "🔨 Building..."
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) ./cmd/main.go
	@echo "✅ Built bin/$(BINARY_NAME)"

# Run all tests
test:
	@echo "🧪 Running tests..."
	$(GO) test $(GOFLAGS) . ./service/... ./pkg/...
	@echo "✅ All tests passed"

# Run tests with coverage
cover:
	@echo "📊 Running tests with coverage..."
	$(GO) test -coverprofile=$(COVERAGE_FILE) -covermode=atomic . ./service/... ./pkg/...
	$(GO) tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "✅ Coverage report: coverage.html"

# Run benchmarks
bench:
	@echo "⚡ Running benchmarks..."
	$(GO) test -bench=. -benchmem . ./service/...
	@echo "✅ Benchmarks completed"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	$(GO) fmt ./...
	@echo "✅ Code formatted"

# Vet code
vet:
	@echo "🔍 Vetting code..."
	$(GO) vet ./...
	@echo "✅ Vet passed"

# Lint code (requires golangci-lint)
lint:
	@echo "🔎 Linting code..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
	@echo "✅ Lint passed"

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	rm -f $(COVERAGE_FILE) coverage.html
	rm -f *.test
	@echo "✅ Cleaned"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "✅ Dependencies installed"

# Generate documentation
docs:
	@echo "📚 Generating documentation..."
	@which godoc > /dev/null || go install golang.org/x/tools/cmd/godoc@latest
	@echo "Open http://localhost:6060/pkg/github.com/infosec554/convert-pdf-go-sdk/"
	godoc -http=:6060

# Run example
run:
	@echo "🚀 Running example..."
	$(GO) run ./cmd/main.go

# Docker build
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t pdfsdk:$(VERSION) .
	@echo "✅ Built pdfsdk:$(VERSION)"

# Docker compose up
docker-up:
	@echo "🐳 Starting services..."
	docker-compose up -d
	@echo "✅ Services started"

# Docker compose down
docker-down:
	@echo "🐳 Stopping services..."
	docker-compose down
	@echo "✅ Services stopped"

# Show version
version:
	@echo "v$(VERSION)"

# Help
help:
	@echo "PDF SDK v$(VERSION) - Available targets:"
	@echo ""
	@echo "  make build      - Build the example binary"
	@echo "  make test       - Run all tests"
	@echo "  make cover      - Run tests with coverage report"
	@echo "  make bench      - Run benchmarks"
	@echo "  make fmt        - Format code"
	@echo "  make vet        - Vet code"
	@echo "  make lint       - Lint code (requires golangci-lint)"
	@echo "  make clean      - Clean build artifacts"
	@echo "  make deps       - Install dependencies"
	@echo "  make docs       - Generate and serve documentation"
	@echo "  make run        - Run example"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-up  - Start Docker services"
	@echo "  make docker-down - Stop Docker services"
	@echo "  make version    - Show version"
	@echo "  make help       - Show this help"
