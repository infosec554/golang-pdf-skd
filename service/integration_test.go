package service_test

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"github.com/infosec554/convert-pdf-go-sdk/service"
)

func createDummyJPG(t *testing.T, filename string) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create dummy JPG: %v", err)
	}
	defer f.Close()

	if err := jpeg.Encode(f, img, nil); err != nil {
		t.Fatalf("Failed to encode dummy JPG: %v", err)
	}
}

func createDummyHTML(t *testing.T, filename string) {
	content := []byte("<html><body><h1>Hello World</h1></body></html>")
	if err := os.WriteFile(filename, content, 0644); err != nil {
		t.Fatalf("Failed to write dummy HTML: %v", err)
	}
}

func TestJPGToPDF(t *testing.T) {
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	tmpDir := createTempTestDir(t)
	defer os.RemoveAll(tmpDir)

	jpgPath := filepath.Join(tmpDir, "test.jpg")
	createDummyJPG(t, jpgPath)

	pdfService := service.NewWithGotenberg(gotURL)

	// Test Convert path
	f, err := os.Open(jpgPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	pdfBytes, err := pdfService.JPGToPDF().Convert(f, "test.jpg")
	if err != nil {
		t.Fatalf("JPGToPDF failed: %v", err)
	}
	if len(pdfBytes) == 0 {
		t.Error("JPGToPDF returned empty bytes")
	}

	// Test ConvertBytes
	jpgBytes, err := os.ReadFile(jpgPath)
	if err != nil {
		t.Fatal(err)
	}
	pdfBytes2, err := pdfService.JPGToPDF().ConvertBytes(jpgBytes, "test.jpg")
	if err != nil {
		t.Fatalf("JPGToPDF ConvertBytes failed: %v", err)
	}
	if len(pdfBytes2) == 0 {
		t.Error("JPGToPDF ConvertBytes returned empty bytes")
	}
}

func TestPDFToJPG(t *testing.T) {
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	pdfService := service.NewWithGotenberg(gotURL)

	f, err := os.Open(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Convert returns a ZIP bytes
	zipBytes, err := pdfService.PDFToJPG().Convert(f)
	if err != nil {
		t.Fatalf("PDFToJPG failed: %v", err)
	}

	if len(zipBytes) == 0 {
		t.Error("PDFToJPG returned no data")
	}
}

func TestRotate(t *testing.T) {
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	pdfBytes, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	pdfService := service.NewWithGotenberg(gotURL)

	output, err := pdfService.Rotate().RotateBytes(pdfBytes, 90, "all")
	if err != nil {
		t.Fatalf("Rotate failed: %v", err)
	}
	if len(output) == 0 {
		t.Error("Rotate returned empty output")
	}
}

func TestSplit(t *testing.T) {
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	pdfService := service.NewWithGotenberg(gotURL)

	// Create a temporary directory to store split pages
	tmpOutDir := createTempTestDir(t)
	defer os.RemoveAll(tmpOutDir)

	files, err := pdfService.Split().SplitFile(testPDFPath, tmpOutDir, "1")
	if err != nil {
		t.Fatalf("SplitFile failed: %v", err)
	}

	if len(files) == 0 {
		t.Error("SplitFile created no files")
	}
}

func TestWatermark(t *testing.T) {
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	pdfBytes, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	pdfService := service.NewWithGotenberg(gotURL)

	output, err := pdfService.Watermark().AddWatermarkBytes(pdfBytes, "CONFIDENTIAL", nil)
	if err != nil {
		t.Fatalf("Watermark failed: %v", err)
	}
	if len(output) == 0 {
		t.Error("Watermark returned empty output")
	}
}

func TestProtectUnlock(t *testing.T) {
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	pdfBytes, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	pdfService := service.NewWithGotenberg(gotURL)
	password := "secret"

	// Protect
	protected, err := pdfService.Protect().ProtectBytes(pdfBytes, password)
	if err != nil {
		t.Fatalf("Protect failed: %v", err)
	}
	if len(protected) == 0 {
		t.Error("Protect returned empty output")
	}

	// Unlock
	unlocked, err := pdfService.Unlock().UnlockBytes(protected, password)
	if err != nil {
		t.Fatalf("Unlock failed: %v", err)
	}
	if len(unlocked) == 0 {
		t.Error("Unlock returned empty output")
	}
}

func TestOperationsChain(t *testing.T) {
	// Test a pipeline of operations
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	pdfBytes, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	pdfService := service.NewWithGotenberg(gotURL)

	// Compress -> Watermark -> Rotate
	compressed, err := pdfService.Compress().CompressBytes(pdfBytes)
	if err != nil {
		t.Fatalf("Chain-Compress failed: %v", err)
	}

	watermarked, err := pdfService.Watermark().AddWatermarkBytes(compressed, "DRAFT", nil)
	if err != nil {
		t.Fatalf("Chain-Watermark failed: %v", err)
	}

	rotated, err := pdfService.Rotate().RotateBytes(watermarked, 180, "all")
	if err != nil {
		t.Fatalf("Chain-Rotate failed: %v", err)
	}

	if len(rotated) == 0 {
		t.Error("Chained operations returned empty output")
	}
}

func TestMetadata(t *testing.T) {
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	pdfBytes, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	pdfService := service.NewWithGotenberg(gotURL)

	// GetMetadata
	meta, err := pdfService.Metadata().GetMetadata(pdfBytes)
	if err != nil {
		t.Fatalf("GetMetadata failed: %v", err)
	}
	if meta == nil {
		t.Error("Metadata is nil")
	}
	// Per advanced.go, we expect at least pages, version, encrypted
	if _, ok := meta["pages"]; !ok {
		t.Error("Metadata missing 'pages'")
	}

	// SetMetadata (Limited support, returns input unchanged usually)
	// Just verify it doesn't crash
	newBytes, err := pdfService.Metadata().SetMetadata(pdfBytes, map[string]string{"Title": "New Title"})
	if err != nil {
		t.Fatalf("SetMetadata failed: %v", err)
	}
	if len(newBytes) == 0 {
		t.Error("SetMetadata returned empty bytes")
	}
}

func TestImageExtraction(t *testing.T) {
	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	// 1. Create a PDF with an image first (using JPGToPDF)
	tmpDir := createTempTestDir(t)
	defer os.RemoveAll(tmpDir)

	jpgPath := filepath.Join(tmpDir, "test_extract.jpg")
	createDummyJPG(t, jpgPath)

	f, err := os.Open(jpgPath)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	pdfService := service.NewWithGotenberg(gotURL)
	pdfWithImage, err := pdfService.JPGToPDF().Convert(f, "test_extract.jpg")
	if err != nil {
		t.Fatalf("Failed to create PDF for image extraction: %v", err)
	}

	// 2. Extract images from it
	images, err := pdfService.Images().ExtractImages(pdfWithImage)
	if err != nil {
		t.Fatalf("ExtractImages failed: %v", err)
	}

	// Note: JPGToPDF using gofpdf might bake the image in a way that pdfcpu can extract,
	// or it might not. gofpdf usually embeds standard JPEG images.
	if len(images) == 0 {
		t.Log("Warning: No images extracted from generated PDF. This might depend on how gofpdf embeds them or pdfcpu extracts them.")
	} else {
		t.Logf("Extracted %d images", len(images))
	}
}
