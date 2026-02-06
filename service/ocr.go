package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

// OCRService provides Optical Character Recognition capabilities
type OCRService interface {
	// ExtractText extracts text from scanned PDF or images using OCR
	ExtractText(ctx context.Context, input []byte, lang string) (string, error)

	// CreateSearchablePDF converts scanned PDF to searchable PDF (adds text layer)
	CreateSearchablePDF(ctx context.Context, input []byte, lang string) ([]byte, error)

	// IsAvailable checks if Tesseract OCR and pdftoppm are installed
	IsAvailable() bool
}

type ocrService struct {
	log logger.ILogger
}

// NewOCRService creates a new OCR service
func NewOCRService(log logger.ILogger) OCRService {
	return &ocrService{log: log}
}

func (s *ocrService) IsAvailable() bool {
	_, err := exec.LookPath("tesseract")
	if err != nil {
		return false
	}
	_, err = exec.LookPath("pdftoppm")
	return err == nil
}

func (s *ocrService) convertPDFToImages(ctx context.Context, pdfPath string, outDir string) ([]string, error) {
	args := []string{"-jpeg", "-r", "300", pdfPath, filepath.Join(outDir, "page")}
	cmd := exec.CommandContext(ctx, "pdftoppm", args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("pdftoppm failed: %v, output: %s", err, string(output))
	}

	files, err := os.ReadDir(outDir)
	if err != nil {
		return nil, err
	}

	var imageFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".jpg") {
			imageFiles = append(imageFiles, filepath.Join(outDir, f.Name()))
		}
	}
	sort.Strings(imageFiles)

	return imageFiles, nil
}

func (s *ocrService) ExtractText(ctx context.Context, input []byte, lang string) (string, error) {
	s.log.Info("OCRService.ExtractText called", logger.String("lang", lang))

	if !s.IsAvailable() {
		return "", fmt.Errorf("dependencies missing: install 'tesseract-ocr' and 'poppler-utils'")
	}

	if lang == "" {
		lang = "eng"
	}

	tmpDir, err := os.MkdirTemp("", "ocr-extract-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return "", err
	}

	images, err := s.convertPDFToImages(ctx, inputPath, tmpDir)
	if err != nil {
		return "", err
	}

	var fullText strings.Builder

	for _, imgPath := range images {
		outBase := filepath.Join(tmpDir, filepath.Base(imgPath)+"_out")
		args := []string{imgPath, outBase, "-l", lang}
		cmd := exec.CommandContext(ctx, "tesseract", args...)

		if output, err := cmd.CombinedOutput(); err != nil {
			s.log.Error("Tesseract failed on page", logger.String("image", filepath.Base(imgPath)), logger.String("error", string(output)))
			continue
		}

		txtPath := outBase + ".txt"
		content, err := os.ReadFile(txtPath)
		if err == nil {
			fullText.Write(content)
			fullText.WriteString("\n\n")
		}
	}

	return fullText.String(), nil
}

func (s *ocrService) CreateSearchablePDF(ctx context.Context, input []byte, lang string) ([]byte, error) {
	s.log.Info("OCRService.CreateSearchablePDF called", logger.String("lang", lang))

	if !s.IsAvailable() {
		return nil, fmt.Errorf("dependencies missing: install 'tesseract-ocr' and 'poppler-utils'")
	}

	if lang == "" {
		lang = "eng"
	}

	tmpDir, err := os.MkdirTemp("", "ocr-pdf-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	images, err := s.convertPDFToImages(ctx, inputPath, tmpDir)
	if err != nil {
		return nil, err
	}

	var pdfPages []string

	for _, imgPath := range images {
		outBase := filepath.Join(tmpDir, filepath.Base(imgPath)+"_out")
		args := []string{imgPath, outBase, "-l", lang, "pdf"}
		cmd := exec.CommandContext(ctx, "tesseract", args...)

		if output, err := cmd.CombinedOutput(); err != nil {
			s.log.Error("Tesseract failed on page", logger.String("image", filepath.Base(imgPath)), logger.String("error", string(output)))
			return nil, fmt.Errorf("ocr failed on page: %w", err)
		}

		pdfPages = append(pdfPages, outBase+".pdf")
	}

	if len(pdfPages) == 0 {
		return nil, fmt.Errorf("no pages processed")
	}

	outputPath := filepath.Join(tmpDir, "final.pdf")
	conf := model.NewDefaultConfiguration()

	if err := api.MergeCreateFile(pdfPages, outputPath, false, conf); err != nil {
		return nil, fmt.Errorf("merge failed: %w", err)
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("Searchable PDF created", logger.Int("size", len(output)))
	return output, nil
}
