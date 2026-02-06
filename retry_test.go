package pdfsdk_test

import (
	"context"
	"testing"
	"time"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func TestRetryConfig(t *testing.T) {
	cfg := pdfsdk.DefaultRetryConfig()

	if cfg.MaxAttempts != 3 {
		t.Errorf("Expected 3 max attempts, got %d", cfg.MaxAttempts)
	}
	if cfg.InitialDelay != 100*time.Millisecond {
		t.Errorf("Expected 100ms initial delay, got %v", cfg.InitialDelay)
	}
	if cfg.MaxDelay != 5*time.Second {
		t.Errorf("Expected 5s max delay, got %v", cfg.MaxDelay)
	}
	if cfg.Multiplier != 2.0 {
		t.Errorf("Expected 2.0 multiplier, got %f", cfg.Multiplier)
	}
	if !cfg.Jitter {
		t.Error("Expected jitter to be enabled")
	}
}

func TestRetrySuccess(t *testing.T) {
	ctx := context.Background()
	cfg := pdfsdk.DefaultRetryConfig()

	attempts := 0
	result, err := pdfsdk.Retry(ctx, cfg, func() (string, error) {
		attempts++
		return "success", nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %s", result)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryWithFailures(t *testing.T) {
	ctx := context.Background()
	cfg := &pdfsdk.RetryConfig{
		MaxAttempts:     3,
		InitialDelay:    10 * time.Millisecond,
		MaxDelay:        50 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          false,
		RetryableErrors: []error{pdfsdk.ErrTimeout},
	}

	attempts := 0
	result, err := pdfsdk.Retry(ctx, cfg, func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", pdfsdk.ErrTimeout // Retryable error
		}
		return "success", nil
	})

	if err != nil {
		t.Errorf("Expected success after retries, got %v", err)
	}
	if result != "success" {
		t.Errorf("Expected 'success', got %s", result)
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryNonRetryableError(t *testing.T) {
	ctx := context.Background()
	cfg := &pdfsdk.RetryConfig{
		MaxAttempts:     3,
		InitialDelay:    10 * time.Millisecond,
		RetryableErrors: []error{pdfsdk.ErrTimeout}, // Only timeout is retryable
	}

	attempts := 0
	_, err := pdfsdk.Retry(ctx, cfg, func() (string, error) {
		attempts++
		return "", pdfsdk.ErrInvalidPDF // Non-retryable
	})

	if err == nil {
		t.Error("Expected error for non-retryable")
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt (no retry for non-retryable), got %d", attempts)
	}
}

func TestRetryContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := &pdfsdk.RetryConfig{
		MaxAttempts:  5,
		InitialDelay: 100 * time.Millisecond,
	}

	attempts := 0
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := pdfsdk.Retry(ctx, cfg, func() (string, error) {
		attempts++
		return "", pdfsdk.ErrTimeout
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestRetryMaxAttempts(t *testing.T) {
	ctx := context.Background()
	cfg := &pdfsdk.RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 10 * time.Millisecond,
	}

	attempts := 0
	_, err := pdfsdk.Retry(ctx, cfg, func() (string, error) {
		attempts++
		return "", pdfsdk.ErrTimeout
	})

	if err == nil {
		t.Error("Expected error after max attempts")
	}
	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

func TestRetryVoid(t *testing.T) {
	ctx := context.Background()
	cfg := pdfsdk.DefaultRetryConfig()

	attempts := 0
	err := pdfsdk.RetryVoid(ctx, cfg, func() error {
		attempts++
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

func TestRetryWrapper(t *testing.T) {
	sdk := pdfsdk.New("http://localhost:3000")
	defer sdk.Close()

	wrapper := sdk.WithRetry(nil) // Use default config
	if wrapper == nil {
		t.Fatal("Expected retry wrapper, got nil")
	}
}
