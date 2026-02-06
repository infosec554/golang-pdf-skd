package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

type PDFToJPGService interface {
	Convert(input io.Reader) ([]byte, error)
	ConvertFile(inputPath, outputDir string) ([]string, error)
	ConvertBytes(input []byte) ([]byte, error)
	ConvertToImages(input []byte) ([][]byte, error)
}

type pdfToJPGService struct {
	log logger.ILogger
}

func NewPDFToJPGService(log logger.ILogger) PDFToJPGService {
	return &pdfToJPGService{
		log: log,
	}
}

func (s *pdfToJPGService) Convert(input io.Reader) ([]byte, error) {
	s.log.Info("PDFToJPGService.Convert called")

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		s.log.Error("Failed to read input", logger.Error(err))
		return nil, err
	}

	return s.ConvertBytes(inputBytes)
}

func (s *pdfToJPGService) ConvertFile(inputPath, outputDir string) ([]string, error) {
	s.log.Info("PDFToJPGService.ConvertFile called", logger.String("input", inputPath))

	if err := os.MkdirAll(outputDir, 0777); err != nil {
		s.log.Error("Failed to create output dir", logger.Error(err))
		return nil, err
	}

	prefix := filepath.Join(outputDir, "page")
	cmd := exec.Command("pdftoppm", "-jpeg", "-r", "150", inputPath, prefix)

	if output, err := cmd.CombinedOutput(); err != nil {
		s.log.Error("pdftoppm execution failed", logger.String("output", string(output)), logger.Error(err))
		return nil, fmt.Errorf("image conversion failed: %w", err)
	}

	var imageFiles []string
	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".jpg" || ext == ".jpeg" {
				imageFiles = append(imageFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(imageFiles)
	s.log.Info("PDF to JPG conversion completed", logger.Int("pages", len(imageFiles)))

	return imageFiles, nil
}

func (s *pdfToJPGService) ConvertBytes(input []byte) ([]byte, error) {
	s.log.Info("PDFToJPGService.ConvertBytes called")

	tmpDir, err := os.MkdirTemp("", "pdf-to-jpg-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	tmpInput := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(tmpInput, input, 0644); err != nil {
		return nil, err
	}

	outputDir := filepath.Join(tmpDir, "output")
	imageFiles, err := s.ConvertFile(tmpInput, outputDir)
	if err != nil {
		return nil, err
	}

	if len(imageFiles) == 0 {
		return nil, fmt.Errorf("no images generated")
	}

	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	for i, imgPath := range imageFiles {
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			continue
		}

		fileName := fmt.Sprintf("page_%d.jpg", i+1)
		w, err := zipWriter.Create(fileName)
		if err != nil {
			continue
		}
		w.Write(imgData)
	}
	zipWriter.Close()

	s.log.Info("PDF to JPG ZIP created", logger.Int("pages", len(imageFiles)))
	return zipBuffer.Bytes(), nil
}

func (s *pdfToJPGService) ConvertToImages(input []byte) ([][]byte, error) {
	s.log.Info("PDFToJPGService.ConvertToImages called")

	tmpDir, err := os.MkdirTemp("", "pdf-to-jpg-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	tmpInput := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(tmpInput, input, 0644); err != nil {
		return nil, err
	}

	outputDir := filepath.Join(tmpDir, "output")
	imageFiles, err := s.ConvertFile(tmpInput, outputDir)
	if err != nil {
		return nil, err
	}

	var images [][]byte
	for _, imgPath := range imageFiles {
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			continue
		}
		images = append(images, imgData)
	}

	s.log.Info("PDF to images conversion completed", logger.Int("pages", len(images)))
	return images, nil
}
