package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"user-crud/cache"
	"user-crud/database"
	"user-crud/middleware"
	"user-crud/routes"
)

func main() {
	// Get connection details from environment variables
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "userdb")
	port := getEnv("PORT", "8080")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")
	redisDB := getEnv("REDIS_DB", "0")

	// Connect to MongoDB
	err := database.ConnectDB(mongoURI, dbName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer database.DisconnectDB()

	// Connect to Redis
	redisDBInt := 0
	if db, err := strconv.Atoi(redisDB); err == nil {
		redisDBInt = db
	}
	err = cache.ConnectRedis(redisAddr, redisPassword, redisDBInt)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		log.Println("Continuing without Redis caching...")
	} else {
		defer cache.DisconnectRedis()
	}

	// Setup routes with middleware
	router := routes.SetupRoutes()

	// Apply performance middleware
	handler := middleware.SecurityHeadersMiddleware(router)
	handler = middleware.CORSMiddleware(handler)
	handler = middleware.CompressionMiddleware(handler)
	handler = middleware.RateLimitMiddleware(100, time.Minute)(handler) // 100 requests per minute
	handler = middleware.TimeoutMiddleware(30 * time.Second)(handler)
	handler = middleware.LoggingMiddleware(handler)

	// Create server with optimized settings
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		log.Printf("Health check available at: http://localhost:%s/health", port)
		log.Printf("API endpoints available at: http://localhost:%s/api/v1", port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
