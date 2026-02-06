package service

import (
	"context"
	"io"
	"os"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/gotenberg"
	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

type PowerPointToPDFService interface {
	Convert(ctx context.Context, input io.Reader, filename string) ([]byte, error)
	ConvertFile(ctx context.Context, inputPath, outputPath string) error
	ConvertBytes(ctx context.Context, input []byte, filename string) ([]byte, error)
}

type powerPointToPDFService struct {
	log       logger.ILogger
	gotClient gotenberg.Client
}

func NewPowerPointToPDFService(log logger.ILogger, gotClient gotenberg.Client) PowerPointToPDFService {
	return &powerPointToPDFService{
		log:       log,
		gotClient: gotClient,
	}
}

func (s *powerPointToPDFService) Convert(ctx context.Context, input io.Reader, filename string) ([]byte, error) {
	s.log.Info("PowerPointToPDFService.Convert called", logger.String("filename", filename))

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		s.log.Error("Failed to read input", logger.Error(err))
		return nil, err
	}

	return s.ConvertBytes(ctx, inputBytes, filename)
}

func (s *powerPointToPDFService) ConvertFile(ctx context.Context, inputPath, outputPath string) error {
	s.log.Info("PowerPointToPDFService.ConvertFile called", logger.String("input", inputPath))

	resultBytes, err := s.gotClient.PowerPointToPDF(ctx, inputPath)
	if err != nil {
		s.log.Error("Gotenberg conversion failed", logger.Error(err))
		return err
	}

	if err := os.WriteFile(outputPath, resultBytes, 0644); err != nil {
		s.log.Error("Failed to write output file", logger.Error(err))
		return err
	}

	s.log.Info("PowerPoint to PDF conversion completed", logger.String("output", outputPath))
	return nil
}

func (s *powerPointToPDFService) ConvertBytes(ctx context.Context, input []byte, filename string) ([]byte, error) {
	s.log.Info("PowerPointToPDFService.ConvertBytes called")

	ext := ".pptx"
	if len(filename) > 4 {
		ext = filename[len(filename)-5:]
	}
	tmpInput, err := os.CreateTemp("", "ppt-input-*"+ext)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())

	if _, err := tmpInput.Write(input); err != nil {
		tmpInput.Close()
		return nil, err
	}
	tmpInput.Close()

	resultBytes, err := s.gotClient.PowerPointToPDF(ctx, tmpInput.Name())
	if err != nil {
		s.log.Error("Gotenberg conversion failed", logger.Error(err))
		return nil, err
	}

	s.log.Info("PowerPoint to PDF conversion completed", logger.Int("outputSize", len(resultBytes)))
	return resultBytes, nil
}
