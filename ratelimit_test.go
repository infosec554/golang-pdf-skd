package pdfsdk_test

import (
	"testing"
	"time"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func TestRateLimiter(t *testing.T) {
	limiter := pdfsdk.NewRateLimiter(5, time.Second)
	defer limiter.Stop()

	// Should have 5 tokens initially
	if limiter.Available() != 5 {
		t.Errorf("Expected 5 available tokens, got %d", limiter.Available())
	}

	// Acquire 3 tokens
	for i := 0; i < 3; i++ {
		limiter.Acquire()
	}

	if limiter.Available() != 2 {
		t.Errorf("Expected 2 available tokens, got %d", limiter.Available())
	}
}

func TestRateLimiterTryAcquire(t *testing.T) {
	limiter := pdfsdk.NewRateLimiter(2, time.Second)
	defer limiter.Stop()

	// First two should succeed
	if !limiter.TryAcquire() {
		t.Error("First TryAcquire should succeed")
	}
	if !limiter.TryAcquire() {
		t.Error("Second TryAcquire should succeed")
	}

	// Third should fail
	if limiter.TryAcquire() {
		t.Error("Third TryAcquire should fail")
	}
}

func TestRateLimiterRefill(t *testing.T) {
	limiter := pdfsdk.NewRateLimiter(3, 100*time.Millisecond)
	defer limiter.Stop()

	// Exhaust all tokens
	for i := 0; i < 3; i++ {
		limiter.Acquire()
	}

	if limiter.Available() != 0 {
		t.Errorf("Expected 0 available tokens, got %d", limiter.Available())
	}

	// Wait for refill
	time.Sleep(150 * time.Millisecond)

	// Should have tokens again
	if limiter.Available() == 0 {
		t.Error("Expected tokens to be refilled")
	}
}

func TestRateLimitedSDK(t *testing.T) {
	sdk := pdfsdk.New("http://localhost:3000")
	defer sdk.Close()

	rls := sdk.WithRateLimiter(10, time.Second)
	defer rls.Close()

	if rls.RateLimiter().Available() != 10 {
		t.Errorf("Expected 10 available tokens, got %d", rls.RateLimiter().Available())
	}
}
