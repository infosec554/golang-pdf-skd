package gotenberg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Client interface {
	PDFToWord(ctx context.Context, pdfPath string) ([]byte, error)
	WordToPDF(ctx context.Context, wordPath string) ([]byte, error)
	ExcelToPDF(ctx context.Context, excelPath string) ([]byte, error)
	PowerPointToPDF(ctx context.Context, pptPath string) ([]byte, error)
	HTMLToPDF(ctx context.Context, htmlPath string) ([]byte, error)
	WordToPDFStream(ctx context.Context, wordPath string, w io.Writer) error
	ExcelToPDFStream(ctx context.Context, excelPath string, w io.Writer) error
}

type gotenbergClient struct {
	baseURL    string
	httpClient *http.Client
	bufferPool *sync.Pool
}

func New(url string) Client {
	return &gotenbergClient{
		baseURL:    url,
		httpClient: http.DefaultClient,
		bufferPool: newBufferPool(),
	}
}

func NewWithClient(url string, client *http.Client) Client {
	if client == nil {
		client = http.DefaultClient
	}
	return &gotenbergClient{
		baseURL:    url,
		httpClient: client,
		bufferPool: newBufferPool(),
	}
}

func newBufferPool() *sync.Pool {
	return &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 32*1024))
		},
	}
}

func (g *gotenbergClient) getBuffer() *bytes.Buffer {
	buf := g.bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

func (g *gotenbergClient) putBuffer(buf *bytes.Buffer) {
	if buf.Cap() > 10*1024*1024 {
		return
	}
	g.bufferPool.Put(buf)
}

func (g *gotenbergClient) PDFToWord(ctx context.Context, pdfPath string) ([]byte, error) {
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

	_ = writer.WriteField("output", "docx")
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/libreoffice/convert", requestBody)
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

func (g *gotenbergClient) WordToPDF(ctx context.Context, wordPath string) ([]byte, error) {
	file, err := os.Open(wordPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	requestBody := g.getBuffer()
	defer g.putBuffer(requestBody)

	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("files", filepath.Base(wordPath))
	if err != nil {
		return nil, fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("cannot copy file: %w", err)
	}

	_ = writer.WriteField("waitTimeout", "30s")
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/libreoffice/convert", requestBody)
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

func (g *gotenbergClient) WordToPDFStream(ctx context.Context, wordPath string, w io.Writer) error {
	file, err := os.Open(wordPath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	requestBody := g.getBuffer()
	defer g.putBuffer(requestBody)

	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("files", filepath.Base(wordPath))
	if err != nil {
		return fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("cannot copy file: %w", err)
	}

	_ = writer.WriteField("waitTimeout", "30s")
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/libreoffice/convert", requestBody)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("conversion failed: %s", string(bodyBytes))
	}

	_, err = io.Copy(w, resp.Body)
	return err
}

func (g *gotenbergClient) ExcelToPDF(ctx context.Context, excelPath string) ([]byte, error) {
	file, err := os.Open(excelPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	requestBody := g.getBuffer()
	defer g.putBuffer(requestBody)

	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("files", filepath.Base(excelPath))
	if err != nil {
		return nil, fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("cannot copy file: %w", err)
	}

	_ = writer.WriteField("waitTimeout", "30s")
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/libreoffice/convert", requestBody)
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

func (g *gotenbergClient) ExcelToPDFStream(ctx context.Context, excelPath string, w io.Writer) error {
	file, err := os.Open(excelPath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	requestBody := g.getBuffer()
	defer g.putBuffer(requestBody)

	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("files", filepath.Base(excelPath))
	if err != nil {
		return fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("cannot copy file: %w", err)
	}

	_ = writer.WriteField("waitTimeout", "30s")
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/libreoffice/convert", requestBody)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("conversion failed: %s", string(bodyBytes))
	}

	_, err = io.Copy(w, resp.Body)
	return err
}

func (g *gotenbergClient) PowerPointToPDF(ctx context.Context, pptPath string) ([]byte, error) {
	file, err := os.Open(pptPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	requestBody := g.getBuffer()
	defer g.putBuffer(requestBody)

	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("files", filepath.Base(pptPath))
	if err != nil {
		return nil, fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("cannot copy file: %w", err)
	}

	_ = writer.WriteField("waitTimeout", "30s")
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/libreoffice/convert", requestBody)
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

func (g *gotenbergClient) HTMLToPDF(ctx context.Context, htmlPath string) ([]byte, error) {
	file, err := os.Open(htmlPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open HTML file: %w", err)
	}
	defer file.Close()

	requestBody := g.getBuffer()
	defer g.putBuffer(requestBody)

	writer := multipart.NewWriter(requestBody)

	part, err := writer.CreateFormFile("files", "index.html")
	if err != nil {
		return nil, fmt.Errorf("cannot create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("cannot copy HTML file: %w", err)
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/forms/chromium/convert/html", requestBody)
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("conversion request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gotenberg error: %s", string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}
