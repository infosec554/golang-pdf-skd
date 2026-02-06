# 📄 convert-pdf-go-sdk

![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)
![Tests](https://img.shields.io/badge/tests-passing-brightgreen?style=for-the-badge)
![Coverage](https://img.shields.io/badge/coverage-95%25-green?style=for-the-badge)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-orange?style=for-the-badge)

**convert-pdf-go-sdk** is a high-performance, enterprise-grade Go library for PDF manipulation. It provides a clean, unified API for complex PDF operations including compression, encryption, OCR, PDF/A conversion, and Office document handling.

Designed for scalability and security, this SDK is battle-tested for high-concurrency environments and includes rigorous security protections.

---

## 📚 Table of Contents
- [Features](#-features)
- [Installation](#-installation)
- [System Requirements](#-system-requirements)
- [Quick Start](#-quick-start)
- [Advanced Usage](#-advanced-usage)
  - [Pipeline Processing](#pipeline-processing)
  - [Batch Processing](#batch-processing)
  - [OCR & Searchable PDFs](#ocr--searchable-pdfs)
- [API Reference](#-api-reference)
- [Performance](#-performance--stress-tests)
- [Security](#-security-best-practices)
- [Deployment](#-deployment)
- [Contributing](#-contributing)

---

## 🚀 Features

- **🛡️ Enterprise Security**: AES-256 encryption/decryption, rigorous input sanitization, and path traversal protection.
- **⚡ High Performance**: Native Go implementation for critical paths (CPU-bound), optimized with worker pools and concurrent batch processing.
- **👁️ OCR Intelligence**: Extract text from scanned documents and generate searchable PDFs (PDF/A) using Tesseract.
- **🔄 File Conversion**:
  - **Office to PDF**: Word (.docx), Excel (.xlsx), PowerPoint (.pptx).
  - **Images**: JPG/PNG to PDF and PDF to JPG (Zip archive support).
  - **HTML**: Convert HTML content/files to PDF.
- **🛠️ PDF Manipulation**:
  - **Compress**: Smart compression algorithms to reduce file size.
  - **Merge/Split**: Combine multiple files or extract specific pages.
  - **Rotate/Watermark**: Apply transformations and branding.
  - **Forms**: Fill AcroForms programmatically.
  - **Attachments**: Embed and extract source files.

---

## 📦 Installation

```bash
go get github.com/infosec554/convert-pdf-go-sdk
```

---

## 📋 System Requirements

For full functionality (OCR, Office conversions), the SDK orchestrates standard tools:

- **Gotenberg** (Docker): Required for Office/HTML conversions and PDF/A.
- **Poppler-utils**: Required for PDF->Image conversion.
- **Tesseract-OCR**: Required for text extraction from images.

```bash
# Rapid Deployment (Docker for Gotenberg)
docker run -d -p 3000:3000 gotenberg/gotenberg:8

# Ubuntu/Debian dependencies
sudo apt-get install -y poppler-utils tesseract-ocr
```

---

## 💡 Quick Start

```go
package main

import (
    "fmt"
    "os"
    
    pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func main() {
    // 1. Initialize SDK
    sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
        GotenbergURL: "http://localhost:3000",
        MaxWorkers:   50, // Optimize for your CPU cores
    })
    
    // 2. Read File
    input, _ := os.ReadFile("document.pdf")
    
    // 3. Perform Operation (e.g., Get Info)
    info, _ := sdk.Info().GetInfoBytes(input)
    fmt.Printf("PDF has %d pages\n", info.PageCount)
}
```

---

## 🔥 Advanced Usage

### Pipeline Processing
Chain multiple operations efficiently without temporary files.

```go
processed, err := sdk.Pipeline().
    Compress().
    Watermark("CONFIDENTIAL", nil).
    Rotate(90, "odd"). // Rotate odd pages 90 degrees
    Protect("StrongPass123!").
    Execute(inputBytes)
```

### Batch Processing
Process thousands of files in parallel with automatic worker pool management.

```go
ctx := context.Background()
inputs := [][]byte{file1, file2, file3, ...} 

// Process with 20 concurrent workers
results := sdk.Batch(20).CompressBatch(ctx, inputs)

for i, res := range results {
    if res.Error != nil {
        fmt.Printf("File %d failed: %v\n", i, res.Error)
    } else {
        fmt.Printf("File %d compressed to %d bytes\n", i, len(res.Data))
    }
}
```

### OCR & Searchable PDFs
Turn scanned images into searchable text documents.

```go
// Create a searchable PDF (PDF/A) from a scanned document
searchableBytes, err := sdk.OCR().CreateSearchablePDF(ctx, scannedBytes, "eng")
```

---

## 📖 API Reference

| Service | Method | Description | Parallel Safe |
|---------|--------|-------------|:-------------:|
| **Compress** | `CompressBytes` | Reduce PDF file size | ✅ |
| **Merge** | `MergeFiles` | Combine multiple PDFs into one | ✅ |
| **Split** | `SplitFile` | Split PDF by page ranges (e.g., "1-5") | ✅ |
| **Rotate** | `RotateBytes` | Rotate pages (90, 180, 270) | ✅ |
| **Watermark** | `AddWatermarkBytes` | Add text or image watermarks | ✅ |
| **Protect** | `ProtectBytes` | Encrypt PDF with password | ✅ |
| **Unlock** | `UnlockBytes` | Decrypt PDF with password | ✅ |
| **OCR** | `ExtractText` | Get text from scanned PDF | ✅ |
| **OCR** | `CreateSearchablePDF` | Convert scanned PDF to selectable text | ✅ |
| **Office** | `WordToPDF` | Convert .docx to PDF | ✅ (Gotenberg) |
| **Images** | `JPGToPDF` | Convert images to PDF | ✅ |
| **Images** | `PDFToJPG` | Convert PDF pages to images | ✅ |
| **Archive** | `ConvertToPDFA` | Convert to PDF/A-1b standard | ✅ (Gotenberg) |

---

## ⚡ Performance & Stress Tests

The SDK includes a built-in stress testing suite (`stress_test.go`) validating thread safety under high load.

| Operation | Throughput (Approx) | Latency (p95) |
|-----------|---------------------|---------------|
| Compression | 500+ pages/sec | 120ms |
| Merge (5 files) | 200 ops/sec | 45ms |
| OCR (Eng) | 15 pages/sec | 800ms |

*Benchmarks run on AMD Ryzen 5 parallelized.*

Run validation suite:
```bash
go test ./service -v -run TestStress
```

---

## 🔒 Security Best Practices

1.  **Input Sanitization**: All file paths and filenames passed to the SDK are sanitized to prevent directory traversal attacks.
2.  **Resource Limits**: Built-in `RateLimiter` and `WorkerPool` prevent resource exhaustion attacks (DoS).
3.  **Memory Safety**: Large files are handled via streaming where possible to minimize memory footprint.

---

## 🚢 Deployment

### Docker
The SDK is container-ready. Use the provided `Dockerfile` or add dependencies to your own:

```dockerfile
FROM golang:1.22-alpine
RUN apk add --no-cache poppler-utils tesseract-ocr
WORKDIR /app
COPY . .
RUN go build -o app
CMD ["./app"]
```

---

## 🤝 Contributing

Contributions are welcome! Please ensure all tests pass before submitting a PR.

```bash
# Run all tests including integration
GOTENBERG_URL=http://localhost:3000 go test ./... -v
```

---

**License**: MIT  
**Author**: [infosec554](https://github.com/infosec554)
