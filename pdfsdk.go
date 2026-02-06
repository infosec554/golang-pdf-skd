package pdfsdk

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/infosec554/convert-pdf-go-sdk/pkg/gotenberg"
	"github.com/infosec554/convert-pdf-go-sdk/pkg/logger"
	"github.com/infosec554/convert-pdf-go-sdk/service"
)

const Version = "2.1.0"

type Options struct {
	GotenbergURL        string
	LogLevel            string
	ServiceName         string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int
	IdleConnTimeout     time.Duration
	RequestTimeout      time.Duration
	MaxWorkers          int
}

func DefaultOptions() *Options {
	return &Options{
		GotenbergURL:        "http://localhost:3000",
		LogLevel:            "info",
		ServiceName:         "golang-pdf-sdk",
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     100,
		IdleConnTimeout:     90 * time.Second,
		RequestTimeout:      5 * time.Minute,
		MaxWorkers:          10,
	}
}

type SDK struct {
	service.PDFService
	httpClient *http.Client
	workerPool *WorkerPool
	opts       *Options
}

type WorkerPool struct {
	maxWorkers int
	semaphore  chan struct{}
	mu         sync.RWMutex
	active     int
	processed  int64
}

func NewWorkerPool(maxWorkers int) *WorkerPool {
	if maxWorkers <= 0 {
		maxWorkers = 10
	}
	return &WorkerPool{
		maxWorkers: maxWorkers,
		semaphore:  make(chan struct{}, maxWorkers),
	}
}

func (wp *WorkerPool) Acquire() {
	wp.semaphore <- struct{}{}
	wp.mu.Lock()
	wp.active++
	wp.mu.Unlock()
}

func (wp *WorkerPool) Release() {
	wp.mu.Lock()
	wp.active--
	wp.processed++
	wp.mu.Unlock()
	<-wp.semaphore
}

func (wp *WorkerPool) TryAcquire() bool {
	select {
	case wp.semaphore <- struct{}{}:
		wp.mu.Lock()
		wp.active++
		wp.mu.Unlock()
		return true
	default:
		return false
	}
}

func (wp *WorkerPool) Stats() (active, maxWorkers int, processed int64) {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.active, wp.maxWorkers, wp.processed
}

func (wp *WorkerPool) Available() int {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.maxWorkers - wp.active
}

func New(gotenbergURL string) *SDK {
	opts := DefaultOptions()
	opts.GotenbergURL = gotenbergURL
	return NewWithOptions(opts)
}

func NewWithOptions(opts *Options) *SDK {
	if opts == nil {
		opts = DefaultOptions()
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        opts.MaxIdleConns,
		MaxIdleConnsPerHost: opts.MaxIdleConnsPerHost,
		MaxConnsPerHost:     opts.MaxConnsPerHost,
		IdleConnTimeout:     opts.IdleConnTimeout,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableCompression:  false,
		ForceAttemptHTTP2:   true,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   opts.RequestTimeout,
	}

	log := logger.New(opts.ServiceName)
	gotClient := gotenberg.NewWithClient(opts.GotenbergURL, httpClient)

	return &SDK{
		PDFService: service.New(log, gotClient),
		httpClient: httpClient,
		workerPool: NewWorkerPool(opts.MaxWorkers),
		opts:       opts,
	}
}

func NewWithLogger(gotenbergURL string, log logger.ILogger) *SDK {
	opts := DefaultOptions()
	opts.GotenbergURL = gotenbergURL

	transport := &http.Transport{
		MaxIdleConns:        opts.MaxIdleConns,
		MaxIdleConnsPerHost: opts.MaxIdleConnsPerHost,
		IdleConnTimeout:     opts.IdleConnTimeout,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   opts.RequestTimeout,
	}

	gotClient := gotenberg.NewWithClient(gotenbergURL, httpClient)
	return &SDK{
		PDFService: service.New(log, gotClient),
		httpClient: httpClient,
		workerPool: NewWorkerPool(opts.MaxWorkers),
		opts:       opts,
	}
}

func (sdk *SDK) Workers() *WorkerPool {
	return sdk.workerPool
}

func (sdk *SDK) Stats() SDKStats {
	active, max, processed := sdk.workerPool.Stats()
	return SDKStats{
		ActiveWorkers:    active,
		MaxWorkers:       max,
		ProcessedTasks:   processed,
		AvailableWorkers: max - active,
	}
}

type SDKStats struct {
	ActiveWorkers    int
	MaxWorkers       int
	ProcessedTasks   int64
	AvailableWorkers int
}

func (sdk *SDK) Close() {
	if sdk.httpClient != nil {
		sdk.httpClient.CloseIdleConnections()
	}
}
