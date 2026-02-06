# 📄 Convert PDF Go SDK v2.1

A powerful, memory-efficient Go SDK for PDF operations with parallel processing support.

[![Go Reference](https://pkg.go.dev/badge/github.com/infosec554/convert-pdf-go-sdk.svg)](https://pkg.go.dev/github.com/infosec554/convert-pdf-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)

## 🆕 What's New in v2.1.0

- **🔄 Worker Pool** - Control parallel operations with configurable worker limits
- **📊 Batch Processing** - Process multiple PDFs in parallel
- **⛓️ Pipeline** - Chain multiple operations (compress → watermark → protect)
- **📋 PDF Info** - Get page count, version, encryption status
- **🔧 Connection Pooling** - Optimized HTTP connections for Gotenberg
- **💾 Buffer Pool** - Memory-efficient buffer reuse
- **📈 Statistics** - Track processed tasks and worker usage

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
- **Gotenberg** (for Word/Excel/PowerPoint conversion)
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

### Available Batch Operations:

```go
batch := sdk.Batch(maxWorkers)

batch.CompressBatch(ctx, inputs)
batch.RotateBatch(ctx, inputs, 90, "all")
batch.WatermarkBatch(ctx, inputs, "DRAFT", nil)
batch.ProtectBatch(ctx, inputs, "password")
```

## 👷 Worker Pool Control

For fine-grained control over concurrency:

```go
// Get worker slot (blocks if all busy)
sdk.Workers().Acquire()
defer sdk.Workers().Release()

// Non-blocking attempt
if sdk.Workers().TryAcquire() {
    defer sdk.Workers().Release()
    // Do work
}

// Check availability
available := sdk.Workers().Available() // e.g., 8 of 10
```

## 📊 Statistics

```go
stats := sdk.Stats()
fmt.Printf("Active: %d/%d workers, Processed: %d tasks\n",
    stats.ActiveWorkers,
    stats.MaxWorkers,
    stats.ProcessedTasks)
```

## 📚 API Reference

### Initialization

```go
// Simple
sdk := pdfsdk.New("http://localhost:3000")

// With options (recommended)
sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
    GotenbergURL:        "http://localhost:3000",
    MaxWorkers:          10,
    MaxIdleConns:        100,
    MaxIdleConnsPerHost: 10,
    MaxConnsPerHost:     100,
    IdleConnTimeout:     90 * time.Second,
    RequestTimeout:      5 * time.Minute,
})
defer sdk.Close()
```

### PDF Info

```go
info, _ := sdk.Info().GetInfoBytes(input)
// info.PageCount, info.Version, info.Encrypted, info.FileSize

pageCount, _ := sdk.Info().GetPageCount(input)
err := sdk.Info().ValidatePDF(input)
isEncrypted, _ := sdk.Info().IsEncrypted(input)
```

### Compress

```go
output, _ := sdk.Compress().CompressBytes(input)
sdk.Compress().CompressFile("input.pdf", "output.pdf")
```

### Merge

```go
output, _ := sdk.Merge().MergeBytes([][]byte{pdf1, pdf2, pdf3})
sdk.Merge().MergeFiles([]string{"1.pdf", "2.pdf"}, "merged.pdf")
```

### Split

```go
zipBytes, _ := sdk.Split().SplitBytes(input, "1-5,6-10")
pages, _ := sdk.Split().SplitToPages(input)
```

### Rotate

```go
output, _ := sdk.Rotate().RotateBytes(input, 90, "all")  // 90, 180, 270
sdk.Rotate().RotateFile("in.pdf", "out.pdf", 180, "1-3")
```

### Watermark

```go
// Default options
output, _ := sdk.Watermark().AddWatermarkBytes(input, "CONFIDENTIAL", nil)

// Custom options
output, _ := sdk.Watermark().AddWatermarkBytes(input, "DRAFT", &service.WatermarkOptions{
    FontSize: 72,
    Position: "diagonal",
    Opacity:  0.5,
    Color:    "red",
})
```

### Protect & Unlock

```go
protected, _ := sdk.Protect().ProtectBytes(input, "password")
unlocked, _ := sdk.Unlock().UnlockBytes(protected, "password")
```

### PDF to JPG

```go
zipBytes, _ := sdk.PDFToJPG().ConvertBytes(input)
images, _ := sdk.PDFToJPG().ConvertToImages(input) // [][]byte
```

### Images to PDF

```go
pdfBytes, _ := sdk.JPGToPDF().ConvertMultipleBytes(images, filenames)
```

### Document Conversion (requires Gotenberg)

```go
ctx := context.Background()
pdfBytes, _ := sdk.WordToPDF().ConvertBytes(ctx, docxBytes, "doc.docx")
pdfBytes, _ := sdk.ExcelToPDF().ConvertBytes(ctx, xlsxBytes, "sheet.xlsx")
pdfBytes, _ := sdk.PowerPointToPDF().ConvertBytes(ctx, pptxBytes, "slides.pptx")
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
| `Batch()` | Parallel batch ops | ✅ | - |
| `Pipeline()` | Chained operations | ✅ | - |

## 🐳 Docker

```bash
# Run with Docker Compose
docker-compose up -d

# Or run Gotenberg separately
docker run -d -p 3000:3000 gotenberg/gotenberg:8
```

## 📦 Project Structure

```
convert-pdf-go-sdk/
├── pdfsdk.go              # Main SDK with WorkerPool
├── Dockerfile
├── docker-compose.yml
├── config/
├── pkg/
│   ├── gotenberg/
│   │   └── client.go      # Gotenberg with connection pool
│   └── logger/
└── service/
    ├── service.go         # Main service interface
    ├── batch.go           # BatchProcessor & Pipeline
    ├── info.go            # PDF info & validation
    ├── compress.go
    ├── merge_service.go
    ├── split.go
    ├── rotate.go
    ├── watermark.go
    ├── protect.go
    ├── unlock.go
    ├── pdf_to_jpg.go
    ├── jpgtopdf.go
    ├── word_to_pdf.go
    ├── excel_to_pdf.go
    └── powerpoint_to_pdf.go
```

## 📧 Contact

- **Telegram:** [@zarifjorayev](https://t.me/zarifjorayev)
- **Email:** infosec554@gmail.com
- **GitHub:** [@infosec554](https://github.com/infosec554)

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
