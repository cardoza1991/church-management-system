package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/study-service/api/handlers"
	"github.com/cardoza1991/church-management-system/services/study-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/study-service/config"
	"github.com/cardoza1991/church-management-system/services/study-service/internal/db"
	"github.com/cardoza1991/church-management-system/services/study-service/internal/models"
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
	lessonRepo := models.NewLessonRepository(database)
	studyRepo := models.NewStudyRepository(database)
	
	// Create handlers
	lessonHandler := &handlers.LessonHandler{
		LessonRepo: lessonRepo,
	}
	studyHandler := &handlers.StudyHandler{
		StudyRepo:  studyRepo,
		LessonRepo: lessonRepo,
	}
	
	// Create router
	r := mux.NewRouter()
	
	// Health check endpoint
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Public endpoints (no authentication required)
	r.HandleFunc("/lessons", lessonHandler.GetAllLessons).Methods("GET")
	r.HandleFunc("/lessons/{id:[0-9]+}", lessonHandler.GetLesson).Methods("GET")
	
	// Protected endpoints (require authentication)
	apiRouter := r.PathPrefix("").Subrouter()
	apiRouter.Use(middleware.AuthMiddleware)
	
	// Studies endpoints
	apiRouter.HandleFunc("/contacts/{contactId:[0-9]+}/studies", studyHandler.GetStudiesByContact).Methods("GET")
	apiRouter.HandleFunc("/contacts/{contactId:[0-9]+}/study-stats", studyHandler.GetContactStudyStats).Methods("GET")
	apiRouter.HandleFunc("/contacts/{contactId:[0-9]+}/completed-lessons", studyHandler.GetCompletedLessons).Methods("GET")
	apiRouter.HandleFunc("/studies", studyHandler.CreateStudy).Methods("POST")
	apiRouter.HandleFunc("/studies/{id:[0-9]+}", studyHandler.GetStudy).Methods("GET")
	apiRouter.HandleFunc("/studies/{id:[0-9]+}", studyHandler.UpdateStudy).Methods("PUT")
	apiRouter.HandleFunc("/studies/{id:[0-9]+}", studyHandler.DeleteStudy).Methods("DELETE")
	
	// Admin-only endpoints
	adminRouter := r.PathPrefix("").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware, middleware.AdminRequired)
	adminRouter.HandleFunc("/lessons", lessonHandler.CreateLesson).Methods("POST")
	adminRouter.HandleFunc("/lessons/{id:[0-9]+}", lessonHandler.UpdateLesson).Methods("PUT")
	adminRouter.HandleFunc("/lessons/{id:[0-9]+}", lessonHandler.DeleteLesson).Methods("DELETE")
	
	// Start server
	log.Printf("Study service starting on port %s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}