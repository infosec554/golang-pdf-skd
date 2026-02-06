package service

import (
	"github.com/infosec554/convert-pdf-go-sdk/pkg/gotenberg"
	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
)

type PDFService interface {
	WordToPDF() WordToPDFService
	ExcelToPDF() ExcelToPDFService
	PowerPointToPDF() PowerPointToPDFService
	JPGToPDF() JPGToPDFService
	PDFToJPG() PDFToJPGService
	Compress() CompressService
	Merge() MergeService
	Split() SplitService
	Rotate() RotateService
	Watermark() WatermarkService
	Protect() ProtectService
	Unlock() UnlockService
	Info() InfoService
	Pages() PageService
	Text() TextService
	Metadata() MetadataService
	Images() ImageExtractService
	Batch(maxWorkers int) *BatchProcessor
	Pipeline() *Pipeline
}

type pdfService struct {
	wordToPDF       WordToPDFService
	excelToPDF      ExcelToPDFService
	powerPointToPDF PowerPointToPDFService
	jpgToPDF        JPGToPDFService
	pdfToJPG        PDFToJPGService
	compress        CompressService
	merge           MergeService
	split           SplitService
	rotate          RotateService
	watermark       WatermarkService
	protect         ProtectService
	unlock          UnlockService
	info            InfoService
	pages           PageService
	text            TextService
	metadata        MetadataService
	images          ImageExtractService
	log             logger.ILogger
	gotClient       gotenberg.Client
}

func New(log logger.ILogger, gotClient gotenberg.Client) PDFService {
	return &pdfService{
		wordToPDF:       NewWordToPDFService(log, gotClient),
		excelToPDF:      NewExcelToPDFService(log, gotClient),
		powerPointToPDF: NewPowerPointToPDFService(log, gotClient),
		jpgToPDF:        NewJPGToPDFService(log),
		pdfToJPG:        NewPDFToJPGService(log),
		compress:        NewCompressService(log),
		merge:           NewMergeService(log),
		split:           NewSplitService(log),
		rotate:          NewRotateService(log),
		watermark:       NewWatermarkService(log),
		protect:         NewProtectService(log),
		unlock:          NewUnlockService(log),
		info:            NewInfoService(log),
		pages:           NewPageService(log),
		text:            NewTextService(log),
		metadata:        NewMetadataService(log),
		images:          NewImageExtractService(log),
		log:             log,
		gotClient:       gotClient,
	}
}

func NewWithGotenberg(gotenbergURL string) PDFService {
	log := logger.New("golang-pdf-sdk")
	gotClient := gotenberg.New(gotenbergURL)
	return New(log, gotClient)
}

func (s *pdfService) WordToPDF() WordToPDFService             { return s.wordToPDF }
func (s *pdfService) ExcelToPDF() ExcelToPDFService           { return s.excelToPDF }
func (s *pdfService) PowerPointToPDF() PowerPointToPDFService { return s.powerPointToPDF }
func (s *pdfService) JPGToPDF() JPGToPDFService               { return s.jpgToPDF }
func (s *pdfService) PDFToJPG() PDFToJPGService               { return s.pdfToJPG }
func (s *pdfService) Compress() CompressService               { return s.compress }
func (s *pdfService) Merge() MergeService                     { return s.merge }
func (s *pdfService) Split() SplitService                     { return s.split }
func (s *pdfService) Rotate() RotateService                   { return s.rotate }
func (s *pdfService) Watermark() WatermarkService             { return s.watermark }
func (s *pdfService) Protect() ProtectService                 { return s.protect }
func (s *pdfService) Unlock() UnlockService                   { return s.unlock }
func (s *pdfService) Info() InfoService                       { return s.info }
func (s *pdfService) Pages() PageService                      { return s.pages }
func (s *pdfService) Text() TextService                       { return s.text }
func (s *pdfService) Metadata() MetadataService               { return s.metadata }
func (s *pdfService) Images() ImageExtractService             { return s.images }

func (s *pdfService) Batch(maxWorkers int) *BatchProcessor {
	return NewBatchProcessor(s, maxWorkers)
}

func (s *pdfService) Pipeline() *Pipeline {
	return NewPipeline(s)
}
