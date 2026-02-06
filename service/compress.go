package service

import (
	"bytes"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

type CompressService interface {
	Compress(input io.Reader) ([]byte, error)
	CompressFile(inputPath, outputPath string) error
	CompressBytes(input []byte) ([]byte, error)
}

type compressService struct {
	log logger.ILogger
}

func NewCompressService(log logger.ILogger) CompressService {
	return &compressService{
		log: log,
	}
}

func (s *compressService) Compress(input io.Reader) ([]byte, error) {
	s.log.Info("CompressService.Compress called")

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		s.log.Error("Failed to read input", logger.Error(err))
		return nil, err
	}

	return s.CompressBytes(inputBytes)
}

func (s *compressService) CompressFile(inputPath, outputPath string) error {
	s.log.Info("CompressService.CompressFile called", logger.String("input", inputPath))

	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.OPTIMIZE

	if err := api.OptimizeFile(inputPath, outputPath, conf); err != nil {
		s.log.Error("pdfcpu optimize failed", logger.Error(err))
		return err
	}

	s.log.Info("PDF compression completed", logger.String("output", outputPath))
	return nil
}

func (s *compressService) CompressBytes(input []byte) ([]byte, error) {
	s.log.Info("CompressService.CompressBytes called")

	tmpInput, err := os.CreateTemp("", "pdf-compress-input-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())
	defer tmpInput.Close()

	tmpOutput, err := os.CreateTemp("", "pdf-compress-output-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpOutput.Name())
	tmpOutput.Close()

	if _, err := tmpInput.Write(input); err != nil {
		return nil, err
	}
	tmpInput.Close()

	if err := s.CompressFile(tmpInput.Name(), tmpOutput.Name()); err != nil {
		return nil, err
	}

	output, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, err
	}

	ratio := float64(len(output)) / float64(len(input)) * 100
	s.log.Info("Compression completed",
		logger.Int("inputSize", len(input)),
		logger.Int("outputSize", len(output)),
		logger.String("ratio", bytes.NewBufferString("").String()),
	)
	_ = ratio

	return output, nil
}
