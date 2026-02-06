package gotenberg

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (g *gotenbergClient) ConvertToPDFA(ctx context.Context, pdfPath string, format string) ([]byte, error) {
	file, err := os.Open(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	requestBody := g.getBuffer()
	defer g.putBuffer(requestBody)

	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("files", filepath.Base(pdfPath))
	if err != nil {
		return nil, fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("cannot copy file: %w", err)
	}

	// Supported formats: PDF/A-1b, PDF/A-2b, PDF/A-3b
	if format == "" {
		format = "PDF/A-1b"
	}
	_ = writer.WriteField("pdfa", format)
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/pdfengines/convert", requestBody)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("conversion failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("conversion failed: %s", string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}
