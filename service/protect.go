package service

import (
	"fmt"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

// ProtectService adds password protection to PDF documents
type ProtectService interface {
	// Protect adds password to PDF
	Protect(input io.Reader, password string) ([]byte, error)
	// ProtectFile adds password to PDF file
	ProtectFile(inputPath, outputPath, password string) error
	// ProtectBytes adds password to PDF bytes
	ProtectBytes(input []byte, password string) ([]byte, error)
}

type protectService struct {
	log logger.ILogger
}

// NewProtectService creates a new protect service
func NewProtectService(log logger.ILogger) ProtectService {
	return &protectService{
		log: log,
	}
}

// Protect adds password to PDF from reader
func (s *protectService) Protect(input io.Reader, password string) ([]byte, error) {
	s.log.Info("ProtectService.Protect called")

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return s.ProtectBytes(inputBytes, password)
}

// ProtectFile adds password to PDF file
func (s *protectService) ProtectFile(inputPath, outputPath, password string) error {
	s.log.Info("ProtectService.ProtectFile called", logger.String("input", inputPath))

	conf := api.LoadConfiguration()
	conf.UserPW = password
	conf.OwnerPW = password
	conf.EncryptUsingAES = true
	conf.EncryptKeyLength = 256
	conf.Permissions = model.PermissionsAll

	if err := api.EncryptFile(inputPath, outputPath, conf); err != nil {
		s.log.Error("pdfcpu encrypt failed", logger.Error(err))
		return fmt.Errorf("encryption failed: %w", err)
	}

	s.log.Info("PDF protected successfully", logger.String("output", outputPath))
	return nil
}

// ProtectBytes adds password to PDF bytes
func (s *protectService) ProtectBytes(input []byte, password string) ([]byte, error) {
	s.log.Info("ProtectService.ProtectBytes called")

	// Create temp files
	tmpInput, err := os.CreateTemp("", "pdf-protect-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())

	tmpOutput, err := os.CreateTemp("", "pdf-protect-out-*.pdf")
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

	// Protect
	if err := s.ProtectFile(tmpInput.Name(), tmpOutput.Name(), password); err != nil {
		return nil, err
	}

	// Read output
	output, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, err
	}

	s.log.Info("PDF protected", logger.Int("outputSize", len(output)))
	return output, nil
}
