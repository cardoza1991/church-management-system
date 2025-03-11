package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/study-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/study-service/internal/models"
)

// LessonHandler handles lesson-related requests
type LessonHandler struct {
	LessonRepo *models.LessonRepository
}

// GetAllLessons returns all lessons
func (h *LessonHandler) GetAllLessons(w http.ResponseWriter, r *http.Request) {
	// Fetch lessons from repository
	lessons, err := h.LessonRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch lessons: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"lessons": lessons,
	})
}

// GetLesson returns a single lesson by ID
func (h *LessonHandler) GetLesson(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}
	
	// Fetch lesson from repository
	lesson, err := h.LessonRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, lesson)
}

// LessonRequest represents a request to create or update a lesson
type LessonRequest struct {
	Title          string `json:"title"`
	Description    string `json:"description,omitempty"`
	SequenceNumber int    `json:"sequence_number"`
}

// CreateLesson handles creating a new lesson
func (h *LessonHandler) CreateLesson(w http.ResponseWriter, r *http.Request) {
	// Only admins can create lessons
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Parse request
	var req LessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	
	if req.SequenceNumber <= 0 {
		http.Error(w, "Sequence number must be greater than zero", http.StatusBadRequest)
		return
	}
	
	// Check for duplicates (by title)
	existingLessons, err := h.LessonRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to check for duplicates: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	for _, lesson := range existingLessons {
		if lesson.Title == req.Title {
			http.Error(w, "A lesson with this title already exists", http.StatusBadRequest)
			return
		}
		if lesson.SequenceNumber == req.SequenceNumber {
			http.Error(w, "A lesson with this sequence number already exists", http.StatusBadRequest)
			return
		}
	}
	
	// Create lesson
	lesson := &models.Lesson{
		Title:          req.Title,
		Description:    req.Description,
		SequenceNumber: req.SequenceNumber,
	}
	
	// Save to database
	if err := h.LessonRepo.Create(lesson); err != nil {
		http.Error(w, "Failed to create lesson: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusCreated, lesson)
}

// UpdateLesson handles updating an existing lesson
func (h *LessonHandler) UpdateLesson(w http.ResponseWriter, r *http.Request) {
	// Only admins can update lessons
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}
	
	// Check if lesson exists
	existingLesson, err := h.LessonRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Parse request
	var req LessonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	
	if req.SequenceNumber <= 0 {
		http.Error(w, "Sequence number must be greater than zero", http.StatusBadRequest)
		return
	}
	
	// Check for duplicates (by title and sequence number) - ignore the current lesson being updated
	existingLessons, err := h.LessonRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to check for duplicates: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	for _, lesson := range existingLessons {
		if lesson.ID != id { // Skip the lesson being updated
			if lesson.Title == req.Title {
				http.Error(w, "A lesson with this title already exists", http.StatusBadRequest)
				return
			}
			if lesson.SequenceNumber == req.SequenceNumber {
				http.Error(w, "A lesson with this sequence number already exists", http.StatusBadRequest)
				return
			}
		}
	}
	
	// Update lesson
	lesson := &models.Lesson{
		ID:             id,
		Title:          req.Title,
		Description:    req.Description,
		SequenceNumber: req.SequenceNumber,
	}
	
	// Save to database
	if err := h.LessonRepo.Update(id, lesson); err != nil {
		http.Error(w, "Failed to update lesson: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the updated lesson to return
	updatedLesson, err := h.LessonRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Lesson updated but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, updatedLesson)
}

// DeleteLesson handles deleting a lesson
func (h *LessonHandler) DeleteLesson(w http.ResponseWriter, r *http.Request) {
	// Only admins can delete lessons
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid lesson ID", http.StatusBadRequest)
		return
	}
	
	// Delete from database (this will fail if lesson is in use)
	if err := h.LessonRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete lesson: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Lesson deleted successfully",
	})
}