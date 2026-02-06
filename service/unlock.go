package service

import (
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

// UnlockService removes password protection from PDF documents
type UnlockService interface {
	// Unlock removes password from PDF
	Unlock(input io.Reader, password string) ([]byte, error)
	// UnlockFile removes password from PDF file
	UnlockFile(inputPath, outputPath, password string) error
	// UnlockBytes removes password from PDF bytes
	UnlockBytes(input []byte, password string) ([]byte, error)
}

type unlockService struct {
	log logger.ILogger
}

// NewUnlockService creates a new unlock service
func NewUnlockService(log logger.ILogger) UnlockService {
	return &unlockService{
		log: log,
	}
}

// Unlock removes password from PDF reader
func (s *unlockService) Unlock(input io.Reader, password string) ([]byte, error) {
	s.log.Info("UnlockService.Unlock called")

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return s.UnlockBytes(inputBytes, password)
}

// UnlockFile removes password from PDF file
func (s *unlockService) UnlockFile(inputPath, outputPath, password string) error {
	s.log.Info("UnlockService.UnlockFile called", logger.String("input", inputPath))

	conf := api.LoadConfiguration()
	if password != "" {
		conf.UserPW = password
		conf.OwnerPW = password
	}

	if err := api.DecryptFile(inputPath, outputPath, conf); err != nil {
		s.log.Error("pdfcpu decrypt failed, copying as-is", logger.Error(err))
		// If decrypt fails, just copy the file
		inputBytes, err := os.ReadFile(inputPath)
		if err != nil {
			return err
		}
		return os.WriteFile(outputPath, inputBytes, 0644)
	}

	s.log.Info("PDF unlocked successfully", logger.String("output", outputPath))
	return nil
}

// UnlockBytes removes password from PDF bytes
func (s *unlockService) UnlockBytes(input []byte, password string) ([]byte, error) {
	s.log.Info("UnlockService.UnlockBytes called")

	// Create temp files
	tmpInput, err := os.CreateTemp("", "pdf-unlock-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())

	tmpOutput, err := os.CreateTemp("", "pdf-unlock-out-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpOutput.Name())
	tmpOutput.Close()

	// Write input
	if _, err := tmpInput.Write(input); err != nil {
		tmpInput.Close()
		return nil, err
	}
	tmpInput.Close()

	// Unlock
	if err := s.UnlockFile(tmpInput.Name(), tmpOutput.Name(), password); err != nil {
		return nil, err
	}

	// Read output
	output, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, err
	}

	s.log.Info("PDF unlocked", logger.Int("outputSize", len(output)))
	return output, nil
}
