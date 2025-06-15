package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "status_code"},
	)

	// Database Metrics
	databaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"operation", "table", "status"},
	)

	databaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "database_operation_duration_seconds",
			Help: "Duration of database operations in seconds",
			Buckets: []float64{.0001, .0005, .001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"operation", "table", "status"},
	)

	cacheOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "Total number of cache operations",
		},
		[]string{"operation", "status"},
	)

	cacheHitRatio = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_hit_ratio",
			Help: "Cache hit ratio",
		},
		[]string{"type"},
	)

	booksTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "books_total",
			Help: "Total number of books in the system",
		},
	)

	usersTotal = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "users_total",
			Help: "Total number of users in the system",
		},
	)

	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	authAttemptsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_attempts_total",
			Help: "Total number of authentication attempts",
		},
		[]string{"type", "status"},
	)

	errorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "errors_total",
			Help: "Total number of errors",
		},
		[]string{"type", "component"},
	)

	bookOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "book_operations_total",
			Help: "Total number of book operations",
		},
		[]string{"operation", "status"},
	)

	goroutinesActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_active",
			Help: "Number of active goroutines",
		},
	)
)

var (
	cacheHits   int64
	cacheMisses int64
)

// RecordHTTPRequest records an HTTP request metric
func RecordHTTPRequest(method, endpoint, statusCode string, duration time.Duration) {
	httpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint, statusCode).Observe(duration.Seconds())
}

// RecordDatabaseQuery records a database operation metric
func RecordDatabaseQuery(operation, table, status string, duration time.Duration) {
	databaseOperationsTotal.WithLabelValues(operation, table, status).Inc()
	databaseOperationDuration.WithLabelValues(operation, table, status).Observe(duration.Seconds())
}

// RecordCacheOperation records a cache operation metric
func RecordCacheOperation(operation, status string) {
	cacheOperationsTotal.WithLabelValues(operation, status).Inc()

	// Update hit ratio for get operations
	if operation == "get" {
		if status == "hit" {
			cacheHits++
		} else if status == "miss" {
			cacheMisses++
		}

		// Calculate and update hit ratio
		total := cacheHits + cacheMisses
		if total > 0 {
			ratio := float64(cacheHits) / float64(total)
			cacheHitRatio.WithLabelValues("overall").Set(ratio)
		}
	}
}

// RecordAuthAttempt records an authentication attempt
func RecordAuthAttempt(authType, status string) {
	authAttemptsTotal.WithLabelValues(authType, status).Inc()
}

// RecordError records an error occurrence
func RecordError(errorType, component string) {
	errorsTotal.WithLabelValues(errorType, component).Inc()
}

// RecordBookOperation records a book-related operation
func RecordBookOperation(operation, status string) {
	bookOperationsTotal.WithLabelValues(operation, status).Inc()
}

// SetBooksTotal sets the total number of books
func SetBooksTotal(count float64) {
	booksTotal.Set(count)
}

// SetUsersTotal sets the total number of users
func SetUsersTotal(count float64) {
	usersTotal.Set(count)
}

// SetActiveConnections sets the number of active connections
func SetActiveConnections(count float64) {
	activeConnections.Set(count)
}

// SetActiveGoroutines sets the number of active goroutines
func SetActiveGoroutines(count float64) {
	goroutinesActive.Set(count)
}

// GetMetricsRegistry returns the Prometheus registry for custom metrics
func GetMetricsRegistry() *prometheus.Registry {
	return prometheus.DefaultRegisterer.(*prometheus.Registry)
}

// MetricsCollector provides methods to collect application metrics
type MetricsCollector struct {
	startTime time.Time
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		startTime: time.Now(),
	}
}

// CollectSystemMetrics collects system-level metrics
func (mc *MetricsCollector) CollectSystemMetrics() {
	// This could be expanded to collect more detailed system metrics
	// For now, we'll focus on the metrics we've already defined
}

// GetUptime returns the application uptime
func (mc *MetricsCollector) GetUptime() time.Duration {
	return time.Since(mc.startTime)
}

// CacheMetrics represents cache-specific metrics
type CacheMetrics struct {
	Hits   int64   `json:"hits"`
	Misses int64   `json:"misses"`
	Ratio  float64 `json:"hit_ratio"`
	Total  int64   `json:"total_operations"`
}

// GetCacheMetrics returns current cache metrics
func GetCacheMetrics() *CacheMetrics {
	total := cacheHits + cacheMisses
	var ratio float64
	if total > 0 {
		ratio = float64(cacheHits) / float64(total)
	}

	return &CacheMetrics{
		Hits:   cacheHits,
		Misses: cacheMisses,
		Ratio:  ratio,
		Total:  total,
	}
}

// ResetCacheMetrics resets cache hit/miss counters
func ResetCacheMetrics() {
	cacheHits = 0
	cacheMisses = 0
	cacheHitRatio.WithLabelValues("overall").Set(0)
}

// HealthMetrics represents application health metrics
type HealthMetrics struct {
	Uptime    string `json:"uptime"`
	Requests  int64  `json:"total_requests"`
	Errors    int64  `json:"total_errors"`
	Books     int64  `json:"total_books"`
	Users     int64  `json:"total_users"`
	Cache     *CacheMetrics `json:"cache"`
}

// GetHealthMetrics returns comprehensive health metrics
func GetHealthMetrics(collector *MetricsCollector) *HealthMetrics {
	return &HealthMetrics{
		Uptime: collector.GetUptime().String(),
		Cache:  GetCacheMetrics(),
	}
}

// IncrementCounter is a helper function to increment a counter metric
func IncrementCounter(name string, labels map[string]string) error {
	metric := prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: name},
		getLabelKeys(labels),
	)

	counter, err := metric.GetMetricWithLabelValues(getLabelValues(labels)...)
	if err != nil {
		return fmt.Errorf("failed to get metric: %w", err)
	}

	counter.Inc()
	return nil
}

// SetGauge is a helper function to set a gauge metric
func SetGauge(name string, value float64, labels map[string]string) error {
	metric := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{Name: name},
		getLabelKeys(labels),
	)

	gauge, err := metric.GetMetricWithLabelValues(getLabelValues(labels)...)
	if err != nil {
		return fmt.Errorf("failed to get metric: %w", err)
	}

	gauge.Set(value)
	return nil
}

// Helper functions for label handling
func getLabelKeys(labels map[string]string) []string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	return keys
}

func getLabelValues(labels map[string]string) []string {
	values := make([]string, 0, len(labels))
	for _, v := range labels {
		values = append(values, v)
	}
	return values
}

// InitMetrics initializes metrics collection
var metricsInitialized bool

func InitMetrics() {
	if metricsInitialized {
		return
	}
	// Initialize any startup metrics here
	prometheus.MustRegister(prometheus.NewGoCollector())
	prometheus.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	metricsInitialized = true
}

// Global exported variables for testing
var (
	HTTPRequestsTotal        = httpRequestsTotal
	HTTPRequestDuration      = httpRequestDuration
	DatabaseQueryDuration    = databaseOperationDuration
	CacheHits               = cacheOperationsTotal
	CacheMisses             = cacheOperationsTotal
	AuthAttempts            = authAttemptsTotal
	ErrorsTotal             = errorsTotal
	ActiveConnections       = activeConnections
)

// Init initializes the metrics
func Init() {
	InitMetrics()
}
