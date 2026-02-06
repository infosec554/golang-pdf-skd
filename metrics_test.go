package pdfsdk_test

import (
	"testing"
	"time"

	pdfsdk "github.com/infosec554/convert-pdf-go-sdk"
)

func TestMetrics(t *testing.T) {
	metrics := pdfsdk.NewMetrics()

	// Initial state
	if metrics.SuccessRate() != 100.0 {
		t.Error("Initial success rate should be 100%")
	}

	// Record successful operation
	metrics.RecordOperation("compress", true, time.Millisecond*100, 1000, 500)

	snapshot := metrics.Snapshot()
	if snapshot.TotalOperations != 1 {
		t.Errorf("Expected 1 total operation, got %d", snapshot.TotalOperations)
	}
	if snapshot.SuccessfulOperations != 1 {
		t.Errorf("Expected 1 successful operation, got %d", snapshot.SuccessfulOperations)
	}
	if snapshot.CompressCount != 1 {
		t.Errorf("Expected 1 compress count, got %d", snapshot.CompressCount)
	}

	// Record failed operation
	metrics.RecordOperation("merge", false, time.Millisecond*50, 2000, 0)

	snapshot = metrics.Snapshot()
	if snapshot.TotalOperations != 2 {
		t.Errorf("Expected 2 total operations, got %d", snapshot.TotalOperations)
	}
	if snapshot.FailedOperations != 1 {
		t.Errorf("Expected 1 failed operation, got %d", snapshot.FailedOperations)
	}

	// Check success rate (1/2 = 50%)
	if snapshot.SuccessRate != 50.0 {
		t.Errorf("Expected 50%% success rate, got %.1f%%", snapshot.SuccessRate)
	}
}

func TestMetricsRecordError(t *testing.T) {
	metrics := pdfsdk.NewMetrics()

	metrics.RecordError("ErrInvalidPDF")
	metrics.RecordError("ErrInvalidPDF")
	metrics.RecordError("ErrTimeout")

	errorCounts := metrics.GetErrorCounts()
	if errorCounts["ErrInvalidPDF"] != 2 {
		t.Errorf("Expected 2 ErrInvalidPDF errors, got %d", errorCounts["ErrInvalidPDF"])
	}
	if errorCounts["ErrTimeout"] != 1 {
		t.Errorf("Expected 1 ErrTimeout errors, got %d", errorCounts["ErrTimeout"])
	}
}

func TestMetricsReset(t *testing.T) {
	metrics := pdfsdk.NewMetrics()

	metrics.RecordOperation("compress", true, time.Second, 1000, 500)
	metrics.RecordError("test")

	metrics.Reset()

	snapshot := metrics.Snapshot()
	if snapshot.TotalOperations != 0 {
		t.Error("Operations should be 0 after reset")
	}

	errorCounts := metrics.GetErrorCounts()
	if len(errorCounts) != 0 {
		t.Error("Error counts should be empty after reset")
	}
}

func TestMetricsAverages(t *testing.T) {
	metrics := pdfsdk.NewMetrics()

	metrics.RecordOperation("compress", true, time.Second, 1000, 500)
	metrics.RecordOperation("compress", true, time.Second*2, 2000, 1000)

	avgDuration := metrics.AverageDuration()
	// Average of 1s and 2s = 1.5s
	if avgDuration < time.Millisecond*1400 || avgDuration > time.Millisecond*1600 {
		t.Errorf("Expected ~1.5s average duration, got %v", avgDuration)
	}

	avgInput := metrics.AverageInputSize()
	// Average of 1000 and 2000 = 1500
	if avgInput != 1500 {
		t.Errorf("Expected 1500 average input size, got %d", avgInput)
	}

	avgOutput := metrics.AverageOutputSize()
	// Average of 500 and 1000 = 750
	if avgOutput != 750 {
		t.Errorf("Expected 750 average output size, got %d", avgOutput)
	}
}

func TestMetricsCompressionRatio(t *testing.T) {
	metrics := pdfsdk.NewMetrics()

	metrics.RecordOperation("compress", true, time.Second, 1000, 500)

	ratio := metrics.CompressionRatio()
	// 500/1000 = 0.5
	if ratio != 0.5 {
		t.Errorf("Expected 0.5 compression ratio, got %f", ratio)
	}
}

func TestPrometheusMetrics(t *testing.T) {
	metrics := pdfsdk.NewMetrics()

	metrics.RecordOperation("compress", true, time.Second, 1000, 500)

	prometheusOutput := metrics.PrometheusMetrics()
	if prometheusOutput == "" {
		t.Error("Prometheus output should not be empty")
	}

	// Check for expected sections
	if len(prometheusOutput) < 100 {
		t.Error("Prometheus output seems too short")
	}
}
