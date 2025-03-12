package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/cardoza1991/church-management-system/services/reservation-service/api/handlers"
	"github.com/cardoza1991/church-management-system/services/reservation-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/reservation-service/config"
	"github.com/cardoza1991/church-management-system/services/reservation-service/internal/db"
	"github.com/cardoza1991/church-management-system/services/reservation-service/internal/models"
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
	roomRepo := models.NewRoomRepository(database)
	reservationRepo := models.NewReservationRepository(database)
	
	// Create handlers
	roomHandler := &handlers.RoomHandler{
		RoomRepo: roomRepo,
	}
	reservationHandler := &handlers.ReservationHandler{
		ReservationRepo: reservationRepo,
		RoomRepo:        roomRepo,
	}
	
	// Create router
	r := mux.NewRouter()
	
	// Health check endpoint
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Public endpoints (no authentication required)
	r.HandleFunc("/rooms", roomHandler.GetAllRooms).Methods("GET")
	r.HandleFunc("/rooms/{id:[0-9]+}", roomHandler.GetRoom).Methods("GET")
	r.HandleFunc("/rooms/available", roomHandler.GetAvailableRooms).Methods("GET")
	
	// Protected endpoints (require authentication)
	apiRouter := r.PathPrefix("").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)
	
	// Room endpoints
	apiRouter.HandleFunc("/rooms/{id:[0-9]+}/availability", roomHandler.CheckRoomAvailability).Methods("GET")
	
	// Reservation endpoints
	apiRouter.HandleFunc("/reservations", reservationHandler.GetAllReservations).Methods("GET")
	apiRouter.HandleFunc("/reservations/by-date", reservationHandler.GetReservationsByDate).Methods("GET")
	apiRouter.HandleFunc("/reservations", reservationHandler.CreateReservation).Methods("POST")
	apiRouter.HandleFunc("/reservations/{id:[0-9]+}", reservationHandler.GetReservation).Methods("GET")
	apiRouter.HandleFunc("/reservations/{id:[0-9]+}", reservationHandler.UpdateReservation).Methods("PUT")
	apiRouter.HandleFunc("/reservations/{id:[0-9]+}", reservationHandler.DeleteReservation).Methods("DELETE")
	
	// Admin-only endpoints
	adminRouter := r.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware, middleware.AdminRequired)
	adminRouter.HandleFunc("/rooms", roomHandler.CreateRoom).Methods("POST")
	adminRouter.HandleFunc("/rooms/{id:[0-9]+}", roomHandler.UpdateRoom).Methods("PUT")
	adminRouter.HandleFunc("/rooms/{id:[0-9]+}", roomHandler.DeleteRoom).Methods("DELETE")
	
	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	
	// Wrap the router with the CORS middleware
	handler := c.Handler(r)
	
	// Start server
	log.Printf("Reservation service starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}