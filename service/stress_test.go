package service_test

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/infosec554/convert-pdf-go-sdk/service"
)

// TestStress is a stress test that runs multiple operations in parallel
// to ensure thread safety and stability under load.
func TestStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	gotURL := os.Getenv("GOTENBERG_URL")
	if gotURL == "" {
		t.Skip("GOTENBERG_URL not set")
	}

	testPDFPath := "testdata/test.pdf"
	if _, err := os.Stat(testPDFPath); os.IsNotExist(err) {
		t.Skip("testdata/test.pdf not found")
	}

	pdfBytes, err := os.ReadFile(testPDFPath)
	if err != nil {
		t.Fatal(err)
	}

	pdfService := service.NewWithGotenberg(gotURL)

	// Create a worker pool to limit concurrency if needed, but here we want to stress it.
	// We'll spawn N goroutines.
	concurrency := 20
	iterations := 5

	var wg sync.WaitGroup
	errCh := make(chan error, concurrency*iterations)

	start := time.Now()

	t.Logf("Starting stress test with %d routines, %d iterations each", concurrency, iterations)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				// Randomly choose an operation
				op := rand.Intn(5)
				var err error

				switch op {
				case 0:
					// Compress
					_, err = pdfService.Compress().CompressBytes(pdfBytes)
				case 1:
					// Rotate
					_, err = pdfService.Rotate().RotateBytes(pdfBytes, 90, "all")
				case 2:
					// Watermark
					_, err = pdfService.Watermark().AddWatermarkBytes(pdfBytes, fmt.Sprintf("Stress-%d-%d", id, j), nil)
				case 3:
					// Metadata
					_, err = pdfService.Metadata().GetMetadata(pdfBytes)
				case 4:
					// Split
					// Note: SplitBytes returns zip bytes
					_, err = pdfService.Split().SplitBytes(pdfBytes, "1")
				}

				if err != nil {
					errCh <- fmt.Errorf("Routine %d Iter %d Op %d failed: %v", id, j, op, err)
				}
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	duration := time.Since(start)
	t.Logf("Stress test completed in %v", duration)

	failCount := 0
	for err := range errCh {
		t.Error(err)
		failCount++
		if failCount >= 10 {
			t.Fatal("Too many errors, aborting stress test eval")
		}
	}
}
