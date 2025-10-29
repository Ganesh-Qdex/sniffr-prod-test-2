package monitoring

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Metrics represents application metrics
type Metrics struct {
	RequestCount      int64     `json:"request_count"`
	ResponseTime      float64   `json:"avg_response_time_ms"`
	ErrorCount        int64     `json:"error_count"`
	ActiveConnections int64     `json:"active_connections"`
	LastUpdated       time.Time `json:"last_updated"`
}

// MetricsCollector collects and stores application metrics
type MetricsCollector struct {
	metrics      Metrics
	mutex        sync.RWMutex
	startTime    time.Time
	requestTimes []time.Duration
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		startTime:    time.Now(),
		requestTimes: make([]time.Duration, 0, 1000), // Keep last 1000 request times
	}
}

// RecordRequest records a request with its response time
func (mc *MetricsCollector) RecordRequest(responseTime time.Duration, isError bool) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.metrics.RequestCount++

	// Add response time to sliding window
	mc.requestTimes = append(mc.requestTimes, responseTime)
	if len(mc.requestTimes) > 1000 {
		mc.requestTimes = mc.requestTimes[1:] // Remove oldest
	}

	// Calculate average response time
	var total time.Duration
	for _, rt := range mc.requestTimes {
		total += rt
	}
	mc.metrics.ResponseTime = float64(total) / float64(len(mc.requestTimes)) / float64(time.Millisecond)

	if isError {
		mc.metrics.ErrorCount++
	}

	mc.metrics.LastUpdated = time.Now()
}

// GetMetrics returns current metrics
func (mc *MetricsCollector) GetMetrics() Metrics {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return mc.metrics
}

// GetUptime returns application uptime
func (mc *MetricsCollector) GetUptime() time.Duration {
	return time.Since(mc.startTime)
}

// MetricsMiddleware records request metrics
func MetricsMiddleware(collector *MetricsCollector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			// Record metrics
			responseTime := time.Since(start)
			isError := wrapped.statusCode >= 400
			collector.RecordRequest(responseTime, isError)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// HealthCheckHandler provides detailed health information
func HealthCheckHandler(collector *MetricsCollector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := collector.GetMetrics()
		uptime := collector.GetUptime()

		health := map[string]interface{}{
			"status":         "healthy",
			"timestamp":      time.Now().UTC(),
			"uptime_seconds": uptime.Seconds(),
			"metrics":        metrics,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(health)
	}
}

// MetricsHandler provides metrics endpoint
func MetricsHandler(collector *MetricsCollector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics := collector.GetMetrics()
		uptime := collector.GetUptime()

		response := map[string]interface{}{
			"metrics":        metrics,
			"uptime_seconds": uptime.Seconds(),
			"timestamp":      time.Now().UTC(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// SetupMonitoringRoutes adds monitoring routes to the router
func SetupMonitoringRoutes(router *mux.Router, collector *MetricsCollector) {
	router.HandleFunc("/health", HealthCheckHandler(collector)).Methods("GET")
	router.HandleFunc("/metrics", MetricsHandler(collector)).Methods("GET")
}
