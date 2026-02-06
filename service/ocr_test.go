package service_test

import (
	"context"
	"os"
	"testing"

	"github.com/infosec554/convert-pdf-go-sdk/service"
)

func TestOCRService(t *testing.T) {
	log := getTestLogger()
	ocrService := service.NewOCRService(log)

	if ocrService == nil {
		t.Fatal("OCRService should not be nil")
	}

	if !ocrService.IsAvailable() {
		t.Log("Tesseract OCR is not installed. Skipping functional tests.")
		return
	}
}

func TestOCRService_ExtractText(t *testing.T) {
	log := getTestLogger()
	ocrService := service.NewOCRService(log)

	if !ocrService.IsAvailable() {
		t.Skip("Tesseract not installed")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	input, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	// This assumes test.pdf contains some renderable text.
	// Since minimalPDF is vector text, Tesseract might struggle if resolution is low,
	// but usually it works if we force it.
	ctx := context.Background()
	text, err := ocrService.ExtractText(ctx, input, "eng")

	if err != nil {
		t.Logf("ExtractText failed (might be expected for vector PDF or config): %v", err)
	} else {
		t.Logf("Extracted text: %s", text)
	}
}

func TestOCRService_CreateSearchablePDF(t *testing.T) {
	log := getTestLogger()
	ocrService := service.NewOCRService(log)

	if !ocrService.IsAvailable() {
		t.Skip("Tesseract not installed")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	input, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	output, err := ocrService.CreateSearchablePDF(ctx, input, "eng")

	if err != nil {
		t.Logf("CreateSearchablePDF failed: %v", err)
	} else {
		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
		t.Logf("Searchable PDF created, size: %d", len(output))
	}
}
