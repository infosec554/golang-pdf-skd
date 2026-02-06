package pdfsdk

import (
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	TotalOperations      int64
	SuccessfulOperations int64
	FailedOperations     int64
	CompressCount        int64
	MergeCount           int64
	SplitCount           int64
	RotateCount          int64
	WatermarkCount       int64
	ProtectCount         int64
	UnlockCount          int64
	ConvertCount         int64
	InfoCount            int64
	PageOpsCount         int64
	TextExtractCount     int64
	ImageExtractCount    int64
	TotalDuration        time.Duration
	LastOperationTime    time.Time
	TotalInputBytes      int64
	TotalOutputBytes     int64
	errorCounts          map[string]int64
	mu                   sync.RWMutex
}

func NewMetrics() *Metrics {
	return &Metrics{
		errorCounts: make(map[string]int64),
	}
}

func (m *Metrics) RecordOperation(service string, success bool, duration time.Duration, inputSize, outputSize int64) {
	atomic.AddInt64(&m.TotalOperations, 1)
	if success {
		atomic.AddInt64(&m.SuccessfulOperations, 1)
	} else {
		atomic.AddInt64(&m.FailedOperations, 1)
	}

	atomic.AddInt64(&m.TotalInputBytes, inputSize)
	atomic.AddInt64(&m.TotalOutputBytes, outputSize)

	m.mu.Lock()
	m.TotalDuration += duration
	m.LastOperationTime = time.Now()
	m.mu.Unlock()

	switch service {
	case "compress":
		atomic.AddInt64(&m.CompressCount, 1)
	case "merge":
		atomic.AddInt64(&m.MergeCount, 1)
	case "split":
		atomic.AddInt64(&m.SplitCount, 1)
	case "rotate":
		atomic.AddInt64(&m.RotateCount, 1)
	case "watermark":
		atomic.AddInt64(&m.WatermarkCount, 1)
	case "protect":
		atomic.AddInt64(&m.ProtectCount, 1)
	case "unlock":
		atomic.AddInt64(&m.UnlockCount, 1)
	case "convert":
		atomic.AddInt64(&m.ConvertCount, 1)
	case "info":
		atomic.AddInt64(&m.InfoCount, 1)
	case "pages":
		atomic.AddInt64(&m.PageOpsCount, 1)
	case "text":
		atomic.AddInt64(&m.TextExtractCount, 1)
	case "images":
		atomic.AddInt64(&m.ImageExtractCount, 1)
	}
}

func (m *Metrics) RecordError(errType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errorCounts[errType]++
}

func (m *Metrics) GetErrorCounts() map[string]int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]int64)
	for k, v := range m.errorCounts {
		result[k] = v
	}
	return result
}

func (m *Metrics) SuccessRate() float64 {
	total := atomic.LoadInt64(&m.TotalOperations)
	if total == 0 {
		return 100.0
	}
	success := atomic.LoadInt64(&m.SuccessfulOperations)
	return float64(success) / float64(total) * 100
}

func (m *Metrics) AverageDuration() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	total := atomic.LoadInt64(&m.TotalOperations)
	if total == 0 {
		return 0
	}
	return time.Duration(int64(m.TotalDuration) / total)
}

func (m *Metrics) AverageInputSize() int64 {
	total := atomic.LoadInt64(&m.TotalOperations)
	if total == 0 {
		return 0
	}
	return atomic.LoadInt64(&m.TotalInputBytes) / total
}

func (m *Metrics) AverageOutputSize() int64 {
	total := atomic.LoadInt64(&m.TotalOperations)
	if total == 0 {
		return 0
	}
	return atomic.LoadInt64(&m.TotalOutputBytes) / total
}

func (m *Metrics) CompressionRatio() float64 {
	input := atomic.LoadInt64(&m.TotalInputBytes)
	if input == 0 {
		return 1.0
	}
	output := atomic.LoadInt64(&m.TotalOutputBytes)
	return float64(output) / float64(input)
}

func (m *Metrics) Reset() {
	atomic.StoreInt64(&m.TotalOperations, 0)
	atomic.StoreInt64(&m.SuccessfulOperations, 0)
	atomic.StoreInt64(&m.FailedOperations, 0)
	atomic.StoreInt64(&m.CompressCount, 0)
	atomic.StoreInt64(&m.MergeCount, 0)
	atomic.StoreInt64(&m.SplitCount, 0)
	atomic.StoreInt64(&m.RotateCount, 0)
	atomic.StoreInt64(&m.WatermarkCount, 0)
	atomic.StoreInt64(&m.ProtectCount, 0)
	atomic.StoreInt64(&m.UnlockCount, 0)
	atomic.StoreInt64(&m.ConvertCount, 0)
	atomic.StoreInt64(&m.InfoCount, 0)
	atomic.StoreInt64(&m.PageOpsCount, 0)
	atomic.StoreInt64(&m.TextExtractCount, 0)
	atomic.StoreInt64(&m.ImageExtractCount, 0)
	atomic.StoreInt64(&m.TotalInputBytes, 0)
	atomic.StoreInt64(&m.TotalOutputBytes, 0)

	m.mu.Lock()
	m.TotalDuration = 0
	m.LastOperationTime = time.Time{}
	m.errorCounts = make(map[string]int64)
	m.mu.Unlock()
}

