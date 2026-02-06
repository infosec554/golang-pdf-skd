# 📄 Convert PDF Go SDK v2.3

A powerful, memory-efficient Go SDK for PDF operations with parallel processing support.

[![Go Reference](https://pkg.go.dev/badge/github.com/infosec554/convert-pdf-go-sdk.svg)](https://pkg.go.dev/github.com/infosec554/convert-pdf-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)

## 🆕 What's New in v2.3.0

- **👁️ OCR Service** - Extract text from scanned PDFs and create searchable PDFs (Tesseract)
- **🗄️ Archive Service** - Convert PDF to PDF/A for long-term archiving (PDF/A-1b, 2b, 3b)
- **📝 Form Service** - Fill PDF forms programmatically
- **📎 Attachment Service** - Add/Extract file attachments from PDF
- **🔄 Worker Pool** - Control parallel operations
- **📊 Batch Processing** - Process multiple PDFs in parallel
- **⛓️ Pipeline** - Chain multiple operations

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
- **Gotenberg** (for Word/Excel/PowerPoint & PDF/A)
- **poppler-utils** (for PDF to Image)
- **tesseract-ocr** (for OCR)

```bash
# Ubuntu/Debian
sudo apt-get install poppler-utils tesseract-ocr tesseract-ocr-eng

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
    "context"

    pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func main() {
    sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
        GotenbergURL:  "http://localhost:3000",
        MaxWorkers:    10,
    })
    defer sdk.Close()

    input, _ := os.ReadFile("document.pdf")
    ctx := context.Background()

    // 1. Get Info
    info, _ := sdk.Info().GetInfoBytes(input)
    fmt.Printf("Pages: %d\n", info.PageCount)

    // 2. Compress
    compressed, _ := sdk.Compress().CompressBytes(input)

    // 3. OCR (Extract Text from Scanned PDF)
    text, err := sdk.OCR().ExtractText(ctx, input, "eng")
    if err == nil {
        fmt.Println("Extracted Text:", text)
    }

    // 4. Create Searchable PDF (from scanned)
    searchable, _ := sdk.OCR().CreateSearchablePDF(ctx, input, "eng")
    os.WriteFile("searchable.pdf", searchable, 0644)
}
```

## 👁️ OCR (Optical Character Recognition)

Process scanned documents:

```go
// Extract text from scanned PDF
text, err := sdk.OCR().ExtractText(ctx, input, "eng")

// Convert scanned PDF to Searchable PDF (adds text layer)
searchableBytes, err := sdk.OCR().CreateSearchablePDF(ctx, input, "eng")
```

## 🗄️ Archive (PDF/A)

Convert documents for long-term preservation:

```go
// Convert to PDF/A-1b
pdfaBytes, err := sdk.Archive().ConvertToPDFA(input, "PDF/A-1b")
```

## 📝 Forms & Attachments

```go
// Fill Form
sdk.Form().FillForm(input, map[string]interface{}{"Name": "John"})

// Add Attachment
sdk.Attachment().AddAttachments(input, map[string][]byte{"note.txt": []byte("Hi")})
```

## ⛓️ Pipeline (Chained Operations)

```go
result, err := sdk.Pipeline().
    Compress().
    Watermark("CONFIDENTIAL", nil).
    Protect("secret123").
    Execute(input)
```

## 🔄 Batch Processing (Parallel)

```go
ctx := context.Background()
inputs := [][]byte{pdf1, pdf2, pdf3}
results := sdk.Batch(5).CompressBatch(ctx, inputs)
```

## 🔧 Available Services

| Service | Description | Parallel | Gotenberg |
|---------|-------------|:--------:|:---------:|
| `Info()` | PDF metadata, page count | ✅ | ❌ |
| `Compress()` | Optimize and compress | ✅ | ❌ |
| `Merge()` | Combine multiple PDFs | ✅ | ❌ |
| `Split()` | Split by page ranges | ✅ | ❌ |
| `Rotate()` | Rotate pages | ✅ | ❌ |
| `Watermark()` | Watermarking | ✅ | ❌ |
| `Protect()` | Password protection | ✅ | ❌ |
| `OCR()` | Extract Text / Searchable PDF | ✅ | ❌ |
| `Archive()` | PDF to PDF/A | ✅ | ✅ |
| `Form()` | Fill Forms | ✅ | ❌ |
| `WordToPDF()` | DOCX to PDF | ✅ | ✅ |

(Full list in code)

## 📧 Contact

- **Telegram:** [@zarifjorayev](https://t.me/zarifjorayev)
- **GitHub:** [@infosec554](https://github.com/infosec554)

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
