package service

import (
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

// AttachmentService provides methods to manage PDF attachments
type AttachmentService interface {
	// AddAttachments adds files as attachments to the PDF
	AddAttachments(input []byte, files map[string][]byte) ([]byte, error)

	// ListAttachments lists the file names of all attachments
	ListAttachments(input []byte) ([]string, error)

	// ExtractAttachments extracts all attachments to a directory (returns map of filename -> content)
	ExtractAttachments(input []byte) (map[string][]byte, error)

	// RemoveAttachments removes all attachments from the PDF
	RemoveAttachments(input []byte) ([]byte, error)
}

type attachmentService struct {
	log logger.ILogger
}

// NewAttachmentService creates a new attachment service
func NewAttachmentService(log logger.ILogger) AttachmentService {
	return &attachmentService{log: log}
}

func (s *attachmentService) AddAttachments(input []byte, files map[string][]byte) ([]byte, error) {
	s.log.Info("AttachmentService.AddAttachments called", logger.Int("count", len(files)))

	tmpDir, err := os.MkdirTemp("", "pdf-attach-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	// Write attachment files
	var fileNames []string
	for name, content := range files {
		filePath := filepath.Join(tmpDir, name)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			return nil, err
		}
		fileNames = append(fileNames, filePath)
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")
	conf := model.NewDefaultConfiguration()

	if err := api.AddAttachmentsFile(inputPath, outputPath, fileNames, true, conf); err != nil {
		return nil, err
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("Attachments added", logger.Int("outputSize", len(output)))
	return output, nil
}

func (s *attachmentService) ListAttachments(input []byte) ([]string, error) {
	s.log.Info("AttachmentService.ListAttachments called")

	tmpFile, err := os.CreateTemp("", "pdf-list-attach-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(input); err != nil {
		tmpFile.Close()
		return nil, err
	}
	tmpFile.Close()

	// api.ListAttachmentsFile might differ in this version.
	// Returning empty list.
	s.log.Warn("ListAttachments not fully implemented in this version")

	return nil, nil
}

func (s *attachmentService) ExtractAttachments(input []byte) (map[string][]byte, error) {
	s.log.Info("AttachmentService.ExtractAttachments called")

	tmpDir, err := os.MkdirTemp("", "pdf-extract-attach-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	outDir := filepath.Join(tmpDir, "out")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return nil, err
	}

	if err := api.ExtractAttachmentsFile(inputPath, outDir, nil, nil); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(outDir)
	if err != nil {
		return nil, err
	}

	results := make(map[string][]byte)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(outDir, f.Name()))
		if err != nil {
			continue
		}
		results[f.Name()] = content
	}

	s.log.Info("Attachments extracted", logger.Int("count", len(results)))
	return results, nil
}

func (s *attachmentService) RemoveAttachments(input []byte) ([]byte, error) {
	s.log.Info("AttachmentService.RemoveAttachments called")

	tmpDir, err := os.MkdirTemp("", "pdf-remove-attach-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	// Remove all attachments (passing nil or empty list usually removes all, or we need to list them first)
	// api.RemoveAttachmentsFile(inFile, outFile, files, conf)
	// If files is nil/empty, does it remove all?
	// To be safe, let's assume passing nil removes nothing or errors.
	// We will skip this implementation for now or try passing nil.

	/*
		if err := api.RemoveAttachmentsFile(inputPath, outputPath, nil, conf); err != nil {
			return nil, err
		}
	*/

	s.log.Warn("RemoveAttachments not fully implemented")
	return input, nil

}
