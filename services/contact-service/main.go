package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/contact-service/api/handlers"
	"github.com/cardoza1991/church-management-system/services/contact-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/contact-service/config"
	"github.com/cardoza1991/church-management-system/services/contact-service/internal/db"
	"github.com/cardoza1991/church-management-system/services/contact-service/internal/models"
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
	
	// Ensure necessary tables exist
	if err := db.EnsureTablesExist(database); err != nil {
		log.Fatalf("Failed to create database tables: %v", err)
	}
	
	// Create repositories
	contactRepo := models.NewContactRepository(database)
	statusRepo := models.NewStatusRepository(database)
	
	// Create handlers
	contactHandler := &handlers.ContactHandler{
		ContactRepo: contactRepo,
		StatusRepo:  statusRepo,
	}
	statusHandler := &handlers.StatusHandler{
		StatusRepo: statusRepo,
	}
	
	// Create router
	r := mux.NewRouter()
	
	// Health check endpoint
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Public endpoints (no authentication required)
	r.HandleFunc("/statuses", statusHandler.GetAllStatuses).Methods("GET")
	r.HandleFunc("/statuses/{id:[0-9]+}", statusHandler.GetStatus).Methods("GET")
	
	// Protected endpoints (require authentication)
	apiRouter := r.PathPrefix("").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)
	
	// Contacts endpoints
	apiRouter.HandleFunc("/contacts", contactHandler.ListContacts).Methods("GET")
	apiRouter.HandleFunc("/contacts", contactHandler.CreateContact).Methods("POST")
	apiRouter.HandleFunc("/contacts/{id:[0-9]+}", contactHandler.GetContact).Methods("GET")
	apiRouter.HandleFunc("/contacts/{id:[0-9]+}", contactHandler.UpdateContact).Methods("PUT")
	apiRouter.HandleFunc("/contacts/{id:[0-9]+}", contactHandler.DeleteContact).Methods("DELETE")
	apiRouter.HandleFunc("/contacts/{id:[0-9]+}/status", contactHandler.UpdateContactStatus).Methods("PUT")
	apiRouter.HandleFunc("/contacts/{id:[0-9]+}/status-history", contactHandler.GetContactStatusHistory).Methods("GET")
	
	// Admin-only endpoints
	adminRouter := r.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware, middleware.AdminRequired)
	adminRouter.HandleFunc("/statuses", statusHandler.CreateStatus).Methods("POST")
	adminRouter.HandleFunc("/statuses/{id:[0-9]+}", statusHandler.UpdateStatus).Methods("PUT")
	adminRouter.HandleFunc("/statuses/{id:[0-9]+}", statusHandler.DeleteStatus).Methods("DELETE")
	
	// Start server
	log.Printf("Contact service starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}