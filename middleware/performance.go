package middleware

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/handlers"
)

// GzipResponseWriter wraps http.ResponseWriter to provide gzip compression
type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write writes compressed data to the response
func (w *GzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// CompressionMiddleware provides gzip compression for responses
func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Create gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Set headers
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// Wrap response writer
		gzw := &GzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}

// RateLimiter provides in-memory rate limiting
type RateLimiter struct {
	visitors map[string]*Visitor
	mutex    sync.RWMutex
	rate     int
	window   time.Duration
}

// Visitor represents a client with request history
type Visitor struct {
	requests []time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		window:   window,
	}

	// Clean up old entries periodically
	go rl.cleanup()
	return rl
}

// cleanup removes old visitor entries
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		for ip, visitor := range rl.visitors {
			// Remove requests older than the window
			var validRequests []time.Time
			for _, reqTime := range visitor.requests {
				if now.Sub(reqTime) < rl.window {
					validRequests = append(validRequests, reqTime)
				}
			}
			visitor.requests = validRequests

			// Remove visitor if no valid requests
			if len(visitor.requests) == 0 {
				delete(rl.visitors, ip)
			}
		}
		rl.mutex.Unlock()
	}
}

// IsAllowed checks if a request is allowed based on rate limiting
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	visitor, exists := rl.visitors[ip]
	if !exists {
		visitor = &Visitor{requests: []time.Time{}}
		rl.visitors[ip] = visitor
	}

	now := time.Now()
	// Remove requests outside the window
	var validRequests []time.Time
	for _, reqTime := range visitor.requests {
		if now.Sub(reqTime) < rl.window {
			validRequests = append(validRequests, reqTime)
		}
	}
	visitor.requests = validRequests

	// Check if under rate limit
	if len(visitor.requests) >= rl.rate {
		return false
	}

	// Add current request
	visitor.requests = append(visitor.requests, now)
	return true
}

// RateLimitMiddleware provides rate limiting functionality
func RateLimitMiddleware(rate int, window time.Duration) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(rate, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			if !limiter.IsAllowed(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rate))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.WriteHeader(http.StatusTooManyRequests)

				response := map[string]interface{}{
					"error":   "Rate limit exceeded",
					"message": "Too many requests, please try again later",
				}
				json.NewEncoder(w).Encode(response)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check for forwarded headers
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

// TimeoutMiddleware adds timeout to requests
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			// Create a channel to signal completion
			done := make(chan struct{})
			go func() {
				next.ServeHTTP(w, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				// Request completed successfully
			case <-ctx.Done():
				// Request timed out
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestTimeout)
				response := map[string]interface{}{
					"error":   "Request timeout",
					"message": "The request took too long to process",
				}
				json.NewEncoder(w).Encode(response)
			}
		})
	}
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware provides request logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return handlers.LoggingHandler(io.Discard, next)
}
