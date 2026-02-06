package service

import (
	"fmt"
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

// WatermarkService adds watermarks to PDF documents
type WatermarkService interface {
	// AddWatermark adds text watermark to PDF
	AddWatermark(input io.Reader, text string, options *WatermarkOptions) ([]byte, error)
	// AddWatermarkFile adds watermark to PDF file
	AddWatermarkFile(inputPath, outputPath, text string, options *WatermarkOptions) error
	// AddWatermarkBytes adds watermark to PDF bytes
	AddWatermarkBytes(input []byte, text string, options *WatermarkOptions) ([]byte, error)
}

// WatermarkOptions configures watermark appearance
type WatermarkOptions struct {
	FontSize int     // Default: 48
	Position string  // "diagonal", "center", "top", "bottom"
	Opacity  float64 // 0.0 to 1.0, default: 0.3
	Color    string  // Default: "gray"
}

type watermarkService struct {
	log logger.ILogger
}

// NewWatermarkService creates a new watermark service
func NewWatermarkService(log logger.ILogger) WatermarkService {
	return &watermarkService{
		log: log,
	}
}

// DefaultWatermarkOptions returns default watermark options
func DefaultWatermarkOptions() *WatermarkOptions {
	return &WatermarkOptions{
		FontSize: 48,
		Position: "diagonal",
		Opacity:  0.3,
		Color:    "gray",
	}
}

// AddWatermark adds watermark from reader
func (s *watermarkService) AddWatermark(input io.Reader, text string, options *WatermarkOptions) ([]byte, error) {
	s.log.Info("WatermarkService.AddWatermark called", logger.String("text", text))

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return s.AddWatermarkBytes(inputBytes, text, options)
}

// AddWatermarkFile adds watermark to file
func (s *watermarkService) AddWatermarkFile(inputPath, outputPath, text string, options *WatermarkOptions) error {
	s.log.Info("WatermarkService.AddWatermarkFile called", logger.String("input", inputPath))

	if options == nil {
		options = DefaultWatermarkOptions()
	}

	// Copy input to output first
	inputBytes, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, inputBytes, 0644); err != nil {
		return err
	}

	// Calculate rotation based on position
	rotation := 45
	if options.Position == "center" {
		rotation = 0
	}

	// Build watermark description
	wmDesc := fmt.Sprintf("%s, font:Helvetica, points:%d, color:%s, opacity:%.1f, rotation:%d, scale:1.0 abs, position:c",
		text, options.FontSize, options.Color, options.Opacity, rotation)

	wm, err := api.TextWatermark(wmDesc, "", true, false, types.POINTS)
	if err != nil {
		s.log.Error("Failed to create watermark", logger.Error(err))
		os.Remove(outputPath)
		return fmt.Errorf("watermark create failed: %w", err)
	}

	if err := api.AddWatermarksFile(outputPath, "", nil, wm, nil); err != nil {
		s.log.Error("pdfcpu watermark failed", logger.Error(err))
		os.Remove(outputPath)
		return fmt.Errorf("watermark failed: %w", err)
	}

	s.log.Info("Watermark added successfully", logger.String("output", outputPath))
	return nil
}

// AddWatermarkBytes adds watermark to PDF bytes
func (s *watermarkService) AddWatermarkBytes(input []byte, text string, options *WatermarkOptions) ([]byte, error) {
	s.log.Info("WatermarkService.AddWatermarkBytes called")

	if options == nil {
		options = DefaultWatermarkOptions()
	}

	// Create temp files
	tmpInput, err := os.CreateTemp("", "pdf-watermark-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpInput.Name())

	tmpOutput, err := os.CreateTemp("", "pdf-watermark-out-*.pdf")
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

	// Add watermark
	if err := s.AddWatermarkFile(tmpInput.Name(), tmpOutput.Name(), text, options); err != nil {
		return nil, err
	}

	// Read output
	output, err := os.ReadFile(tmpOutput.Name())
	if err != nil {
		return nil, err
	}

	s.log.Info("Watermark added", logger.Int("outputSize", len(output)))
	return output, nil
}
