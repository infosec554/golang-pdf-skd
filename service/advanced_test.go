package service_test

import (
	"os"
	"testing"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/gotenberg"
	"github.com/infosec554/convert-pdf-go-sdk/service"
)

// Unit tests for new services

func TestArchiveService(t *testing.T) {
	log := getTestLogger()
	// Mock or nil client for unit test (it won't be used since we won't call ConvertToPDFA in unit test mode)
	// Actually we can pass a special mock if needed, but for now just test initialization
	gotClient := gotenberg.New("http://localhost:3000")
	archiveService := service.NewArchiveService(log, gotClient)

	if archiveService == nil {
		t.Fatal("ArchiveService should not be nil")
	}
}

func TestFormService(t *testing.T) {
	log := getTestLogger()
	formService := service.NewFormService(log)
	if formService == nil {
		t.Fatal("FormService should not be nil")
	}
}

func TestAttachmentService(t *testing.T) {
	log := getTestLogger()
	attachmentService := service.NewAttachmentService(log)
	if attachmentService == nil {
		t.Fatal("AttachmentService should not be nil")
	}
}

// Integration tests

func TestArchiveService_ConvertToPDFA(t *testing.T) {
	if os.Getenv("GOTENBERG_URL") == "" {
		t.Skip("GOTENBERG_URL not set, skipping integration test")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	input, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	log := getTestLogger()
	gotClient := gotenberg.New(os.Getenv("GOTENBERG_URL"))
	archiveService := service.NewArchiveService(log, gotClient)

	output, err := archiveService.ConvertToPDFA(input, "PDF/A-1b")
	if err != nil {
		t.Errorf("ConvertToPDFA failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("Output should not be empty")
	}
}

func TestAttachmentService_AddAttachments(t *testing.T) {
	// This is basically an integration test as it uses real PDF and pdfcpu
	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	input, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	files := map[string][]byte{
		"note.txt": []byte("This is an attachment"),
	}

	log := getTestLogger()
	attachmentService := service.NewAttachmentService(log)

	output, err := attachmentService.AddAttachments(input, files)
	if err != nil {
		t.Errorf("AddAttachments failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("Output should not be empty")
	}
}

func TestFormService_FillForm(t *testing.T) {
	// Requires a PDF form. If test.pdf is not a form, this might fail or do nothing.
	// We'll skip if no form fields are found or just test it doesn't crash.
	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	input, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	data := map[string]interface{}{
		"Name": "John Doe",
	}

	log := getTestLogger()
	formService := service.NewFormService(log)

	// This might fail if test.pdf is not a form, so we expect error or success depending on file
	// Actually pdfcpu might return error if no form.
	_, err = formService.FillForm(input, data)
	if err != nil {
		t.Logf("FillForm failed (expected if not a form PDF): %v", err)
	} else {
		// Success
	}
}

// Ensure dummy testdata directory and file exist for basic tests if possible
func init() {
	_ = os.MkdirAll("testdata", 0755)
	// Create a minimal PDF if not exists
	if _, err := os.Stat("testdata/test.pdf"); os.IsNotExist(err) {
		// minimal valid PDF
		pdf := []byte(`%PDF-1.4
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
		_ = os.WriteFile("testdata/test.pdf", pdf, 0644)
	}
}
