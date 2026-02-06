//go:build examples
// +build examples

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

// This file contains complete examples of using the PDF SDK.
// Run with: go run -tags=examples examples/complete_example.go

func main() {
	fmt.Println("🚀 PDF SDK Complete Examples")
	fmt.Println("============================")

	// 1. Basic Initialization
	basicExample()

	// 2. With Options
	optionsExample()

	// 3. Rate Limiting
	rateLimitExample()

	// 4. Retry Mechanism
	retryExample()

	// 5. Metrics
	metricsExample()

	// 6. Pipeline
	pipelineExample()

	// 7. Batch Processing
	batchExample()

	fmt.Println("\n✅ All examples completed!")
}

func basicExample() {
	fmt.Println("\n📦 1. Basic Initialization:")

	sdk := pdfsdk.New("http://localhost:3000")
	defer sdk.Close()

	fmt.Printf("   Version: %s\n", pdfsdk.Version)
	fmt.Printf("   Workers: %d\n", sdk.Stats().MaxWorkers)
	fmt.Println("   ✅ SDK initialized")
}

func optionsExample() {
	fmt.Println("\n⚙️  2. With Options:")

	sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
		GotenbergURL:        "http://localhost:3000",
		MaxWorkers:          20,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		RequestTimeout:      10 * time.Minute,
		ServiceName:         "my-pdf-service",
	})
	defer sdk.Close()

	stats := sdk.Stats()
	fmt.Printf("   Max Workers: %d\n", stats.MaxWorkers)
	fmt.Printf("   Available: %d\n", stats.AvailableWorkers)
	fmt.Println("   ✅ Custom options applied")
}

func rateLimitExample() {
	fmt.Println("\n🚦 3. Rate Limiting:")

	sdk := pdfsdk.New("http://localhost:3000")
	defer sdk.Close()

	// 100 operations per second
	rls := sdk.WithRateLimiter(100, time.Second)
	defer rls.Close()

	fmt.Printf("   Rate: 100 ops/second\n")
	fmt.Printf("   Available tokens: %d\n", rls.RateLimiter().Available())

	// Use rate limiter
	for i := 0; i < 5; i++ {
		rls.RateLimiter().Acquire()
		fmt.Printf("   Acquired token %d\n", i+1)
	}

	fmt.Println("   ✅ Rate limiting works")
}

func retryExample() {
	fmt.Println("\n🔄 4. Retry Mechanism:")

	cfg := &pdfsdk.RetryConfig{
		MaxAttempts:  3,
		InitialDelay: 100 * time.Millisecond,
		MaxDelay:     time.Second,
		Multiplier:   2.0,
		Jitter:       true,
	}

	fmt.Printf("   Max Attempts: %d\n", cfg.MaxAttempts)
	fmt.Printf("   Initial Delay: %v\n", cfg.InitialDelay)
	fmt.Printf("   Backoff Multiplier: %.1f\n", cfg.Multiplier)

	// Example retry with success
	attempts := 0
	result, err := pdfsdk.Retry(context.Background(), cfg, func() (string, error) {
		attempts++
		if attempts < 2 {
			return "", fmt.Errorf("temporary error")
		}
		return "success", nil
	})

	if err != nil {
		fmt.Printf("   ❌ Error: %v\n", err)
	} else {
		fmt.Printf("   Result: %s (after %d attempts)\n", result, attempts)
	}
	fmt.Println("   ✅ Retry mechanism works")
}

func metricsExample() {
	fmt.Println("\n📊 5. Metrics:")

	metrics := pdfsdk.NewMetrics()

	// Record some operations
	metrics.RecordOperation("compress", true, 100*time.Millisecond, 10000, 5000)
	metrics.RecordOperation("compress", true, 150*time.Millisecond, 20000, 8000)
	metrics.RecordOperation("merge", true, 200*time.Millisecond, 30000, 30000)
	metrics.RecordOperation("rotate", false, 50*time.Millisecond, 5000, 0)

	snapshot := metrics.Snapshot()
	fmt.Printf("   Total Operations: %d\n", snapshot.TotalOperations)
	fmt.Printf("   Success Rate: %.1f%%\n", snapshot.SuccessRate)
	fmt.Printf("   Average Duration: %v\n", snapshot.AverageDuration)
	fmt.Printf("   Compression Ratio: %.2f\n", metrics.CompressionRatio())

	// Prometheus format preview
	prom := metrics.PrometheusMetrics()
	fmt.Printf("   Prometheus Output: %d bytes\n", len(prom))
	fmt.Println("   ✅ Metrics collection works")
}

func pipelineExample() {
	fmt.Println("\n⛓️  6. Pipeline:")

	sdk := pdfsdk.New("http://localhost:3000")
	defer sdk.Close()

	// Create pipeline (won't execute without real PDF)
	pipeline := sdk.Pipeline().
		Compress().
		Rotate(90, "all").
		Watermark("CONFIDENTIAL", nil).
		Protect("password123")

	fmt.Println("   Pipeline created with:")
	fmt.Println("   - Compress")
	fmt.Println("   - Rotate 90°")
	fmt.Println("   - Watermark 'CONFIDENTIAL'")
	fmt.Println("   - Protect with password")
	fmt.Println("   ✅ Pipeline ready (use .Execute(pdfBytes) to run)")

	_ = pipeline // prevent unused warning
}

func batchExample() {
	fmt.Println("\n🔄 7. Batch Processing:")

	sdk := pdfsdk.New("http://localhost:3000")
	defer sdk.Close()

	batch := sdk.Batch(5)
	fmt.Println("   Batch processor created with 5 workers")
	fmt.Println("   Available methods:")
	fmt.Println("   - CompressBatch(ctx, [][]byte)")
	fmt.Println("   - RotateBatch(ctx, [][]byte, angle, pages)")
	fmt.Println("   - WatermarkBatch(ctx, [][]byte, text, opts)")
	fmt.Println("   - ProtectBatch(ctx, [][]byte, password)")
	fmt.Println("   ✅ Batch processor ready")

	_ = batch // prevent unused warning
}

// PDF Operations Example (requires real PDF file)
func pdfOperationsExample() {
	sdk := pdfsdk.New("http://localhost:3000")
	defer sdk.Close()

	// Read PDF
	input, err := os.ReadFile("test.pdf")
	if err != nil {
		fmt.Println("No test.pdf found, skipping PDF operations")
		return
	}

	// Get PDF info
	info, _ := sdk.Info().GetInfoBytes(input)
	fmt.Printf("Pages: %d, Encrypted: %v\n", info.PageCount, info.Encrypted)

	// Compress
	compressed, _ := sdk.Compress().CompressBytes(input)
	fmt.Printf("Compressed: %d -> %d bytes\n", len(input), len(compressed))

	// Extract pages
	pages, _ := sdk.Pages().ExtractPages(input, "1-3")
	fmt.Printf("Extracted pages: %d bytes\n", len(pages))

	// Extract text
	text, _ := sdk.Text().ExtractText(input)
	fmt.Printf("Extracted text: %d chars\n", len(text))

	// Extract images
	images, _ := sdk.Images().ExtractImages(input)
	fmt.Printf("Extracted images: %d\n", len(images))
}
