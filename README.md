# 📄 Convert PDF Go SDK

A powerful Go SDK for PDF operations - compression, merging, splitting, rotation, watermarking, and document conversion.

[![Go Reference](https://pkg.go.dev/badge/github.com/infosec554/convert-pdf-go-sdk.svg)](https://pkg.go.dev/github.com/infosec554/convert-pdf-go-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🚀 Installation

```bash
go get github.com/infosec554/convert-pdf-go-sdk
```

## 📋 Requirements

- **Go 1.21+**
- **Gotenberg** (for Word/Excel/PowerPoint conversion)
- **pdftoppm** (for PDF to JPG)

```bash
# Ubuntu/Debian
sudo apt-get install poppler-utils

# macOS
brew install poppler

# Run with Docker Compose
docker-compose up -d
```

## 💡 Quick Start

```go
package main

import (
    "fmt"
    "os"

    pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func main() {
    pdf := pdfsdk.New("http://localhost:3000")

    input, _ := os.ReadFile("document.pdf")

    output, _ := pdf.Compress().CompressBytes(input)

    os.WriteFile("compressed.pdf", output, 0644)

    fmt.Printf("Size: %d -> %d bytes\n", len(input), len(output))
}
```

## 📚 API Reference

### Compress
```go
output, _ := pdf.Compress().CompressBytes(input)
pdf.Compress().CompressFile("input.pdf", "output.pdf")
```

### Merge
```go
output, _ := pdf.Merge().MergeBytes([][]byte{pdf1, pdf2, pdf3})
pdf.Merge().MergeFiles([]string{"1.pdf", "2.pdf"}, "merged.pdf")
```

### Split
```go
zipBytes, _ := pdf.Split().SplitBytes(input, "1-5,6-10")
pages, _ := pdf.Split().SplitToPages(input)
```

### Rotate
```go
output, _ := pdf.Rotate().RotateBytes(input, 90, "all")
pdf.Rotate().RotateFile("input.pdf", "output.pdf", 180, "1-3")
```

### Watermark
```go
output, _ := pdf.Watermark().AddWatermarkBytes(input, "CONFIDENTIAL", nil)
```

### Protect & Unlock
```go
protected, _ := pdf.Protect().ProtectBytes(input, "password")
unlocked, _ := pdf.Unlock().UnlockBytes(protected, "password")
```

### PDF to JPG
```go
zipBytes, _ := pdf.PDFToJPG().ConvertBytes(input)
images, _ := pdf.PDFToJPG().ConvertToImages(input)
```

### Images to PDF
```go
pdfBytes, _ := pdf.JPGToPDF().ConvertMultipleBytes(images, filenames)
```

### Document Conversion (requires Gotenberg)
```go
pdfBytes, _ := pdf.WordToPDF().ConvertBytes(ctx, docxBytes, "doc.docx")
pdfBytes, _ := pdf.ExcelToPDF().ConvertBytes(ctx, xlsxBytes, "sheet.xlsx")
pdfBytes, _ := pdf.PowerPointToPDF().ConvertBytes(ctx, pptxBytes, "slides.pptx")
```

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
├── pdfsdk.go
├── Dockerfile
├── docker-compose.yml
├── config/
├── pkg/
│   ├── gotenberg/
│   └── logger/
└── service/
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

## 🤝 Contributing

Pull requests are welcome! Please open an issue first to discuss what you would like to change.
