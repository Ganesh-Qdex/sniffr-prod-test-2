package routes

import (
	"user-crud/handlers"
	"user-crud/monitoring"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	userHandler := handlers.NewUserHandler()

	// Create metrics collector
	metricsCollector := monitoring.NewMetricsCollector()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Monitoring routes
	monitoring.SetupMonitoringRoutes(router, metricsCollector)

	return router
}
