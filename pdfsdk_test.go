package pdfsdk_test

import (
	"testing"
	"time"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

// Test SDK initialization
func TestNew(t *testing.T) {
	sdk := pdfsdk.New("http://localhost:3000")
	if sdk == nil {
		t.Fatal("Expected SDK instance, got nil")
	}
	defer sdk.Close()
}

func TestNewWithOptions(t *testing.T) {
	opts := &pdfsdk.Options{
		GotenbergURL:        "http://localhost:3000",
		MaxWorkers:          5,
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 5,
		RequestTimeout:      time.Minute,
	}

	sdk := pdfsdk.NewWithOptions(opts)
	if sdk == nil {
		t.Fatal("Expected SDK instance, got nil")
	}
	defer sdk.Close()

	stats := sdk.Stats()
	if stats.MaxWorkers != 5 {
		t.Errorf("Expected MaxWorkers=5, got %d", stats.MaxWorkers)
	}
}

func TestNewWithNilOptions(t *testing.T) {
	sdk := pdfsdk.NewWithOptions(nil)
	if sdk == nil {
		t.Fatal("Expected SDK instance with default options, got nil")
	}
	defer sdk.Close()
}

func TestDefaultOptions(t *testing.T) {
	opts := pdfsdk.DefaultOptions()

	if opts.GotenbergURL != "http://localhost:3000" {
		t.Errorf("Expected default GotenbergURL, got %s", opts.GotenbergURL)
	}
	if opts.MaxWorkers != 10 {
		t.Errorf("Expected default MaxWorkers=10, got %d", opts.MaxWorkers)
	}
	if opts.MaxIdleConns != 100 {
		t.Errorf("Expected default MaxIdleConns=100, got %d", opts.MaxIdleConns)
	}
}

func TestVersion(t *testing.T) {
	if pdfsdk.Version == "" {
		t.Error("Expected version string, got empty")
	}
	if pdfsdk.Version != "2.1.0" {
		t.Errorf("Expected version 2.1.0, got %s", pdfsdk.Version)
	}
}

// Test Worker Pool
func TestWorkerPool(t *testing.T) {
	sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
		MaxWorkers: 3,
	})
	defer sdk.Close()

	// Initial state
	if sdk.Workers().Available() != 3 {
		t.Errorf("Expected 3 available workers, got %d", sdk.Workers().Available())
	}

	// Acquire worker
	sdk.Workers().Acquire()
	if sdk.Workers().Available() != 2 {
		t.Errorf("Expected 2 available workers after acquire, got %d", sdk.Workers().Available())
	}

	// Release worker
	sdk.Workers().Release()
	if sdk.Workers().Available() != 3 {
		t.Errorf("Expected 3 available workers after release, got %d", sdk.Workers().Available())
	}
}

func TestWorkerPoolTryAcquire(t *testing.T) {
	sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
		MaxWorkers: 1,
	})
	defer sdk.Close()

	// First acquire should succeed
	if !sdk.Workers().TryAcquire() {
		t.Error("Expected TryAcquire to succeed on first call")
	}

	// Second acquire should fail (only 1 worker)
	if sdk.Workers().TryAcquire() {
		t.Error("Expected TryAcquire to fail when all workers busy")
	}

	// Release and try again
	sdk.Workers().Release()
	if !sdk.Workers().TryAcquire() {
		t.Error("Expected TryAcquire to succeed after release")
	}
	sdk.Workers().Release()
}

func TestWorkerPoolStats(t *testing.T) {
	sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
		MaxWorkers: 5,
	})
	defer sdk.Close()

	active, max, processed := sdk.Workers().Stats()
	if active != 0 {
		t.Errorf("Expected 0 active workers, got %d", active)
	}
	if max != 5 {
		t.Errorf("Expected 5 max workers, got %d", max)
	}
	if processed != 0 {
		t.Errorf("Expected 0 processed, got %d", processed)
	}

	// Acquire and release to increment processed
	sdk.Workers().Acquire()
	sdk.Workers().Release()

	_, _, processed = sdk.Workers().Stats()
	if processed != 1 {
		t.Errorf("Expected 1 processed after release, got %d", processed)
	}
}

// Test SDK Stats
func TestSDKStats(t *testing.T) {
	sdk := pdfsdk.NewWithOptions(&pdfsdk.Options{
		MaxWorkers: 10,
	})
	defer sdk.Close()

	stats := sdk.Stats()
	if stats.MaxWorkers != 10 {
		t.Errorf("Expected 10 max workers, got %d", stats.MaxWorkers)
	}
	if stats.ActiveWorkers != 0 {
		t.Errorf("Expected 0 active workers, got %d", stats.ActiveWorkers)
	}
	if stats.AvailableWorkers != 10 {
		t.Errorf("Expected 10 available workers, got %d", stats.AvailableWorkers)
	}
}
