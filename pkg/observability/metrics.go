package observability

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel/metric"
)

var (
	// RequestCounter tracks total HTTP requests
	RequestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// ResponseTime tracks request duration
	ResponseTime = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_time_seconds",
			Help:    "Response time in seconds",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "endpoint"},
	)

	// ErrorCounter tracks total errors
	ErrorCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "endpoint", "error_type"},
	)

	// APIKeyUsage tracks API key usage
	APIKeyUsage = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_key_usage_total",
			Help: "Total number of API key usages",
		},
		[]string{"key_id"},
	)

	// Track API key creation/deletion
	APIKeyOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_key_operations_total",
			Help: "Total number of API key operations",
		},
		[]string{"operation", "user_id"},
	)

	// Track authentication attempts
	AuthAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"success", "method"},
	)

	// New OpenTelemetry-aligned metrics
	HTTPServerDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_server_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{.005, .01, .025, .05, .075, .1, .25, .5, .75, 1, 2.5, 5, 7.5, 10},
		},
		[]string{"method", "route", "status_code"},
	)

	HTTPServerRequestSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_server_request_size_bytes",
			Help: "Size of HTTP request in bytes",
		},
		[]string{"method", "path"},
	)

	HTTPServerActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_server_active_requests",
			Help: "Number of in-flight HTTP requests",
		},
		[]string{"method", "route"},
	)

	HTTPServerResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_server_response_size_bytes",
			Help: "Size of HTTP response in bytes",
		},
		[]string{"method", "path"},
	)
)

// RuntimeMetrics holds runtime metric instruments
type RuntimeMetrics struct {
	// Memory metrics
	heapAlloc   metric.Int64UpDownCounter
	heapIdle    metric.Int64UpDownCounter
	heapInuse   metric.Int64UpDownCounter
	heapObjects metric.Int64UpDownCounter

	// Goroutine metrics
	goroutines metric.Int64UpDownCounter

	// GC metrics
	gcCount      metric.Int64Counter
	gcPauseTotal metric.Int64Counter
	gcPauseNs    metric.Float64Histogram
}

// InitRuntimeMetrics initializes and starts runtime metrics collection
func InitRuntimeMetrics(ctx context.Context) error {
	meter := GetMeter()

	rm := &RuntimeMetrics{}
	var err error

	// Initialize memory metrics
	if rm.heapAlloc, err = meter.Int64UpDownCounter(
		"process.runtime.go.mem.heap_alloc",
		metric.WithDescription("Bytes of allocated heap objects"),
		metric.WithUnit("bytes"),
	); err != nil {
		return fmt.Errorf("failed to create heap_alloc metric: %w", err)
	}

	if rm.heapIdle, err = meter.Int64UpDownCounter(
		"process.runtime.go.mem.heap_idle",
		metric.WithDescription("Bytes in idle (unused) spans"),
		metric.WithUnit("bytes"),
	); err != nil {
		return fmt.Errorf("failed to create heap_idle metric: %w", err)
	}

	if rm.heapInuse, err = meter.Int64UpDownCounter(
		"process.runtime.go.mem.heap_inuse",
		metric.WithDescription("Bytes in in-use spans"),
		metric.WithUnit("bytes"),
	); err != nil {
		return fmt.Errorf("failed to create heap_inuse metric: %w", err)
	}

	if rm.heapObjects, err = meter.Int64UpDownCounter(
		"process.runtime.go.mem.heap_objects",
		metric.WithDescription("Number of allocated heap objects"),
	); err != nil {
		return fmt.Errorf("failed to create heap_objects metric: %w", err)
	}

	// Initialize goroutine metrics
	if rm.goroutines, err = meter.Int64UpDownCounter(
		"process.runtime.go.goroutines",
		metric.WithDescription("Number of goroutines that currently exist"),
	); err != nil {
		return fmt.Errorf("failed to create goroutines metric: %w", err)
	}

	// Initialize GC metrics
	if rm.gcCount, err = meter.Int64Counter(
		"process.runtime.go.gc.count",
		metric.WithDescription("Number of completed GC cycles"),
	); err != nil {
		return fmt.Errorf("failed to create gc_count metric: %w", err)
	}

	if rm.gcPauseTotal, err = meter.Int64Counter(
		"process.runtime.go.gc.pause_total",
		metric.WithDescription("Cumulative nanoseconds in GC stop-the-world pauses"),
		metric.WithUnit("ns"),
	); err != nil {
		return fmt.Errorf("failed to create gc_pause_total metric: %w", err)
	}

	if rm.gcPauseNs, err = meter.Float64Histogram(
		"process.runtime.go.gc.pause",
		metric.WithDescription("Amount of time in GC stop-the-world pauses"),
		metric.WithUnit("ns"),
	); err != nil {
		return fmt.Errorf("failed to create gc_pause metric: %w", err)
	}

	// Start collection loop
	go rm.collect(ctx)

	return nil
}

func (rm *RuntimeMetrics) collect(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	var stats runtime.MemStats
	var lastPauseNs uint64
	var lastNumGC uint32

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runtime.ReadMemStats(&stats)

			// Update memory metrics
			rm.heapAlloc.Add(ctx, int64(stats.HeapAlloc))
			rm.heapIdle.Add(ctx, int64(stats.HeapIdle))
			rm.heapInuse.Add(ctx, int64(stats.HeapInuse))
			rm.heapObjects.Add(ctx, int64(stats.HeapObjects))

			// Update goroutine count
			rm.goroutines.Add(ctx, int64(runtime.NumGoroutine()))

			// Update GC metrics
			if stats.NumGC > lastNumGC {
				rm.gcCount.Add(ctx, int64(stats.NumGC-lastNumGC))

				// Record pause times
				for i := lastNumGC; i < stats.NumGC; i++ {
					pauseNs := stats.PauseNs[i%256]
					if pauseNs > lastPauseNs {
						rm.gcPauseTotal.Add(ctx, int64(pauseNs-lastPauseNs))
						rm.gcPauseNs.Record(ctx, float64(pauseNs-lastPauseNs))
						lastPauseNs = pauseNs
					}
				}
				lastNumGC = stats.NumGC
			}
		}
	}
}
