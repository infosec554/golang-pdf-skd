# 📄 Golang PDF SDK

A powerful and easy-to-use Go SDK for PDF operations including compression, merging, splitting, rotation, watermarking, and more.

[![Go Reference](https://pkg.go.dev/badge/github.com/infosec554/golang-pdf-sdk.svg)](https://pkg.go.dev/github.com/infosec554/golang-pdf-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🚀 Installation

```bash
go get github.com/infosec554/golang-pdf-sdk
```

## 📋 Requirements

- **Go 1.21+**
- **Gotenberg** (for Word/Excel/PowerPoint conversion) - [gotenberg.dev](https://gotenberg.dev)
- **pdftoppm** (for PDF to JPG) - part of poppler-utils

```bash
# Ubuntu/Debian
sudo apt-get install poppler-utils

# macOS
brew install poppler

# Run Gotenberg with Docker
docker run -d -p 3000:3000 gotenberg/gotenberg:8
```

## 💡 Quick Start

```go
package main

import (
    "fmt"
    "os"

    pdfsdk "github.com/infosec554/golang-pdf-sdk"
)

func main() {
    // Initialize SDK
    pdf := pdfsdk.New("http://localhost:3000")

    // Read input PDF
    input, _ := os.ReadFile("document.pdf")

    // Compress the PDF
    output, err := pdf.Compress().CompressBytes(input)
    if err != nil {
        panic(err)
    }

    // Save result
    os.WriteFile("compressed.pdf", output, 0644)

    fmt.Printf("Size: %d -> %d bytes\n", len(input), len(output))
}
```

---

## 📚 API Reference

### 1️⃣ Compress PDF

Reduce PDF file size:

```go
pdf := pdfsdk.New("http://localhost:3000")

// From bytes
input, _ := os.ReadFile("large.pdf")
output, err := pdf.Compress().CompressBytes(input)

// From file
err := pdf.Compress().CompressFile("input.pdf", "output.pdf")
```

---

### 2️⃣ Merge PDFs

Combine multiple PDFs into one:

```go
pdf := pdfsdk.New("http://localhost:3000")

// From bytes
pdf1, _ := os.ReadFile("doc1.pdf")
pdf2, _ := os.ReadFile("doc2.pdf")
pdf3, _ := os.ReadFile("doc3.pdf")

output, err := pdf.Merge().MergeBytes([][]byte{pdf1, pdf2, pdf3})
os.WriteFile("merged.pdf", output, 0644)

// From files
err := pdf.Merge().MergeFiles(
    []string{"doc1.pdf", "doc2.pdf", "doc3.pdf"},
    "merged.pdf",
)
```

---

### 3️⃣ Split PDF

Extract pages from PDF:

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("book.pdf")

// Split by page ranges (returns ZIP)
zipBytes, err := pdf.Split().SplitBytes(input, "1-5,6-10,11-20")
os.WriteFile("parts.zip", zipBytes, 0644)

// Split into individual pages
pages, err := pdf.Split().SplitToPages(input)
for i, page := range pages {
    os.WriteFile(fmt.Sprintf("page_%d.pdf", i+1), page, 0644)
}
```

---

### 4️⃣ Rotate PDF

Rotate pages by 90°, 180°, or 270°:

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("document.pdf")

// Rotate all pages 90°
output, err := pdf.Rotate().RotateBytes(input, 90, "all")

// Rotate specific pages 180°
output, err := pdf.Rotate().RotateBytes(input, 180, "1-3")

// From file
err := pdf.Rotate().RotateFile("input.pdf", "output.pdf", 270, "all")
```

---

### 5️⃣ Add Watermark

Add text watermark to PDF:

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("document.pdf")

// Simple watermark
output, err := pdf.Watermark().AddWatermarkBytes(input, "CONFIDENTIAL", nil)

// With custom options
options := &service.WatermarkOptions{
    FontSize: 72,           // Font size
    Position: "diagonal",   // "diagonal" or "center"
    Opacity:  0.3,          // 0.0 to 1.0
    Color:    "red",        // Color name
}
output, err := pdf.Watermark().AddWatermarkBytes(input, "DRAFT", options)
```

---

### 6️⃣ Protect PDF (Add Password)

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("secret.pdf")

// Add password protection
protected, err := pdf.Protect().ProtectBytes(input, "mypassword123")
os.WriteFile("protected.pdf", protected, 0644)

// From file
err := pdf.Protect().ProtectFile("input.pdf", "protected.pdf", "mypassword123")
```

---

### 7️⃣ Unlock PDF (Remove Password)

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("protected.pdf")

// Remove password
unlocked, err := pdf.Unlock().UnlockBytes(input, "mypassword123")
os.WriteFile("unlocked.pdf", unlocked, 0644)
```

---

### 8️⃣ PDF to Images (JPG)

Convert PDF pages to JPG images:

```go
pdf := pdfsdk.New("http://localhost:3000")
input, _ := os.ReadFile("presentation.pdf")

// Get ZIP file containing all images
zipBytes, err := pdf.PDFToJPG().ConvertBytes(input)
os.WriteFile("images.zip", zipBytes, 0644)

// Get individual images
images, err := pdf.PDFToJPG().ConvertToImages(input)
for i, img := range images {
    os.WriteFile(fmt.Sprintf("page_%d.jpg", i+1), img, 0644)
}
```

---

### 9️⃣ Images to PDF

Combine JPG/PNG images into a single PDF:

```go
pdf := pdfsdk.New("http://localhost:3000")

// Read images
img1, _ := os.ReadFile("photo1.jpg")
img2, _ := os.ReadFile("photo2.png")
img3, _ := os.ReadFile("photo3.jpg")

// Create PDF
pdfBytes, err := pdf.JPGToPDF().ConvertMultipleBytes(
    [][]byte{img1, img2, img3},
    []string{"photo1.jpg", "photo2.png", "photo3.jpg"},
)
os.WriteFile("album.pdf", pdfBytes, 0644)

// From files
err := pdf.JPGToPDF().ConvertFiles(
    []string{"1.jpg", "2.jpg", "3.png"},
    "album.pdf",
)
```

---

### 🔟 Document Conversion (Requires Gotenberg)

Convert Word, Excel, PowerPoint to PDF:

```go
import "context"

pdf := pdfsdk.New("http://localhost:3000")
ctx := context.Background()

// Word to PDF
docxBytes, _ := os.ReadFile("document.docx")
pdfBytes, err := pdf.WordToPDF().ConvertBytes(ctx, docxBytes, "document.docx")

// Excel to PDF
xlsxBytes, _ := os.ReadFile("spreadsheet.xlsx")
pdfBytes, err := pdf.ExcelToPDF().ConvertBytes(ctx, xlsxBytes, "spreadsheet.xlsx")

// PowerPoint to PDF
pptxBytes, _ := os.ReadFile("slides.pptx")
pdfBytes, err := pdf.PowerPointToPDF().ConvertBytes(ctx, pptxBytes, "slides.pptx")
```

---

## 🔧 Configuration

Using `.env` file:

```env
GOTENBERG_URL=http://localhost:3000
SERVICE_NAME=my-pdf-app
LOGGER_LEVEL=info
```

---

## 📦 Project Structure

```
golang-pdf-sdk/
├── pdfsdk.go              # Main SDK entry point
├── config/config.go       # Configuration
├── pkg/
│   ├── gotenberg/         # Gotenberg HTTP client
│   └── logger/            # Logging utilities
└── service/
    ├── compress.go        # PDF compression
    ├── merge_service.go   # Merge multiple PDFs
    ├── split.go           # Split PDF by pages
    ├── rotate.go          # Rotate PDF pages
    ├── watermark.go       # Add text watermark
    ├── protect.go         # Password protect PDF
    ├── unlock.go          # Remove password
    ├── pdf_to_jpg.go      # PDF to JPG images
    ├── jpgtopdf.go        # Images to PDF
    ├── word_to_pdf.go     # Word to PDF
    ├── excel_to_pdf.go    # Excel to PDF
    └── powerpoint_to_pdf.go # PowerPoint to PDF
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Pull requests are welcome!
