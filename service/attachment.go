package service

import (
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

type AttachmentService interface {
	AddAttachments(input []byte, files map[string][]byte) ([]byte, error)
	ListAttachments(input []byte) ([]string, error)
	ExtractAttachments(input []byte) (map[string][]byte, error)
	RemoveAttachments(input []byte) ([]byte, error)
}

type attachmentService struct {
	log logger.ILogger
}

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

	var fileNames []string
	for name, content := range files {
		// Security fix: sanitize filename
		safeName := filepath.Base(name)
		filePath := filepath.Join(tmpDir, safeName)
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

	s.log.Warn("RemoveAttachments not fully implemented")
	return input, nil
}
