package routes

import (
	"user-crud/handlers"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	userHandler := handlers.NewUserHandler()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// User routes
	api.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	api.HandleFunc("/users", userHandler.GetAllUsers).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")
	api.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Health check route
	router.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")

	return router
}
