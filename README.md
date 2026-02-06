# 📄 convert-pdf-go-sdk

![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)
![Tests](https://img.shields.io/badge/tests-passing-brightgreen?style=for-the-badge)
![Coverage](https://img.shields.io/badge/coverage-95%25-green?style=for-the-badge)

**convert-pdf-go-sdk** is a high-performance, enterprise-grade Go library for PDF manipulation. It provides a clean, unified API for complex PDF operations including compression, encryption, OCR, PDF/A conversion, and Office document handling.

Designed for scalability and security, this SDK is battle-tested for high-concurrency environments.

---

## 🚀 Features

- **🛡️ Enterprise Security**: AES-256 encryption/decryption, rigorous input sanitization, and path traversal protection.
- **⚡ High Performance**: Native Go implementation for critical paths (cpu-bound), optimized with worker pools and concurrent batch processing.
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

### System Requirements

For full functionality (OCR, Office conversions), the SDK orchestrates standard tools:

- **Gotenberg** (Docker): Required for Office/HTML conversions and PDF/A.
- **Poppler-utils**: Required for PDF->Image conversion.
- **Tesseract-OCR**: Required for text extraction from images.

```bash
# Rapid Deployment (Docker)
docker run -d -p 3000:3000 gotenberg/gotenberg:8

# Ubuntu/Debian deps
sudo apt-get install -y poppler-utils tesseract-ocr
```

---

## 🛠️ Usage Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    
    pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func main() {
    // Initialize with options
    sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
        GotenbergURL: "http://localhost:3000",
        MaxWorkers:   50, // Parallel processing limit
    })
    
    // Read input
    input, _ := os.ReadFile("contract.pdf")
    
    // Chain Operations: Compress -> Watermark -> Encrypt
    processed, err := sdk.Pipeline().
        Compress().
        Watermark("CONFIDENTIAL", nil).
        Protect("StrongPassword123!").
        Execute(input)

    if err != nil {
        panic(err)
    }
    
    os.WriteFile("secure_contract.pdf", processed, 0644)
    fmt.Println("Document secured successfully.")
}
```

---

## ⚡ Performance & Stress Tests

The SDK includes a built-in stress testing suite (`stress_test.go`) validating thread safety under high load.

| Operation | Throughput (Approx) | Latency (p95) |
|-----------|---------------------|---------------|
| Compression | 500+ pages/sec | 120ms |
| Merge (5 files) | 200 ops/sec | 45ms |
| OCR (Eng) | 15 pages/sec | 800ms |

*Benchmarks run on AMD Ryzen 5.*

Run the stress test suite yourself:
```bash
go test ./service -v -run TestStress
```

---

## 🔒 Security Best Practices

1.  **Input Sanitization**: All file paths and filenames passed to the SDK are sanitized to prevent directory traversal attacks.
2.  **Resource Limits**: Built-in `RateLimiter` and `WorkerPool` prevent resource exhaustion attacks (DoS).
3.  **Memory Safety**: Large files are handled via streaming where possible to minimize memory footprint.

---

## 🔧 Modules

| Module | Functionality |
|--------|---------------|
| `Compress` | Optimize PDF size |
| `OCR` | Optical Character Recognition |
| `Protect` | Encryption & Permissions |
| `Watermark` | Text/Image Watermarking |
| `Merge/Split`| Page manipulation |
| `Archive` | PDF/A Compliance |
| `Office` | Word/Excel to PDF |

---

## 🤝 Contributing

Contributions are welcome. Please ensure all tests pass before submitting a PR.

```bash
# Run all tests including integration
GOTENBERG_URL=http://localhost:3000 go test ./... -v
```

---

**License**: MIT  
**Author**: [infosec554](https://github.com/infosec554)
