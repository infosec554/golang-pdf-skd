// Package pdfsdk provides a comprehensive Go SDK for PDF operations.
//
// This SDK allows you to:
//   - Convert documents (Word, Excel, PowerPoint) to PDF
//   - Convert PDFs to images (JPG)
//   - Convert images (JPG, PNG) to PDF
//   - Compress PDFs to reduce file size
//   - Merge multiple PDFs into one
//   - Split PDFs by page ranges
//   - Rotate PDF pages
//   - Add text watermarks to PDFs
//   - Protect PDFs with password encryption
//   - Unlock password-protected PDFs
//
// Quick Start:
//
//	import "github.com/infosec554/golang-pdf-sdk/service"
//
//	// Create a new PDF service
//	pdf := service.NewWithGotenberg("http://localhost:3000")
//
//	// Compress a PDF
//	output, err := pdf.Compress().CompressBytes(inputBytes)
//
//	// Rotate pages
//	output, err := pdf.Rotate().RotateBytes(inputBytes, 90, "all")
//
//	// Add watermark
//	output, err := pdf.Watermark().AddWatermarkBytes(inputBytes, "CONFIDENTIAL", nil)
//
// Requirements:
//
// Some operations require Gotenberg to be running:
//   - Word/Excel/PowerPoint to PDF conversion
//
// Local PDF operations (no Gotenberg required):
//   - Compress, Merge, Split, Rotate, Watermark, Protect, Unlock
//   - PDF to JPG conversion (requires pdftoppm)
//   - JPG/PNG to PDF conversion
package pdfsdk

import (
	"github.com/infosec554/golang-pdf-sdk/pkg/gotenberg"
	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
	"github.com/infosec554/golang-pdf-sdk/service"
)

// Version of the SDK
const Version = "1.0.0"

// New creates a new PDF service with the given Gotenberg URL.
// Use "http://localhost:3000" for local Gotenberg instance.
func New(gotenbergURL string) service.PDFService {
	return service.NewWithGotenberg(gotenbergURL)
}

// NewWithLogger creates a PDF service with custom logger.
func NewWithLogger(gotenbergURL string, log logger.ILogger) service.PDFService {
	gotClient := gotenberg.New(gotenbergURL)
	return service.New(log, gotClient)
}