func (m *Metrics) Snapshot() MetricsSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return MetricsSnapshot{
		TotalOperations:      atomic.LoadInt64(&m.TotalOperations),
		SuccessfulOperations: atomic.LoadInt64(&m.SuccessfulOperations),
		FailedOperations:     atomic.LoadInt64(&m.FailedOperations),
		CompressCount:        atomic.LoadInt64(&m.CompressCount),
		MergeCount:           atomic.LoadInt64(&m.MergeCount),
		SplitCount:           atomic.LoadInt64(&m.SplitCount),
		RotateCount:          atomic.LoadInt64(&m.RotateCount),
		WatermarkCount:       atomic.LoadInt64(&m.WatermarkCount),
		ProtectCount:         atomic.LoadInt64(&m.ProtectCount),
		UnlockCount:          atomic.LoadInt64(&m.UnlockCount),
		ConvertCount:         atomic.LoadInt64(&m.ConvertCount),
		TotalInputBytes:      atomic.LoadInt64(&m.TotalInputBytes),
		TotalOutputBytes:     atomic.LoadInt64(&m.TotalOutputBytes),
		TotalDuration:        m.TotalDuration,
		LastOperationTime:    m.LastOperationTime,
		SuccessRate:          m.SuccessRate(),
		AverageDuration:      m.AverageDuration(),
	}
}

type MetricsSnapshot struct {
	TotalOperations      int64
	SuccessfulOperations int64
	FailedOperations     int64
	CompressCount        int64
	MergeCount           int64
	SplitCount           int64
	RotateCount          int64
	WatermarkCount       int64
	ProtectCount         int64
	UnlockCount          int64
	ConvertCount         int64
	TotalInputBytes      int64
	TotalOutputBytes     int64
	TotalDuration        time.Duration
	LastOperationTime    time.Time
	SuccessRate          float64
	AverageDuration      time.Duration
}

func (m *Metrics) PrometheusMetrics() string {
	snapshot := m.Snapshot()

	return `# HELP pdfsdk_operations_total Total number of PDF operations
# TYPE pdfsdk_operations_total counter
pdfsdk_operations_total{status="success"} ` + formatInt64(snapshot.SuccessfulOperations) + `
pdfsdk_operations_total{status="failed"} ` + formatInt64(snapshot.FailedOperations) + `

# HELP pdfsdk_operations_by_service Operations count by service
# TYPE pdfsdk_operations_by_service counter
pdfsdk_operations_by_service{service="compress"} ` + formatInt64(snapshot.CompressCount) + `
pdfsdk_operations_by_service{service="merge"} ` + formatInt64(snapshot.MergeCount) + `
pdfsdk_operations_by_service{service="split"} ` + formatInt64(snapshot.SplitCount) + `
pdfsdk_operations_by_service{service="rotate"} ` + formatInt64(snapshot.RotateCount) + `
pdfsdk_operations_by_service{service="watermark"} ` + formatInt64(snapshot.WatermarkCount) + `
pdfsdk_operations_by_service{service="protect"} ` + formatInt64(snapshot.ProtectCount) + `
pdfsdk_operations_by_service{service="unlock"} ` + formatInt64(snapshot.UnlockCount) + `
pdfsdk_operations_by_service{service="convert"} ` + formatInt64(snapshot.ConvertCount) + `

# HELP pdfsdk_bytes_processed_total Total bytes processed
# TYPE pdfsdk_bytes_processed_total counter
pdfsdk_bytes_processed_total{direction="input"} ` + formatInt64(snapshot.TotalInputBytes) + `
pdfsdk_bytes_processed_total{direction="output"} ` + formatInt64(snapshot.TotalOutputBytes) + `

# HELP pdfsdk_success_rate Success rate percentage
# TYPE pdfsdk_success_rate gauge
pdfsdk_success_rate ` + formatFloat64(snapshot.SuccessRate) + `

# HELP pdfsdk_average_duration_seconds Average operation duration
# TYPE pdfsdk_average_duration_seconds gauge
pdfsdk_average_duration_seconds ` + formatFloat64(snapshot.AverageDuration.Seconds()) + `
`
}

func formatInt64(v int64) string {
	return strconv.FormatInt(v, 10)
}

func formatFloat64(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}
