package pdfsdk_test

import (
	"context"
	"testing"
	"time"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

// Benchmark WorkerPool
func BenchmarkWorkerPoolAcquireRelease(b *testing.B) {
	sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
		MaxWorkers: 100,
	})
	defer sdk.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sdk.Workers().Acquire()
			sdk.Workers().Release()
		}
	})
}

func BenchmarkWorkerPoolTryAcquire(b *testing.B) {
	sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
		MaxWorkers: 100,
	})
	defer sdk.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if sdk.Workers().TryAcquire() {
			sdk.Workers().Release()
		}
	}
}

// Benchmark RateLimiter
func BenchmarkRateLimiterAcquire(b *testing.B) {
	limiter := pdfsdk.NewRateLimiter(10000, time.Second)
	defer limiter.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if limiter.TryAcquire() {
			// token acquired
		}
	}
}

// Benchmark Metrics
func BenchmarkMetricsRecordOperation(b *testing.B) {
	metrics := pdfsdk.NewMetrics()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.RecordOperation("compress", true, time.Millisecond, 1000, 500)
		}
	})
}

func BenchmarkMetricsSnapshot(b *testing.B) {
	metrics := pdfsdk.NewMetrics()

	// Pre-populate with some data
	for i := 0; i < 100; i++ {
		metrics.RecordOperation("compress", true, time.Millisecond, 1000, 500)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = metrics.Snapshot()
	}
}

func BenchmarkMetricsPrometheus(b *testing.B) {
	metrics := pdfsdk.NewMetrics()

	// Pre-populate
	for i := 0; i < 100; i++ {
		metrics.RecordOperation("compress", true, time.Millisecond, 1000, 500)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = metrics.PrometheusMetrics()
	}
}

// Benchmark SDK initialization
func BenchmarkNewSDK(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sdk := pdfsdk.New("http://localhost:3000")
		sdk.Close()
	}
}

func BenchmarkNewSDKWithOptions(b *testing.B) {
	opts := &pdfsdk.Options{
		GotenbergURL: "http://localhost:3000",
		MaxWorkers:   10,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sdk := pdfsdk.NewWithOptions(opts)
		sdk.Close()
	}
}

// Benchmark Retry
func BenchmarkRetrySuccess(b *testing.B) {
	cfg := pdfsdk.DefaultRetryConfig()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = pdfsdk.Retry(ctx, cfg, func() (int, error) {
			return 42, nil
		})
	}
}

// Benchmark Error creation
func BenchmarkNewError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = pdfsdk.NewError("compress", pdfsdk.ErrInvalidPDF)
	}
}

func BenchmarkWrapError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = pdfsdk.WrapError("compress", "test.pdf", pdfsdk.ErrInvalidPDF)
	}
}
