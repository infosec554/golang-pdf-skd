package service

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

type FormService interface {
	FillForm(input []byte, data map[string]interface{}) ([]byte, error)
	ListFormFields(input []byte) ([]string, error)
	RemoveFormFields(input []byte) ([]byte, error)
}

type formService struct {
	log logger.ILogger
}

func NewFormService(log logger.ILogger) FormService {
	return &formService{log: log}
}

func (s *formService) FillForm(input []byte, data map[string]interface{}) ([]byte, error) {
	s.log.Info("FormService.FillForm called")

	tmpDir, err := os.MkdirTemp("", "pdf-form-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	jsonPath := filepath.Join(tmpDir, "data.json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return nil, err
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	if err := api.FillFormFile(inputPath, jsonPath, outputPath, nil); err != nil {
		return nil, err
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("Form filled", logger.Int("outputSize", len(output)))
	return output, nil
}

func (s *formService) ListFormFields(input []byte) ([]string, error) {
	s.log.Info("FormService.ListFormFields called")

	s.log.Warn("ListFormFields is limited in this version")
	return []string{}, nil
}

func (s *formService) RemoveFormFields(input []byte) ([]byte, error) {
	s.log.Info("FormService.RemoveFormFields called")

	tmpDir, err := os.MkdirTemp("", "pdf-remove-form-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	if err := api.OptimizeFile(inputPath, "", nil); err != nil {
		return nil, err
	}

	output, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	return output, nil
}
