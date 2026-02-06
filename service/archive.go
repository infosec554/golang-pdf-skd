package service

import (
	"context"
	"os"
	"path/filepath"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/gotenberg"
	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

// ArchiveService provides PDF/A conversion operations
type ArchiveService interface {
	// ConvertToPDFA converts PDF to PDF/A format (v1b, v2b, v3b)
	ConvertToPDFA(input []byte, format string) ([]byte, error)
}

type archiveService struct {
	log       logger.ILogger
	gotClient gotenberg.Client
}

// NewArchiveService creates a new archive service
func NewArchiveService(log logger.ILogger, gotClient gotenberg.Client) ArchiveService {
	return &archiveService{
		log:       log,
		gotClient: gotClient,
	}
}

func (s *archiveService) ConvertToPDFA(input []byte, format string) ([]byte, error) {
	s.log.Info("ArchiveService.ConvertToPDFA called", logger.String("format", format))

	tmpDir, err := os.MkdirTemp("", "pdf-archive-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	output, err := s.gotClient.ConvertToPDFA(context.Background(), inputPath, format)
	if err != nil {
		s.log.Error("PDF/A conversion failed", logger.Error(err))
		return nil, err
	}

	s.log.Info("Converted to PDF/A", logger.Int("outputSize", len(output)))
	return output, nil
}
