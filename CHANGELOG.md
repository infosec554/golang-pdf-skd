# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0] - 2026-02-06

### Added
- **Worker Pool** - Control parallel operations with configurable worker limits
- **Connection Pool** - Optimized HTTP connections for Gotenberg
- **Buffer Pool** - Memory-efficient buffer reuse
- **Rate Limiter** - Token bucket rate limiting
- **Retry Mechanism** - Exponential backoff with jitter
- **Metrics** - SDK usage metrics with Prometheus format
- **Custom Errors** - Typed errors (ErrInvalidPDF, ErrTimeout, etc.)
- **Health Check** - Gotenberg server health monitoring
- **BatchProcessor** - Process multiple PDFs in parallel
- **Pipeline** - Chain multiple operations
- **InfoService** - PDF info, page count, validation
- **PageService** - Extract, delete, insert, reorder pages
- **TextService** - Extract text from PDF
- **MetadataService** - Get/set PDF metadata
- **ImageExtractService** - Extract images from PDF
- **Unit Tests** - Comprehensive test coverage
- **Makefile** - Build, test, lint commands
- **GitHub Actions** - CI/CD workflow

### Changed
- Updated version to 2.1.0
- Improved Gotenberg client with connection pooling
- Enhanced logging with Debug and Warn levels

## [2.0.0] - 2026-02-06

### Added
- **Options struct** - SDK configuration with custom options
- **NewWithOptions()** - Initialize SDK with custom settings
- **Debug/Warn logging** - Additional log levels

### Changed
- Updated version to 2.0.0
- Go version updated to 1.22
- Renamed `pkg/gotenberg/clint.go` to `client.go`

### Fixed
- PowerPoint extension detection
- Excel extension detection
- Compression ratio logging

### Removed
- Old database-related code (data/pg directory)

## [1.0.0] - 2026-02-05

### Added
- Initial release
- PDF compression
- PDF merging
- PDF splitting
- PDF rotation
- Watermark support
- Password protection
- Password removal (unlock)
- PDF to JPG conversion
- JPG to PDF conversion
- Word to PDF conversion (Gotenberg)
- Excel to PDF conversion (Gotenberg)
- PowerPoint to PDF conversion (Gotenberg)
- Zap logger integration
