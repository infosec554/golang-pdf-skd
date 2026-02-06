package pdfsdk

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidPDF           = errors.New("invalid PDF format")
	ErrEncryptedPDF         = errors.New("PDF is encrypted")
	ErrWrongPassword        = errors.New("incorrect password")
	ErrEmptyInput           = errors.New("empty input")
	ErrPageOutOfRange       = errors.New("page number out of range")
	ErrGotenbergUnavailable = errors.New("Gotenberg server unavailable")
	ErrTimeout              = errors.New("operation timed out")
	ErrWorkerPoolFull       = errors.New("worker pool is full")
	ErrOperationCanceled    = errors.New("operation canceled")
)

type PDFError struct {
	Op      string
	Input   string
	Err     error
	Details string
}

func (e *PDFError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s): %v", e.Op, e.Input, e.Details, e.Err)
	}
	if e.Input != "" {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Input, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *PDFError) Unwrap() error {
	return e.Err
}

func NewError(op string, err error) *PDFError {
	return &PDFError{Op: op, Err: err}
}

func WrapError(op, input string, err error) *PDFError {
	return &PDFError{Op: op, Input: input, Err: err}
}

func IsInvalidPDF(err error) bool {
	return errors.Is(err, ErrInvalidPDF)
}

func IsEncrypted(err error) bool {
	return errors.Is(err, ErrEncryptedPDF)
}

func IsTimeout(err error) bool {
	return errors.Is(err, ErrTimeout)
}

func IsGotenbergUnavailable(err error) bool {
	return errors.Is(err, ErrGotenbergUnavailable)
}
