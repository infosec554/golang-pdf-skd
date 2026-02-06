package service

import (
	"bytes"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

// MergeService merges multiple PDF files into one
type MergeService interface {
	// Merge merges multiple PDF readers into single PDF bytes
	Merge(inputs []io.Reader) ([]byte, error)
	// MergeFiles merges multiple PDF files into one
	MergeFiles(inputPaths []string, outputPath string) error
	// MergeBytes merges multiple PDF byte slices into one
	MergeBytes(inputs [][]byte) ([]byte, error)
}

type mergeService struct {
	log logger.ILogger
}

// NewMergeService creates a new merge service
func NewMergeService(log logger.ILogger) MergeService {
	return &mergeService{log: log}
}

// Merge merges multiple PDF readers into single PDF bytes
func (s *mergeService) Merge(inputs []io.Reader) ([]byte, error) {
	s.log.Info("MergeService.Merge called", logger.Int("inputCount", len(inputs)))

	var inputBytes [][]byte
	for _, r := range inputs {
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		inputBytes = append(inputBytes, data)
	}

	return s.MergeBytes(inputBytes)
}

// MergeFiles merges multiple PDF files into one output file
func (s *mergeService) MergeFiles(inputPaths []string, outputPath string) error {
	s.log.Info("MergeService.MergeFiles called", logger.Int("inputCount", len(inputPaths)))

	conf := model.NewDefaultConfiguration()

	if err := api.MergeCreateFile(inputPaths, outputPath, false, conf); err != nil {
		s.log.Error("pdfcpu merge failed", logger.Error(err))
		return err
	}

	s.log.Info("PDF merge completed", logger.String("output", outputPath))
	return nil
}

// MergeBytes merges multiple PDF byte slices into one
func (s *mergeService) MergeBytes(inputs [][]byte) ([]byte, error) {
	s.log.Info("MergeService.MergeBytes called", logger.Int("inputCount", len(inputs)))

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "pdf-merge-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// Write inputs to temp files
	var inputPaths []string
	for i, data := range inputs {
		tmpPath := tmpDir + "/" + string(rune('a'+i)) + ".pdf"
		if err := os.WriteFile(tmpPath, data, 0644); err != nil {
			return nil, err
		}
		inputPaths = append(inputPaths, tmpPath)
	}

	// Merge to output file
	outputPath := tmpDir + "/merged.pdf"
	if err := s.MergeFiles(inputPaths, outputPath); err != nil {
		return nil, err
	}

	// Read output
	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("PDF merge completed", logger.Int("outputSize", len(output)))
	return output, nil
}

// Utility for creating temp filename
func tempPDFName(prefix string, index int) string {
	buf := bytes.NewBufferString(prefix)
	buf.WriteByte(byte('0' + index%10))
	buf.WriteString(".pdf")
	return buf.String()
}
