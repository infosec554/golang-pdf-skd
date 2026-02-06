# Contributing to PDF SDK

Thank you for your interest in contributing to PDF SDK! 🎉

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/infosec554/convert-pdf-go-sdk/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Go version and OS
   - Sample PDF (if relevant)

### Suggesting Features

1. Check existing issues for similar suggestions
2. Create a new issue with the `enhancement` label
3. Describe the feature and its use case

### Pull Requests

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Run tests: `make test`
5. Run linter: `make lint`
6. Commit with clear message: `git commit -m 'Add amazing feature'`
7. Push: `git push origin feature/amazing-feature`
8. Open a Pull Request

## Development Setup

```bash
# Clone the repository
git clone https://github.com/infosec554/convert-pdf-go-sdk.git
cd convert-pdf-go-sdk

# Install dependencies
make deps

# Run tests
make test

# Run with coverage
make cover

# Format code
make fmt

# Run linter
make lint
```

## Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Add comments for exported functions
- Keep functions small and focused
- Write tests for new features

## Testing

```bash
# Run all tests
make test

# Run with verbose output
go test -v . ./service/... ./pkg/...

# Run with coverage
make cover

# Run benchmarks
make bench

# Run specific test
go test -v -run TestWorkerPool .
```

## Commit Messages

Use clear, descriptive commit messages:

- `feat: add PDF/A conversion support`
- `fix: correct page extraction for multi-digit pages`
- `docs: update API documentation`
- `test: add tests for rate limiter`
- `refactor: simplify batch processor logic`

## Project Structure

```
convert-pdf-go-sdk/
├── pdfsdk.go          # Main entry point
├── errors.go          # Custom error types
├── retry.go           # Retry mechanism
├── ratelimit.go       # Rate limiter
├── metrics.go         # Metrics collection
├── pkg/
│   ├── gotenberg/     # Gotenberg client
│   └── logger/        # Logger
├── service/
│   ├── service.go     # Main service interface
│   ├── batch.go       # Batch processing
│   ├── advanced.go    # Advanced operations
│   └── *.go           # Individual services
├── cmd/
│   └── main.go        # Example usage
└── *_test.go          # Tests
```

## Need Help?

- Open an issue
- Contact: [@zarifjorayev](https://t.me/zarifjorayev)
- Email: infosec554@gmail.com

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
