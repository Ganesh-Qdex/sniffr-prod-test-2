package routes

import (
	"user-crud/handlers"
	"user-crud/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()
	userHandler := handlers.NewUserHandler()
	authHandler := handlers.NewAuthHandler()

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Public authentication routes
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Protected authentication routes
	auth := api.PathPrefix("/auth").Subrouter()
	auth.Use(middleware.AuthMiddleware)
	auth.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")
	auth.HandleFunc("/profile", authHandler.UpdateProfile).Methods("PUT")
	auth.HandleFunc("/change-password", authHandler.ChangePassword).Methods("PUT")

	// Protected user routes
	users := api.PathPrefix("/users").Subrouter()
	users.Use(middleware.AuthMiddleware)
	users.HandleFunc("", userHandler.CreateUser).Methods("POST")
	users.HandleFunc("", userHandler.GetAllUsers).Methods("GET")
	users.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	users.HandleFunc("/{id}", userHandler.UpdateUser).Methods("PUT")
	users.HandleFunc("/{id}", userHandler.DeleteUser).Methods("DELETE")

	// Health check route (public)
	router.HandleFunc("/health", userHandler.HealthCheck).Methods("GET")

	return router
}
