package main

import (
	"context"
	"fmt"
	"os"

	"github.com/infosec554/golang-pdf-sdk/service"
)

func main() {
	// Create PDF service with Gotenberg URL
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	// Example: Compress a PDF
	exampleCompress(pdfService)

	// Example: Rotate a PDF
	exampleRotate(pdfService)

	// Example: Split a PDF
	exampleSplit(pdfService)

	fmt.Println("All examples completed!")
}

func exampleCompress(pdf service.PDFService) {
	fmt.Println("=== PDF Compression Example ===")

	// Read input PDF
	inputBytes, err := os.ReadFile("input.pdf")
	if err != nil {
		fmt.Println("No input.pdf found, skipping compression example")
		return
	}

	// Compress
	outputBytes, err := pdf.Compress().CompressBytes(inputBytes)
	if err != nil {
		fmt.Println("Compression failed:", err)
		return
	}

	// Save output
	if err := os.WriteFile("compressed.pdf", outputBytes, 0644); err != nil {
		fmt.Println("Failed to save:", err)
		return
	}

	fmt.Printf("Compressed: %d -> %d bytes (%.1f%% reduction)\n",
		len(inputBytes), len(outputBytes),
		(1-float64(len(outputBytes))/float64(len(inputBytes)))*100)
}

func exampleRotate(pdf service.PDFService) {
	fmt.Println("\n=== PDF Rotation Example ===")

	inputBytes, err := os.ReadFile("input.pdf")
	if err != nil {
		fmt.Println("No input.pdf found, skipping rotation example")
		return
	}

	// Rotate 90 degrees
	outputBytes, err := pdf.Rotate().RotateBytes(inputBytes, 90, "all")
	if err != nil {
		fmt.Println("Rotation failed:", err)
		return
	}

	if err := os.WriteFile("rotated.pdf", outputBytes, 0644); err != nil {
		fmt.Println("Failed to save:", err)
		return
	}

	fmt.Println("Rotated PDF saved as rotated.pdf")
}

func exampleSplit(pdf service.PDFService) {
	fmt.Println("\n=== PDF Split Example ===")

	inputBytes, err := os.ReadFile("input.pdf")
	if err != nil {
		fmt.Println("No input.pdf found, skipping split example")
		return
	}

	// Split first 2 pages
	zipBytes, err := pdf.Split().SplitBytes(inputBytes, "1-2")
	if err != nil {
		fmt.Println("Split failed:", err)
		return
	}

	if err := os.WriteFile("split_pages.zip", zipBytes, 0644); err != nil {
		fmt.Println("Failed to save:", err)
		return
	}

	fmt.Println("Split pages saved as split_pages.zip")
}

func exampleWordToPDF(pdf service.PDFService) {
	fmt.Println("\n=== Word to PDF Example ===")

	inputBytes, err := os.ReadFile("document.docx")
	if err != nil {
		fmt.Println("No document.docx found, skipping")
		return
	}

	ctx := context.Background()
	outputBytes, err := pdf.WordToPDF().ConvertBytes(ctx, inputBytes, "document.docx")
	if err != nil {
		fmt.Println("Conversion failed:", err)
		return
	}

	if err := os.WriteFile("document.pdf", outputBytes, 0644); err != nil {
		fmt.Println("Failed to save:", err)
		return
	}

	fmt.Println("Converted to document.pdf")
}
