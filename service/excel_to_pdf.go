package service

import (
	"context"
	"io"
	"os"

	"github.com/infosec554/golang-pdf-sdk/pkg/gotenberg"
	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

// ExcelToPDFService converts Excel spreadsheets to PDF
type ExcelToPDFService interface {
	// Convert converts Excel to PDF bytes
	Convert(ctx context.Context, input io.Reader, filename string) ([]byte, error)
	// ConvertFile converts Excel file to PDF file
	ConvertFile(ctx context.Context, inputPath, outputPath string) error
	// ConvertBytes converts Excel bytes to PDF bytes
	ConvertBytes(ctx context.Context, input []byte, filename string) ([]byte, error)
}

type excelToPDFService struct {
	log       logger.ILogger
	gotClient gotenberg.Client
}

// NewExcelToPDFService creates a new Excel to PDF service
func NewExcelToPDFService(log logger.ILogger, gotClient gotenberg.Client) ExcelToPDFService {
	return &excelToPDFService{
		log:       log,
		gotClient: gotClient,
	}
}

// Convert converts Excel from reader to PDF bytes
func (s *excelToPDFService) Convert(ctx context.Context, input io.Reader, filename string) ([]byte, error) {
	s.log.Info("ExcelToPDFService.Convert called", logger.String("filename", filename))

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		s.log.Error("Failed to read input", logger.Error(err))
		return nil, err
	}

	return s.ConvertBytes(ctx, inputBytes, filename)
}

// ConvertFile converts Excel file to PDF file
func (s *excelToPDFService) ConvertFile(ctx context.Context, inputPath, outputPath string) error {
	s.log.Info("ExcelToPDFService.ConvertFile called", logger.String("input", inputPath))

	resultBytes, err := s.gotClient.ExcelToPDF(ctx, inputPath)
	if err != nil {
		s.log.Error("Gotenberg conversion failed", logger.Error(err))
		return err
	}

	if err := os.WriteFile(outputPath, resultBytes, 0644); err != nil {
		s.log.Error("Failed to write output file", logger.Error(err))
		return err
	}

	s.log.Info("Excel to PDF conversion completed", logger.String("output", outputPath))
	return nil
}

// ConvertBytes converts Excel bytes to PDF bytes
func (s *excelToPDFService) ConvertBytes(ctx context.Context, input []byte, filename string) ([]byte, error) {
	s.log.Info("ExcelToPDFService.ConvertBytes called")

	// Create temp file for input
	ext := ".xlsx"
	if len(filename) > 4 {
		ext = filename[len(filename)-5:]
	}
	tmpInput, err := os.CreateTemp("", "excel-input-*"+ext)
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
	resultBytes, err := s.gotClient.ExcelToPDF(ctx, tmpInput.Name())
	if err != nil {
		s.log.Error("Gotenberg conversion failed", logger.Error(err))
		return nil, err
	}

	s.log.Info("Excel to PDF conversion completed", logger.Int("outputSize", len(resultBytes)))
	return resultBytes, nil
}
