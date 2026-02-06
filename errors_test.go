package pdfsdk_test

import (
	"testing"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrInvalidPDF", pdfsdk.ErrInvalidPDF},
		{"ErrEncryptedPDF", pdfsdk.ErrEncryptedPDF},
		{"ErrWrongPassword", pdfsdk.ErrWrongPassword},
		{"ErrEmptyInput", pdfsdk.ErrEmptyInput},
		{"ErrPageOutOfRange", pdfsdk.ErrPageOutOfRange},
		{"ErrGotenbergUnavailable", pdfsdk.ErrGotenbergUnavailable},
		{"ErrTimeout", pdfsdk.ErrTimeout},
		{"ErrWorkerPoolFull", pdfsdk.ErrWorkerPoolFull},
		{"ErrOperationCanceled", pdfsdk.ErrOperationCanceled},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s should not be nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Errorf("%s should have error message", tt.name)
			}
		})
	}
}

func TestPDFError(t *testing.T) {
	err := pdfsdk.NewError("compress", pdfsdk.ErrInvalidPDF)

	if err.Op != "compress" {
		t.Errorf("Expected op=compress, got %s", err.Op)
	}
	if err.Err != pdfsdk.ErrInvalidPDF {
		t.Error("Expected underlying error to be ErrInvalidPDF")
	}

	// Test error message
	msg := err.Error()
	if msg == "" {
		t.Error("Error message should not be empty")
	}
}

func TestWrapError(t *testing.T) {
	err := pdfsdk.WrapError("merge", "file1.pdf", pdfsdk.ErrEmptyInput)

	if err.Op != "merge" {
		t.Errorf("Expected op=merge, got %s", err.Op)
	}
	if err.Input != "file1.pdf" {
		t.Errorf("Expected input=file1.pdf, got %s", err.Input)
	}

	// Test Unwrap
	if err.Unwrap() != pdfsdk.ErrEmptyInput {
		t.Error("Unwrap should return underlying error")
	}
}

func TestIsErrorHelpers(t *testing.T) {
	if !pdfsdk.IsInvalidPDF(pdfsdk.ErrInvalidPDF) {
		t.Error("IsInvalidPDF should return true for ErrInvalidPDF")
	}
	if pdfsdk.IsInvalidPDF(pdfsdk.ErrTimeout) {
		t.Error("IsInvalidPDF should return false for ErrTimeout")
	}

	if !pdfsdk.IsEncrypted(pdfsdk.ErrEncryptedPDF) {
		t.Error("IsEncrypted should return true for ErrEncryptedPDF")
	}

	if !pdfsdk.IsTimeout(pdfsdk.ErrTimeout) {
		t.Error("IsTimeout should return true for ErrTimeout")
	}

	if !pdfsdk.IsGotenbergUnavailable(pdfsdk.ErrGotenbergUnavailable) {
		t.Error("IsGotenbergUnavailable should return true for ErrGotenbergUnavailable")
	}
}
