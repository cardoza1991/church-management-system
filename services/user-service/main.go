package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/user-service/api/handlers"
	"github.com/cardoza1991/church-management-system/services/user-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/user-service/config"
	"github.com/cardoza1991/church-management-system/services/user-service/internal/db"
	"github.com/cardoza1991/church-management-system/services/user-service/internal/models"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Connect to database
	database, err := db.Connect(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	
	// Create repositories
	userRepo := models.NewUserRepository(database)
	
	// Create handlers
	authHandler := &handlers.AuthHandler{UserRepo: userRepo}
	userHandler := &handlers.UserHandler{UserRepo: userRepo}
	
	// Create router
	r := mux.NewRouter()
	
	// Health check endpoint
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Auth endpoints
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")
	
	// Protected user endpoints
	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.Use(middleware.AuthMiddleware)
	userRouter.HandleFunc("/me", userHandler.GetSelf).Methods("GET")
	userRouter.HandleFunc("/{id}", userHandler.GetUser).Methods("GET")
	
	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
