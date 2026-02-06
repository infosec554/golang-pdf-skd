package service

import (
	"context"
	"io"
	"os"

	"github.com/infosec554/golang-pdf-sdk/pkg/gotenberg"
	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

// WordToPDFService converts Word documents to PDF
type WordToPDFService interface {
	// Convert converts Word document to PDF bytes
	Convert(ctx context.Context, input io.Reader, filename string) ([]byte, error)
	// ConvertFile converts Word file to PDF
	ConvertFile(ctx context.Context, inputPath, outputPath string) error
	// ConvertBytes converts Word bytes to PDF bytes
	ConvertBytes(ctx context.Context, input []byte, filename string) ([]byte, error)
}

type wordToPDFService struct {
	log       logger.ILogger
	gotClient gotenberg.Client
}

// NewWordToPDFService creates a new Word to PDF service
func NewWordToPDFService(log logger.ILogger, gotClient gotenberg.Client) WordToPDFService {
	return &wordToPDFService{
		log:       log,
		gotClient: gotClient,
	}
}

// Convert converts Word document from reader to PDF bytes
func (s *wordToPDFService) Convert(ctx context.Context, input io.Reader, filename string) ([]byte, error) {
	s.log.Info("WordToPDFService.Convert called", logger.String("filename", filename))

	// Read input to bytes
	inputBytes, err := io.ReadAll(input)
	if err != nil {
		s.log.Error("Failed to read input", logger.Error(err))
		return nil, err
	}

	return s.ConvertBytes(ctx, inputBytes, filename)
}

// ConvertFile converts Word file to PDF file
func (s *wordToPDFService) ConvertFile(ctx context.Context, inputPath, outputPath string) error {
	s.log.Info("WordToPDFService.ConvertFile called", logger.String("input", inputPath))

	resultBytes, err := s.gotClient.WordToPDF(ctx, inputPath)
	if err != nil {
		s.log.Error("Gotenberg conversion failed", logger.Error(err))
		return err
	}

	if err := os.WriteFile(outputPath, resultBytes, 0644); err != nil {
		s.log.Error("Failed to write output file", logger.Error(err))
		return err
	}

	s.log.Info("Word to PDF conversion completed", logger.String("output", outputPath))
	return nil
}

// ConvertBytes converts Word bytes to PDF bytes
func (s *wordToPDFService) ConvertBytes(ctx context.Context, input []byte, filename string) ([]byte, error) {
	s.log.Info("WordToPDFService.ConvertBytes called")

	// Create temp file for input
	tmpInput, err := os.CreateTemp("", "word-input-*"+getExtension(filename))
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())

	if _, err := tmpInput.Write(input); err != nil {
		tmpInput.Close()
		return nil, err
	}
	tmpInput.Close()

	// Convert using Gotenberg
	resultBytes, err := s.gotClient.WordToPDF(ctx, tmpInput.Name())
	if err != nil {
		s.log.Error("Gotenberg conversion failed", logger.Error(err))
		return nil, err
	}

	s.log.Info("Word to PDF conversion completed", logger.Int("outputSize", len(resultBytes)))
	return resultBytes, nil
}

// getExtension extracts file extension
func getExtension(filename string) string {
	for i := len(filename) - 1; i >= 0; i-- {
		if filename[i] == '.' {
			return filename[i:]
		}
	}
	return ".docx"
}
