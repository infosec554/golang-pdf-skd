package pdfsdk

import (
	"context"
	"math/rand"
	"time"
)

type RetryConfig struct {
	MaxAttempts     int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	Multiplier      float64
	Jitter          bool
	RetryableErrors []error
}

func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     5 * time.Second,
		Multiplier:   2.0,
		Jitter:       true,
		RetryableErrors: []error{
			ErrGotenbergUnavailable,
			ErrTimeout,
		},
	}
}

func Retry[T any](ctx context.Context, cfg *RetryConfig, fn func() (T, error)) (T, error) {
	if cfg == nil {
		cfg = DefaultRetryConfig()
	}

	var lastErr error
	var zero T
	delay := cfg.InitialDelay

	for attempt := 1; attempt <= cfg.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		default:
		}

		result, err := fn()
		if err == nil {
			return result, nil
		}

		lastErr = err

		if !isRetryable(err, cfg.RetryableErrors) {
			return zero, err
		}

		if attempt == cfg.MaxAttempts {
			break
		}

		sleepDuration := delay
		if cfg.Jitter {
			jitter := time.Duration(rand.Int63n(int64(delay) / 2))
			sleepDuration = delay + jitter
		}

		timer := time.NewTimer(sleepDuration)
		select {
		case <-ctx.Done():
			timer.Stop()
			return zero, ctx.Err()
		case <-timer.C:
		}

		delay = time.Duration(float64(delay) * cfg.Multiplier)
		if delay > cfg.MaxDelay {
			delay = cfg.MaxDelay
		}
	}

	return zero, &PDFError{
		Op:      "retry",
		Err:     lastErr,
		Details: "max attempts exceeded",
	}
}

func RetryVoid(ctx context.Context, cfg *RetryConfig, fn func() error) error {
	_, err := Retry(ctx, cfg, func() (struct{}, error) {
		return struct{}{}, fn()
	})
	return err
}

func isRetryable(err error, retryableErrors []error) bool {
	if len(retryableErrors) == 0 {
		return true
	}

	for _, retryableErr := range retryableErrors {
		if err == retryableErr {
			return true
		}
		if pdfErr, ok := err.(*PDFError); ok {
			if pdfErr.Err == retryableErr {
				return true
			}
		}
	}
	return false
}

type RetryWrapper struct {
	sdk *SDK
	cfg *RetryConfig
}

func (sdk *SDK) WithRetry(cfg *RetryConfig) *RetryWrapper {
	if cfg == nil {
		cfg = DefaultRetryConfig()
	}
	return &RetryWrapper{sdk: sdk, cfg: cfg}
}

func (rw *RetryWrapper) CompressBytes(ctx context.Context, input []byte) ([]byte, error) {
	return Retry(ctx, rw.cfg, func() ([]byte, error) {
		return rw.sdk.Compress().CompressBytes(input)
	})
}

func (rw *RetryWrapper) MergeBytes(ctx context.Context, inputs [][]byte) ([]byte, error) {
	return Retry(ctx, rw.cfg, func() ([]byte, error) {
		return rw.sdk.Merge().MergeBytes(inputs)
	})
}

func (rw *RetryWrapper) RotateBytes(ctx context.Context, input []byte, angle int, pages string) ([]byte, error) {
	return Retry(ctx, rw.cfg, func() ([]byte, error) {
		return rw.sdk.Rotate().RotateBytes(input, angle, pages)
	})
}
