package service_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
	"github.com/infosec554/convert-pdf-go-sdk/service"
)

var minimalPDF = []byte(`%PDF-1.4
1 0 obj
<< /Type /Catalog /Pages 2 0 R >>
endobj
2 0 obj
<< /Type /Pages /Kids [3 0 R] /Count 1 >>
endobj
3 0 obj
<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>
endobj
xref
0 4
0000000000 65535 f 
0000000009 00000 n 
0000000058 00000 n 
0000000115 00000 n 
trailer
<< /Size 4 /Root 1 0 R >>
startxref
196
%%EOF`)

func getTestLogger() logger.ILogger {
	return logger.New("test")
}

func TestBatchProcessor(t *testing.T) {
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	batch := service.NewBatchProcessor(pdfService, 2)
	if batch == nil {
		t.Fatal("Expected batch processor, got nil")
	}
}

func TestBatchProcessorWithZeroWorkers(t *testing.T) {
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	batch := service.NewBatchProcessor(pdfService, 0)
	if batch == nil {
		t.Fatal("Expected batch processor with default workers")
	}
}

func TestPipeline(t *testing.T) {
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	pipeline := service.NewPipeline(pdfService)
	if pipeline == nil {
		t.Fatal("Expected pipeline, got nil")
	}

	p := pipeline.Compress().Rotate(90, "all").Watermark("TEST", nil).Protect("password")
	if p == nil {
		t.Fatal("Pipeline chaining should return pipeline")
	}
}

func TestPipelineReset(t *testing.T) {
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	pipeline := service.NewPipeline(pdfService)
	pipeline.Compress().Watermark("TEST", nil)

	pipeline.Reset()
}

func TestInfoService(t *testing.T) {
	log := getTestLogger()
	infoService := service.NewInfoService(log)
	if infoService == nil {
		t.Fatal("Expected info service, got nil")
	}
}

func TestPageService(t *testing.T) {
	log := getTestLogger()
	pageService := service.NewPageService(log)
	if pageService == nil {
		t.Fatal("Expected page service, got nil")
	}
}

func TestTextService(t *testing.T) {
	log := getTestLogger()
	textService := service.NewTextService(log)
	if textService == nil {
		t.Fatal("Expected text service, got nil")
	}
}

func TestMetadataService(t *testing.T) {
	log := getTestLogger()
	metadataService := service.NewMetadataService(log)
	if metadataService == nil {
		t.Fatal("Expected metadata service, got nil")
	}
}

func TestImageExtractService(t *testing.T) {
	log := getTestLogger()
	imageService := service.NewImageExtractService(log)
	if imageService == nil {
		t.Fatal("Expected image extract service, got nil")
	}
}

func TestPDFServiceInterface(t *testing.T) {
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	if pdfService.Compress() == nil {
		t.Error("Compress() should not return nil")
	}
	if pdfService.Merge() == nil {
		t.Error("Merge() should not return nil")
	}
	if pdfService.Split() == nil {
		t.Error("Split() should not return nil")
	}
	if pdfService.Rotate() == nil {
		t.Error("Rotate() should not return nil")
	}
	if pdfService.Watermark() == nil {
		t.Error("Watermark() should not return nil")
	}
	if pdfService.Protect() == nil {
		t.Error("Protect() should not return nil")
	}
	if pdfService.Unlock() == nil {
		t.Error("Unlock() should not return nil")
	}
	if pdfService.Info() == nil {
		t.Error("Info() should not return nil")
	}
	if pdfService.Pages() == nil {
		t.Error("Pages() should not return nil")
	}
	if pdfService.Text() == nil {
		t.Error("Text() should not return nil")
	}
	if pdfService.Metadata() == nil {
		t.Error("Metadata() should not return nil")
	}
	if pdfService.Images() == nil {
		t.Error("Images() should not return nil")
	}
	if pdfService.JPGToPDF() == nil {
		t.Error("JPGToPDF() should not return nil")
	}
	if pdfService.PDFToJPG() == nil {
		t.Error("PDFToJPG() should not return nil")
	}
	if pdfService.WordToPDF() == nil {
		t.Error("WordToPDF() should not return nil")
	}
	if pdfService.ExcelToPDF() == nil {
		t.Error("ExcelToPDF() should not return nil")
	}
	if pdfService.PowerPointToPDF() == nil {
		t.Error("PowerPointToPDF() should not return nil")
	}
}

func TestBatch(t *testing.T) {
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	batch := pdfService.Batch(5)
	if batch == nil {
		t.Error("Batch() should not return nil")
	}
}

func TestPipelineFromService(t *testing.T) {
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	pipeline := pdfService.Pipeline()
	if pipeline == nil {
		t.Error("Pipeline() should not return nil")
	}
}

func TestCompressWithRealPDF(t *testing.T) {
	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found, skipping integration test")
	}

	input, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Skip("Cannot read test.pdf")
	}

	pdfService := service.NewWithGotenberg("http://localhost:3000")
	output, err := pdfService.Compress().CompressBytes(input)
	if err != nil {
		t.Errorf("Compress failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("Output should not be empty")
	}
}

func TestBatchWithContextCancel(t *testing.T) {
	pdfService := service.NewWithGotenberg("http://localhost:3000")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	inputs := [][]byte{minimalPDF, minimalPDF}
	results := pdfService.Batch(2).CompressBatch(ctx, inputs)

	for _, r := range results {
		if r.Error == nil {
			continue
		}
		if r.Error != context.Canceled {
			t.Logf("Error (expected): %v", r.Error)
		}
	}
}

func createTempTestDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "pdfsdk-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

func TestMergeFiles(t *testing.T) {
	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	tmpDir := createTempTestDir(t)
	defer os.RemoveAll(tmpDir)

	outputPath := filepath.Join(tmpDir, "merged.pdf")

	pdfService := service.NewWithGotenberg("http://localhost:3000")
	err := pdfService.Merge().MergeFiles([]string{testPDFPath, testPDFPath}, outputPath)
	if err != nil {
		t.Errorf("MergeFiles failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}
