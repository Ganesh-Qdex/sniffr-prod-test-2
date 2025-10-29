package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"user-crud/database"
	"user-crud/routes"
)

func main() {
	// Get MongoDB connection details from environment variables
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	dbName := getEnv("DB_NAME", "userdb")
	port := getEnv("PORT", "8080")

	// Connect to MongoDB
	err := database.ConnectDB(mongoURI, dbName)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer database.DisconnectDB()

	// Setup routes
	router := routes.SetupRoutes()

	// Create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
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
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
