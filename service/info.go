package service

import (
	"io"
	"os"

	"github.com/pdfcpu/pdfcpu/pkg/api"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

type PDFInfo struct {
	PageCount int
	Version   string
	Encrypted bool
	FileSize  int64
	IsValid   bool
}

type InfoService interface {
	GetInfo(input io.Reader) (*PDFInfo, error)
	GetInfoFile(inputPath string) (*PDFInfo, error)
	GetInfoBytes(input []byte) (*PDFInfo, error)
	GetPageCount(input []byte) (int, error)
	ValidatePDF(input []byte) error
	IsEncrypted(input []byte) (bool, error)
}

type infoService struct {
	log logger.ILogger
}

func NewInfoService(log logger.ILogger) InfoService {
	return &infoService{log: log}
}

func (s *infoService) GetInfo(input io.Reader) (*PDFInfo, error) {
	s.log.Info("InfoService.GetInfo called")

	inputBytes, err := io.ReadAll(input)
	if err != nil {
		return nil, err
	}

	return s.GetInfoBytes(inputBytes)
}

func (s *infoService) GetInfoFile(inputPath string) (*PDFInfo, error) {
	s.log.Info("InfoService.GetInfoFile called", logger.String("input", inputPath))

	ctx, err := api.ReadContextFile(inputPath)
	if err != nil {
		return nil, err
	}

	info := &PDFInfo{
		PageCount: ctx.PageCount,
		Version:   ctx.HeaderVersion.String(),
		Encrypted: ctx.Encrypt != nil,
		IsValid:   true,
	}

	if stat, err := os.Stat(inputPath); err == nil {
		info.FileSize = stat.Size()
	}

	s.log.Info("PDF info retrieved", logger.Int("pages", info.PageCount))
	return info, nil
}

func (s *infoService) GetInfoBytes(input []byte) (*PDFInfo, error) {
	s.log.Info("InfoService.GetInfoBytes called")

	tmpFile, err := os.CreateTemp("", "pdf-info-*.pdf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(input); err != nil {
		tmpFile.Close()
		return nil, err
	}
	tmpFile.Close()

	info, err := s.GetInfoFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	info.FileSize = int64(len(input))
	return info, nil
}

func (s *infoService) GetPageCount(input []byte) (int, error) {
	s.log.Info("InfoService.GetPageCount called")

	info, err := s.GetInfoBytes(input)
	if err != nil {
		return 0, err
	}

	return info.PageCount, nil
}

func (s *infoService) ValidatePDF(input []byte) error {
	s.log.Info("InfoService.ValidatePDF called")

	tmpFile, err := os.CreateTemp("", "pdf-validate-*.pdf")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(input); err != nil {
		tmpFile.Close()
		return err
	}
	tmpFile.Close()

	return api.ValidateFile(tmpFile.Name(), nil)
}

func (s *infoService) IsEncrypted(input []byte) (bool, error) {
	s.log.Info("InfoService.IsEncrypted called")

	info, err := s.GetInfoBytes(input)
	if err != nil {
		return false, err
	}

	return info.Encrypted, nil
}
