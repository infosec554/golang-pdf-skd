# Golang PDF SDK

A powerful and easy-to-use Go SDK for PDF operations including conversion, manipulation, and processing.

[![Go Reference](https://pkg.go.dev/badge/github.com/infosec554/golang-pdf-sdk.svg)](https://pkg.go.dev/github.com/infosec554/golang-pdf-sdk)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

### 📄 Document Conversion
- **Word to PDF** - Convert .docx, .doc files to PDF
- **Excel to PDF** - Convert .xlsx, .xls spreadsheets to PDF
- **PowerPoint to PDF** - Convert .pptx, .ppt presentations to PDF

### 🖼️ Image Operations
- **PDF to JPG** - Extract PDF pages as JPG images
- **JPG/PNG to PDF** - Combine images into a PDF document

### ⚙️ PDF Manipulation
- **Compress** - Reduce PDF file size with optimization
- **Merge** - Combine multiple PDFs into one
- **Split** - Extract page ranges from a PDF
- **Rotate** - Rotate pages by 90°, 180°, or 270°

### 🔒 Security
- **Watermark** - Add text watermarks to PDFs
- **Protect** - Add password encryption to PDFs
- **Unlock** - Remove password protection from PDFs

## Installation

```bash
go get github.com/infosec554/golang-pdf-sdk
```

## Requirements

- Go 1.21 or higher
- **For document conversion (Word/Excel/PowerPoint):** [Gotenberg](https://gotenberg.dev/) running locally or remotely
- **For PDF to JPG:** `pdftoppm` (part of poppler-utils)

```bash
# Install pdftoppm on Ubuntu/Debian
sudo apt-get install poppler-utils

# Install pdftoppm on macOS
brew install poppler
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"
    
    pdfsdk "github.com/infosec554/golang-pdf-sdk"
)

func main() {
    // Create PDF service
    pdf := pdfsdk.New("http://localhost:3000") // Gotenberg URL
    
    // Read input PDF
    input, _ := os.ReadFile("input.pdf")
    
    // Compress the PDF
    output, err := pdf.Compress().CompressBytes(input)
    if err != nil {
        panic(err)
    }
    
    // Save compressed PDF
    os.WriteFile("compressed.pdf", output, 0644)
    
    fmt.Printf("Compressed: %d -> %d bytes\n", len(input), len(output))
}
```

## API Reference

### Compression

```go
// Compress PDF bytes
output, err := pdf.Compress().CompressBytes(input)

// Compress PDF file
err := pdf.Compress().CompressFile("input.pdf", "output.pdf")
```

### Merge

```go
// Merge multiple PDF byte slices
output, err := pdf.Merge().MergeBytes([][]byte{pdf1, pdf2, pdf3})

// Merge PDF files
err := pdf.Merge().MergeFiles([]string{"a.pdf", "b.pdf"}, "merged.pdf")
```

### Split

```go
// Split by page ranges (returns ZIP)
zipBytes, err := pdf.Split().SplitBytes(input, "1-3,4-6,7")

// Split into individual pages
pages, err := pdf.Split().SplitToPages(input)
```

### Rotate

```go
// Rotate all pages 90 degrees
output, err := pdf.Rotate().RotateBytes(input, 90, "all")

// Rotate specific pages
output, err := pdf.Rotate().RotateBytes(input, 180, "1-3")
```

### Watermark

```go
// Add watermark with default options
output, err := pdf.Watermark().AddWatermarkBytes(input, "CONFIDENTIAL", nil)

// Add watermark with custom options
options := &service.WatermarkOptions{
    FontSize: 72,
    Position: "diagonal",
    Opacity:  0.5,
    Color:    "red",
}
output, err := pdf.Watermark().AddWatermarkBytes(input, "DRAFT", options)
```

### Protect & Unlock

```go
// Add password protection
protected, err := pdf.Protect().ProtectBytes(input, "secretpassword")

// Remove password protection
unlocked, err := pdf.Unlock().UnlockBytes(protected, "secretpassword")
```

### PDF to JPG

```go
// Convert to ZIP containing JPG images
zipBytes, err := pdf.PDFToJPG().ConvertBytes(input)

// Get individual images
images, err := pdf.PDFToJPG().ConvertToImages(input)
for i, img := range images {
    os.WriteFile(fmt.Sprintf("page_%d.jpg", i+1), img, 0644)
}
```

### Images to PDF

```go
// Convert multiple images to PDF
pdfBytes, err := pdf.JPGToPDF().ConvertMultipleBytes(
    [][]byte{image1, image2, image3},
    []string{"photo1.jpg", "photo2.jpg", "photo3.png"},
)
```

### Document Conversion (requires Gotenberg)

```go
ctx := context.Background()

// Word to PDF
pdfBytes, err := pdf.WordToPDF().ConvertBytes(ctx, docxBytes, "document.docx")

// Excel to PDF
pdfBytes, err := pdf.ExcelToPDF().ConvertBytes(ctx, xlsxBytes, "spreadsheet.xlsx")

// PowerPoint to PDF
pdfBytes, err := pdf.PowerPointToPDF().ConvertBytes(ctx, pptxBytes, "presentation.pptx")
```

## Running Gotenberg

For document conversion features, you need Gotenberg running:

```bash
docker run -d -p 3000:3000 gotenberg/gotenberg:8
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
