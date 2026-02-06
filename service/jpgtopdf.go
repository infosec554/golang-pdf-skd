package service

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/jung-kurt/gofpdf"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

// JPGToPDFService converts JPG/PNG images to PDF
type JPGToPDFService interface {
	// Convert converts single image to PDF bytes
	Convert(input io.Reader, filename string) ([]byte, error)
	// ConvertMultiple converts multiple images to single PDF
	ConvertMultiple(inputs []io.Reader, filenames []string) ([]byte, error)
	// ConvertFiles converts image files to PDF file
	ConvertFiles(inputPaths []string, outputPath string) error
	// ConvertBytes converts image bytes to PDF bytes
	ConvertBytes(input []byte, filename string) ([]byte, error)
	// ConvertMultipleBytes converts multiple image bytes to PDF
	ConvertMultipleBytes(inputs [][]byte, filenames []string) ([]byte, error)
}

type jpgToPDFService struct {
	log logger.ILogger
}

// NewJPGToPDFService creates a new JPG to PDF service
func NewJPGToPDFService(log logger.ILogger) JPGToPDFService {
	return &jpgToPDFService{
		log: log,
	}
}

// Convert converts single image to PDF
func (s *jpgToPDFService) Convert(input io.Reader, filename string) ([]byte, error) {
	s.log.Info("JPGToPDFService.Convert called", logger.String("filename", filename))

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return s.ConvertBytes(inputBytes, filename)
}

// ConvertMultiple converts multiple images to single PDF
func (s *jpgToPDFService) ConvertMultiple(inputs []io.Reader, filenames []string) ([]byte, error) {
	s.log.Info("JPGToPDFService.ConvertMultiple called", logger.Int("count", len(inputs)))

	var inputBytes [][]byte
	for _, r := range inputs {
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		inputBytes = append(inputBytes, data)
	}

	return s.ConvertMultipleBytes(inputBytes, filenames)
}

// ConvertFiles converts image files to PDF
func (s *jpgToPDFService) ConvertFiles(inputPaths []string, outputPath string) error {
	s.log.Info("JPGToPDFService.ConvertFiles called", logger.Int("count", len(inputPaths)))

	// Sort files for consistent ordering
	sort.Strings(inputPaths)

	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, imgPath := range inputPaths {
		ext := filepath.Ext(imgPath)
		imgType := "JPG"
		if ext == ".png" || ext == ".PNG" {
			imgType = "PNG"
		}

		// Add new page
		pdf.AddPage()

		// Get page dimensions
		pageW, pageH := pdf.GetPageSize()

		// Register and add image
		pdf.RegisterImageOptions(imgPath, gofpdf.ImageOptions{ImageType: imgType, ReadDpi: true})

		// Fit image to page
		pdf.ImageOptions(imgPath, 0, 0, pageW, pageH, false, gofpdf.ImageOptions{ImageType: imgType}, 0, "")
	}

	if err := pdf.OutputFileAndClose(outputPath); err != nil {
		s.log.Error("Failed to create PDF", logger.Error(err))
		return err
	}

	s.log.Info("Images to PDF conversion completed", logger.String("output", outputPath))
	return nil
}

// ConvertBytes converts single image bytes to PDF bytes
func (s *jpgToPDFService) ConvertBytes(input []byte, filename string) ([]byte, error) {
	return s.ConvertMultipleBytes([][]byte{input}, []string{filename})
}

// ConvertMultipleBytes converts multiple image bytes to PDF bytes
func (s *jpgToPDFService) ConvertMultipleBytes(inputs [][]byte, filenames []string) ([]byte, error) {
	s.log.Info("JPGToPDFService.ConvertMultipleBytes called", logger.Int("count", len(inputs)))

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "jpg-to-pdf-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Write images to temp files
	var inputPaths []string
	for i, data := range inputs {
		filename := fmt.Sprintf("image_%d.jpg", i)
		if i < len(filenames) {
			filename = filenames[i]
		}
		tmpPath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(tmpPath, data, 0644); err != nil {
			return nil, err
		}
		inputPaths = append(inputPaths, tmpPath)
	}

	// Convert to PDF
	outputPath := filepath.Join(tmpDir, "output.pdf")
	if err := s.ConvertFiles(inputPaths, outputPath); err != nil {
		return nil, err
	}

	// Read output
	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("Images to PDF conversion completed", logger.Int("outputSize", len(output)))
	return output, nil
}

// Helper to detect image type from bytes
func detectImageType(data []byte) string {
	if len(data) < 3 {
		return "JPG"
	}
	// PNG signature
	if bytes.HasPrefix(data, []byte{0x89, 0x50, 0x4E}) {
		return "PNG"
	}
	return "JPG"
}
