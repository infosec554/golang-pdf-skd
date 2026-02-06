package service

import (
	"fmt"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

type RotateService interface {
	Rotate(input io.Reader, angle int, pages string) ([]byte, error)
	RotateFile(inputPath, outputPath string, angle int, pages string) error
	RotateBytes(input []byte, angle int, pages string) ([]byte, error)
}

type rotateService struct {
	log logger.ILogger
}

func NewRotateService(log logger.ILogger) RotateService {
	return &rotateService{
		log: log,
	}
}

func (s *rotateService) Rotate(input io.Reader, angle int, pages string) ([]byte, error) {
	s.log.Info("RotateService.Rotate called", logger.Int("angle", angle))

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return s.RotateBytes(inputBytes, angle, pages)
}

func (s *rotateService) RotateFile(inputPath, outputPath string, angle int, pages string) error {
	s.log.Info("RotateService.RotateFile called", logger.String("input", inputPath), logger.Int("angle", angle))

	if angle != 90 && angle != 180 && angle != 270 {
		return fmt.Errorf("invalid angle: %d (must be 90, 180 or 270)", angle)
	}

	inputBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, inputBytes, 0644); err != nil {
		return err
	}

	var selectedPages []string
	if pages != "" && pages != "all" {
		selectedPages = []string{pages}
	}

	if err := api.RotateFile(outputPath, "", angle, selectedPages, nil); err != nil {
		s.log.Error("pdfcpu rotate failed", logger.Error(err))
		os.Remove(outputPath)
		return fmt.Errorf("rotate failed: %w", err)
	}

	s.log.Info("PDF rotation completed", logger.String("output", outputPath))
	return nil
}

func (s *rotateService) RotateBytes(input []byte, angle int, pages string) ([]byte, error) {
	s.log.Info("RotateService.RotateBytes called", logger.Int("angle", angle))

	tmpInput, err := os.CreateTemp("", "pdf-rotate-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())

	tmpOutput, err := os.CreateTemp("", "pdf-rotate-out-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpOutput.Name())
	tmpOutput.Close()

	if _, err := tmpInput.Write(input); err != nil {
		tmpInput.Close()
		return nil, err
	}
	tmpInput.Close()

	if err := s.RotateFile(tmpInput.Name(), tmpOutput.Name(), angle, pages); err != nil {
		return nil, err
	}

	output, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, err
	}

	s.log.Info("PDF rotation completed", logger.Int("outputSize", len(output)))
	return output, nil
}
