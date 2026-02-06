package service

import (
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

type UnlockService interface {
	Unlock(input io.Reader, password string) ([]byte, error)
	UnlockFile(inputPath, outputPath, password string) error
	UnlockBytes(input []byte, password string) ([]byte, error)
}

type unlockService struct {
	log logger.ILogger
}

func NewUnlockService(log logger.ILogger) UnlockService {
	return &unlockService{
		log: log,
	}
}

func (s *unlockService) Unlock(input io.Reader, password string) ([]byte, error) {
	s.log.Info("UnlockService.Unlock called")

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return s.UnlockBytes(inputBytes, password)
}

func (s *unlockService) UnlockFile(inputPath, outputPath, password string) error {
	s.log.Info("UnlockService.UnlockFile called", logger.String("input", inputPath))

	conf := api.LoadConfiguration()
	if password != "" {
		conf.UserPW = password
		conf.OwnerPW = password
	}

	if err := api.DecryptFile(inputPath, outputPath, conf); err != nil {
		s.log.Error("pdfcpu decrypt failed, copying as-is", logger.Error(err))
		inputBytes, err := os.ReadFile(inputPath)
		if err != nil {
			return err
		}
		return os.WriteFile(outputPath, inputBytes, 0644)
	}

	s.log.Info("PDF unlocked successfully", logger.String("output", outputPath))
	return nil
}

func (s *unlockService) UnlockBytes(input []byte, password string) ([]byte, error) {
	s.log.Info("UnlockService.UnlockBytes called")

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

	if _, err := tmpInput.Write(input); err != nil {
		tmpInput.Close()
		return nil, err
	}
	tmpInput.Close()

	if err := s.UnlockFile(tmpInput.Name(), tmpOutput.Name(), password); err != nil {
		return nil, err
	}

	output, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, err
	}

	s.log.Info("PDF unlocked", logger.Int("outputSize", len(output)))
	return output, nil
}
