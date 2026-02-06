package service

import (
	"bytes"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

// CompressService provides PDF compression functionality
type CompressService interface {
	// Compress compresses a PDF file and returns the compressed bytes
	Compress(input io.Reader) ([]byte, error)
	// CompressFile compresses a PDF from file path
	CompressFile(inputPath, outputPath string) error
	// CompressBytes compresses PDF bytes
	CompressBytes(input []byte) ([]byte, error)
}

type compressService struct {
	log logger.ILogger
}

// NewCompressService creates a new compress service
func NewCompressService(log logger.ILogger) CompressService {
	return &compressService{
		log: log,
	}
}

// Compress compresses PDF from reader
func (s *compressService) Compress(input io.Reader) ([]byte, error) {
	s.log.Info("CompressService.Compress called")

	// Read input to bytes
	inputBytes, err := io.ReadAll(input)
	if err != nil {
		s.log.Error("Failed to read input", logger.Error(err))
		return nil, err
	}

	return s.CompressBytes(inputBytes)
}

// CompressFile compresses PDF file
func (s *compressService) CompressFile(inputPath, outputPath string) error {
	s.log.Info("CompressService.CompressFile called", logger.String("input", inputPath))

	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.OPTIMIZE

	if err := api.OptimizeFile(inputPath, outputPath, conf); err != nil {
		s.log.Error("pdfcpu optimize failed", logger.Error(err))
		return err
	}

	s.log.Info("PDF compression completed", logger.String("output", outputPath))
	return nil
}

// CompressBytes compresses PDF bytes
func (s *compressService) CompressBytes(input []byte) ([]byte, error) {
	s.log.Info("CompressService.CompressBytes called")

	// Create temp files for processing
	tmpInput, err := os.CreateTemp("", "pdf-compress-input-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())
	defer tmpInput.Close()

	tmpOutput, err := os.CreateTemp("", "pdf-compress-output-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpOutput.Name())
	tmpOutput.Close()

	// Write input to temp file
	if _, err := tmpInput.Write(input); err != nil {
		return nil, err
	}
	tmpInput.Close()

	// Compress
	if err := s.CompressFile(tmpInput.Name(), tmpOutput.Name()); err != nil {
		return nil, err
	}

	// Read output
	output, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, err
	}

	// Log compression ratio
	ratio := float64(len(output)) / float64(len(input)) * 100
	s.log.Info("Compression completed",
		logger.Int("inputSize", len(input)),
		logger.Int("outputSize", len(output)),
		logger.String("ratio", bytes.NewBufferString("").String()),
	)
	_ = ratio

	return output, nil
}
