# golang-pdf-sdk

A powerful Go SDK for PDF operations including conversion, manipulation, and processing.

## Features

- **PDF to Image Conversion** - Convert PDF pages to JPG/PNG images
- **Image to PDF** - Convert images to PDF documents
- **PDF Manipulation** - Merge, split, and process PDF files
- **High Performance** - Built with Gotenberg for reliable PDF processing

## Installation

```bash
go get github.com/infosec554/golang-pdf-sdk
```

## Quick Start

```go
package main

import (
    "github.com/infosec554/golang-pdf-sdk/service"
    "github.com/infosec554/golang-pdf-sdk/pkg/gotenberg"
)

func main() {
    // Initialize Gotenberg client
    client := gotenberg.New("http://localhost:3000")
    
    // Use PDF services
    // ... 
}
```

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/pdf/convert` | POST | Convert documents to PDF |
| `/pdf/to-jpg` | POST | Convert PDF to JPG images |
| `/pdf/merge` | POST | Merge multiple PDFs |

## License

MIT License - see [LICENSE](LICENSE) for details.
