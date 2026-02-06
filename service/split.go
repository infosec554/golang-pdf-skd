package service

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
)

type SplitService interface {
	Split(input io.Reader, ranges string) ([]byte, error)
	SplitFile(inputPath, outputDir string, ranges string) ([]string, error)
	SplitBytes(input []byte, ranges string) ([]byte, error)
	SplitToPages(input []byte) ([][]byte, error)
}

type splitService struct {
	log logger.ILogger
}

func NewSplitService(log logger.ILogger) SplitService {
	return &splitService{
		log: log,
	}
}

func (s *splitService) Split(input io.Reader, ranges string) ([]byte, error) {
	s.log.Info("SplitService.Split called", logger.String("ranges", ranges))

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return s.SplitBytes(inputBytes, ranges)
}

func (s *splitService) SplitFile(inputPath, outputDir string, ranges string) ([]string, error) {
	s.log.Info("SplitService.SplitFile called", logger.String("input", inputPath))

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, err
	}

	rangeList := strings.Split(ranges, ",")
	var outputFiles []string

	for i, r := range rangeList {
		r = strings.TrimSpace(r)
		partDir := filepath.Join(outputDir, fmt.Sprintf("part_%d", i+1))
		if err := os.MkdirAll(partDir, os.ModePerm); err != nil {
			continue
		}

		if err := api.ExtractPagesFile(inputPath, partDir, []string{r}, nil); err != nil {
			s.log.Error("pdfcpu extract failed", logger.String("range", r), logger.Error(err))
			continue
		}

		files, err := os.ReadDir(partDir)
		if err != nil || len(files) == 0 {
			continue
		}

		generatedPath := filepath.Join(partDir, files[0].Name())
		outputPath := filepath.Join(outputDir, fmt.Sprintf("part_%d.pdf", i+1))

		if err := os.Rename(generatedPath, outputPath); err != nil {
			continue
		}
		os.RemoveAll(partDir)

		outputFiles = append(outputFiles, outputPath)
	}

	s.log.Info("PDF split completed", logger.Int("parts", len(outputFiles)))
	return outputFiles, nil
}

func (s *splitService) SplitBytes(input []byte, ranges string) ([]byte, error) {
	s.log.Info("SplitService.SplitBytes called")

	tmpDir, err := os.MkdirTemp("", "pdf-split-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	outputDir := filepath.Join(tmpDir, "output")
	outputFiles, err := s.SplitFile(inputPath, outputDir, ranges)
	if err != nil {
		return nil, err
	}

	if len(outputFiles) == 0 {
		return nil, fmt.Errorf("no output files generated")
	}

	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	for _, filePath := range outputFiles {
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		w, err := zipWriter.Create(filepath.Base(filePath))
		if err != nil {
			continue
		}
		w.Write(data)
	}
	zipWriter.Close()

	return zipBuffer.Bytes(), nil
}

func (s *splitService) SplitToPages(input []byte) ([][]byte, error) {
	s.log.Info("SplitService.SplitToPages called")

	tmpDir, err := os.MkdirTemp("", "pdf-split-pages-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	outputDir := filepath.Join(tmpDir, "output")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, err
	}

	if err := api.SplitFile(inputPath, outputDir, 1, nil); err != nil {
		s.log.Error("pdfcpu split failed", logger.Error(err))
		return nil, err
	}

	files, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, err
	}

	var pages [][]byte
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(outputDir, f.Name()))
		if err != nil {
			continue
		}
		pages = append(pages, data)
	}

	s.log.Info("PDF split to pages completed", logger.Int("pages", len(pages)))
	return pages, nil
}
