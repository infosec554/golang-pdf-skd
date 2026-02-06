package service

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

type PageService interface {
	ExtractPages(input []byte, pages string) ([]byte, error)
	DeletePages(input []byte, pages string) ([]byte, error)
	InsertPages(base []byte, insert []byte, afterPage int) ([]byte, error)
	ReorderPages(input []byte, order []int) ([]byte, error)
	GetPageCount(input []byte) (int, error)
}

type pageService struct {
	log logger.ILogger
}

func NewPageService(log logger.ILogger) PageService {
	return &pageService{log: log}
}

func (s *pageService) ExtractPages(input []byte, pages string) ([]byte, error) {
	s.log.Info("PageService.ExtractPages called", logger.String("pages", pages))

	tmpDir, err := os.MkdirTemp("", "pdf-extract-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")

	if err := api.TrimFile(inputPath, outputPath, []string{pages}, nil); err != nil {
		s.log.Error("pdfcpu trim failed", logger.Error(err))
		return nil, err
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("Pages extracted", logger.Int("outputSize", len(output)))
	return output, nil
}

func (s *pageService) DeletePages(input []byte, pages string) ([]byte, error) {
	s.log.Info("PageService.DeletePages called", logger.String("pages", pages))

	tmpDir, err := os.MkdirTemp("", "pdf-delete-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	if err := api.RemovePagesFile(inputPath, "", []string{pages}, nil); err != nil {
		s.log.Error("pdfcpu remove failed", logger.Error(err))
		return nil, err
	}

	output, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("Pages deleted", logger.Int("outputSize", len(output)))
	return output, nil
}

func (s *pageService) InsertPages(base []byte, insert []byte, afterPage int) ([]byte, error) {
	s.log.Info("PageService.InsertPages called", logger.Int("afterPage", afterPage))

	tmpDir, err := os.MkdirTemp("", "pdf-insert-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	basePath := filepath.Join(tmpDir, "base.pdf")
	insertPath := filepath.Join(tmpDir, "insert.pdf")
	outputPath := filepath.Join(tmpDir, "output.pdf")

	if err := os.WriteFile(basePath, base, 0644); err != nil {
		return nil, err
	}
	if err := os.WriteFile(insertPath, insert, 0644); err != nil {
		return nil, err
	}

	baseCtx, err := api.ReadContextFile(basePath)
	if err != nil {
		return nil, err
	}

	if afterPage <= 0 {
		conf := model.NewDefaultConfiguration()
		if err := api.MergeCreateFile([]string{insertPath, basePath}, outputPath, false, conf); err != nil {
			return nil, err
		}
	} else if afterPage >= baseCtx.PageCount {
		conf := model.NewDefaultConfiguration()
		if err := api.MergeCreateFile([]string{basePath, insertPath}, outputPath, false, conf); err != nil {
			return nil, err
		}
	} else {
		part1Path := filepath.Join(tmpDir, "part1.pdf")
		part2Path := filepath.Join(tmpDir, "part2.pdf")

		if err := api.TrimFile(basePath, part1Path, []string{fmt.Sprintf("1-%d", afterPage)}, nil); err != nil {
			return nil, err
		}

		if err := api.TrimFile(basePath, part2Path, []string{fmt.Sprintf("%d-", afterPage+1)}, nil); err != nil {
			return nil, err
		}

		conf := model.NewDefaultConfiguration()
		if err := api.MergeCreateFile([]string{part1Path, insertPath, part2Path}, outputPath, false, conf); err != nil {
			return nil, err
		}
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("Pages inserted", logger.Int("outputSize", len(output)))
	return output, nil
}

func (s *pageService) ReorderPages(input []byte, order []int) ([]byte, error) {
	s.log.Info("PageService.ReorderPages called", logger.Int("newOrderLen", len(order)))

	tmpDir, err := os.MkdirTemp("", "pdf-reorder-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	pagesDir := filepath.Join(tmpDir, "pages")
	if err := os.MkdirAll(pagesDir, 0755); err != nil {
		return nil, err
	}

	if err := api.SplitFile(inputPath, pagesDir, 1, nil); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(pagesDir)
	if err != nil {
		return nil, err
	}

	var pageFiles []string
	for _, f := range files {
		if !f.IsDir() {
			pageFiles = append(pageFiles, filepath.Join(pagesDir, f.Name()))
		}
	}
	sort.Strings(pageFiles)

	var reorderedFiles []string
	for _, pageNum := range order {
		if pageNum > 0 && pageNum <= len(pageFiles) {
			reorderedFiles = append(reorderedFiles, pageFiles[pageNum-1])
		}
	}

	outputPath := filepath.Join(tmpDir, "output.pdf")
	conf := model.NewDefaultConfiguration()
	if err := api.MergeCreateFile(reorderedFiles, outputPath, false, conf); err != nil {
		return nil, err
	}

	output, err := os.ReadFile(outputPath)
	if err != nil {
		return nil, err
	}

	s.log.Info("Pages reordered", logger.Int("outputSize", len(output)))
	return output, nil
}

func (s *pageService) GetPageCount(input []byte) (int, error) {
	s.log.Info("PageService.GetPageCount called")

	tmpFile, err := os.CreateTemp("", "pdf-count-*.pdf")
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(input); err != nil {
		tmpFile.Close()
		return 0, err
	}
	tmpFile.Close()

	ctx, err := api.ReadContextFile(tmpFile.Name())
	if err != nil {
		return 0, err
	}

	return ctx.PageCount, nil
}

type TextService interface {
	ExtractText(input []byte) (string, error)
	ExtractTextFromPage(input []byte, page int) (string, error)
}

type textService struct {
	log logger.ILogger
}

func NewTextService(log logger.ILogger) TextService {
	return &textService{log: log}
}

func (s *textService) ExtractText(input []byte) (string, error) {
	s.log.Info("TextService.ExtractText called")

	tmpFile, err := os.CreateTemp("", "pdf-text-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(input); err != nil {
		tmpFile.Close()
		return "", err
	}
	tmpFile.Close()

	textFile, err := os.CreateTemp("", "pdf-text-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(textFile.Name())
	textFile.Close()

	if err := api.ExtractContentFile(tmpFile.Name(), textFile.Name(), nil, nil); err != nil {
		s.log.Error("pdfcpu extract text failed", logger.Error(err))
		return "", err
	}

	textBytes, err := os.ReadFile(textFile.Name())
	if err != nil {
		return "", err
	}

	s.log.Info("Text extracted", logger.Int("length", len(textBytes)))
	return string(textBytes), nil
}

func (s *textService) ExtractTextFromPage(input []byte, page int) (string, error) {
	s.log.Info("TextService.ExtractTextFromPage called", logger.Int("page", page))

	pageService := &pageService{log: s.log}
	pageBytes, err := pageService.ExtractPages(input, string(rune('0'+page)))
	if err != nil {
		return "", err
	}

	return s.ExtractText(pageBytes)
}

type MetadataService interface {
	GetMetadata(input []byte) (map[string]string, error)
	SetMetadata(input []byte, metadata map[string]string) ([]byte, error)
}

type metadataService struct {
	log logger.ILogger
}

func NewMetadataService(log logger.ILogger) MetadataService {
	return &metadataService{log: log}
}

func (s *metadataService) GetMetadata(input []byte) (map[string]string, error) {
	s.log.Info("MetadataService.GetMetadata called")

	tmpFile, err := os.CreateTemp("", "pdf-meta-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(input); err != nil {
		tmpFile.Close()
		return nil, err
	}
	tmpFile.Close()

	ctx, err := api.ReadContextFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	metadata := make(map[string]string)
	metadata["pages"] = string(rune('0' + ctx.PageCount))
	metadata["version"] = ctx.HeaderVersion.String()
	metadata["encrypted"] = "false"
	if ctx.Encrypt != nil {
		metadata["encrypted"] = "true"
	}

	s.log.Info("Metadata retrieved")
	return metadata, nil
}

func (s *metadataService) SetMetadata(input []byte, metadata map[string]string) ([]byte, error) {
	s.log.Info("MetadataService.SetMetadata called")
	s.log.Warn("SetMetadata has limited support - returning input unchanged")

	output := make([]byte, len(input))
	copy(output, input)

	s.log.Info("Metadata set (limited support)", logger.Int("outputSize", len(output)))
	return output, nil
}

type ImageExtractService interface {
	ExtractImages(input []byte) ([][]byte, error)
	ExtractImagesFromPage(input []byte, page int) ([][]byte, error)
}

type imageExtractService struct {
	log logger.ILogger
}

func NewImageExtractService(log logger.ILogger) ImageExtractService {
	return &imageExtractService{log: log}
}

func (s *imageExtractService) ExtractImages(input []byte) ([][]byte, error) {
	s.log.Info("ImageExtractService.ExtractImages called")

	tmpDir, err := os.MkdirTemp("", "pdf-img-extract-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	inputPath := filepath.Join(tmpDir, "input.pdf")
	if err := os.WriteFile(inputPath, input, 0644); err != nil {
		return nil, err
	}

	outputDir := filepath.Join(tmpDir, "images")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}

	if err := api.ExtractImagesFile(inputPath, outputDir, nil, nil); err != nil {
		s.log.Error("pdfcpu extract images failed", logger.Error(err))
		return nil, err
	}

	files, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, err
	}

	var images [][]byte
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		imgData, err := os.ReadFile(filepath.Join(outputDir, f.Name()))
		if err != nil {
			continue
		}
		images = append(images, imgData)
	}

	s.log.Info("Images extracted", logger.Int("count", len(images)))
	return images, nil
}

func (s *imageExtractService) ExtractImagesFromPage(input []byte, page int) ([][]byte, error) {
	s.log.Info("ImageExtractService.ExtractImagesFromPage called", logger.Int("page", page))

	pageService := &pageService{log: s.log}
	pageBytes, err := pageService.ExtractPages(input, string(rune('0'+page)))
	if err != nil {
		return nil, err
	}

	return s.ExtractImages(pageBytes)
}
