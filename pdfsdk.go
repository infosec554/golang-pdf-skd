package pdfsdk

import (
	"github.com/infosec554/golang-pdf-sdk/pkg/gotenberg"
	"github.com/infosec554/golang-pdf-sdk/pkg/logger"
	"github.com/infosec554/golang-pdf-sdk/service"
)

const Version = "1.0.0"

func New(gotenbergURL string) service.PDFService {
	return service.NewWithGotenberg(gotenbergURL)
}

func NewWithLogger(gotenbergURL string, log logger.ILogger) service.PDFService {
	gotClient := gotenberg.New(gotenbergURL)
	return service.New(log, gotClient)
}
