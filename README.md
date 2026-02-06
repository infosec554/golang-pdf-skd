# 📄 Convert PDF Go SDK v2.2

A powerful, memory-efficient Go SDK for PDF operations with parallel processing support.

[![Go Reference](https://pkg.go.dev/badge/github.com/infosec554/convert-pdf-go-sdk.svg)](https://pkg.go.dev/github.com/infosec554/convert-pdf-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)

## 🆕 What's New in v2.2.0

- **�️ Archive Service** - Convert PDF to PDF/A for long-term archiving (PDF/A-1b, 2b, 3b)
- **� Form Service** - Fill PDF forms programmatically (limited support)
- **� Attachment Service** - Add/Extract file attachments from PDF
- **🔄 Worker Pool** - Control parallel operations with configurable worker limits
- **📊 Batch Processing** - Process multiple PDFs in parallel
- **⛓️ Pipeline** - Chain multiple operations (compress → watermark → protect)
- **📋 PDF Info** - Get page count, version, encryption status
- **🔧 Connection Pooling** - Optimized HTTP connections for Gotenberg
- **💾 Buffer Pool** - Memory-efficient buffer reuse

## 📊 Performance & Limits

| Setting | Default | Description |
|---------|---------|-------------|
| `MaxWorkers` | 10 | Maximum parallel PDF operations |
| `MaxIdleConns` | 100 | HTTP connection pool size |
| `MaxConnsPerHost` | 100 | Max connections per Gotenberg host |
| `RequestTimeout` | 5 min | Request timeout |

## 🚀 Installation

```bash
go get github.com/infosec554/convert-pdf-go-sdk
```

## 📋 Requirements

- **Go 1.22+**
- **Gotenberg** (for Word/Excel/PowerPoint conversion & PDF/A)
- **pdftoppm** (for PDF to JPG)

```bash
# Ubuntu/Debian
sudo apt-get install poppler-utils

# Run Gotenberg with Docker
docker run -d -p 3000:3000 gotenberg/gotenberg:8
```

## 💡 Quick Start

```go
package main

import (
    "fmt"
    "os"
    "time"

    pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func main() {
    // Initialize with optimized settings
    sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
        GotenbergURL:  "http://localhost:3000",
        MaxWorkers:    10,   // Max 10 parallel operations
        MaxIdleConns:  100,  // Connection pool
        RequestTimeout: 5 * time.Minute,
    })
    defer sdk.Close() // Clean up connections

    input, _ := os.ReadFile("document.pdf")

    // Get PDF info
    info, _ := sdk.Info().GetInfoBytes(input)
    fmt.Printf("Pages: %d, Encrypted: %v\n", info.PageCount, info.Encrypted)

    // Compress
    output, _ := sdk.Compress().CompressBytes(input)
    fmt.Printf("Size: %d -> %d bytes\n", len(input), len(output))
}
```

## 🗄️ Archive (PDF/A)

Convert documents for long-term preservation:

```go
// Convert to PDF/A-1b
pdfaBytes, err := sdk.Archive().ConvertToPDFA(input, "PDF/A-1b")
```

## 📝 Forms

Fill PDF forms:

```go
data := map[string]interface{}{
    "Name": "John Doe",
    "Age":  30,
}
filledBytes, err := sdk.Form().FillForm(input, data)
```

## � Attachments

Add attachments to PDF:

```go
files := map[string][]byte{
    "invoice.xml": xmlBytes,
    "notes.txt":   txtBytes,
}
result, err := sdk.Attachment().AddAttachments(input, files)
```

## ⛓️ Pipeline (Chained Operations)

Execute multiple operations in sequence:

```go
// Compress → Add Watermark → Protect with password
result, err := sdk.Pipeline().
    Compress().
    Watermark("CONFIDENTIAL", nil).
    Protect("secret123").
    Execute(input)
```

## 🔄 Batch Processing (Parallel)

Process multiple PDFs concurrently:

```go
ctx := context.Background()

// Compress 100 PDFs with max 5 workers
inputs := [][]byte{pdf1, pdf2, pdf3, /* ... */}
results := sdk.Batch(5).CompressBatch(ctx, inputs)

for _, r := range results {
    if r.Error != nil {
        log.Printf("PDF %d failed: %v", r.Index, r.Error)
    } else {
        // Use r.Data
    }
}
```

## 🔧 Available Services

| Service | Description | Parallel | Gotenberg |
|---------|-------------|:--------:|:---------:|
| `Info()` | PDF metadata, page count, validation | ✅ | ❌ |
| `Compress()` | Optimize and compress | ✅ | ❌ |
| `Merge()` | Combine multiple PDFs | ✅ | ❌ |
| `Split()` | Split by page ranges | ✅ | ❌ |
| `Rotate()` | Rotate pages | ✅ | ❌ |
| `Watermark()` | Add text watermarks | ✅ | ❌ |
| `Protect()` | Password protection | ✅ | ❌ |
| `Unlock()` | Remove passwords | ✅ | ❌ |
| `PDFToJPG()` | Convert to images | ✅ | ❌ |
| `JPGToPDF()` | Images to PDF | ✅ | ❌ |
| `WordToPDF()` | DOCX to PDF | ✅ | ✅ |
| `ExcelToPDF()` | XLSX to PDF | ✅ | ✅ |
| `PowerPointToPDF()` | PPTX to PDF | ✅ | ✅ |
| `Archive()` | PDF to PDF/A | ✅ | ✅ |
| `Form()` | Fill Forms | ✅ | ❌ |
| `Attachment()` | Add/Extract Attachments | ✅ | ❌ |
| `Batch()` | Parallel batch ops | ✅ | - |
| `Pipeline()` | Chained operations | ✅ | - |

## 📧 Contact

- **Telegram:** [@zarifjorayev](https://t.me/zarifjorayev)
- **Email:** infosec554@gmail.com
- **GitHub:** [@infosec554](https://github.com/infosec554)

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
