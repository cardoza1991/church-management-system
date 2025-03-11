package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/study-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/study-service/internal/models"
)

// StudyHandler handles study-related requests
type StudyHandler struct {
	StudyRepo  *models.StudyRepository
	LessonRepo *models.LessonRepository
}

// GetStudiesByContact returns all studies for a specific contact
func (h *StudyHandler) GetStudiesByContact(w http.ResponseWriter, r *http.Request) {
	// Get contact ID from URL
	vars := mux.Vars(r)
	contactID, err := strconv.Atoi(vars["contactId"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	
	// Fetch studies from repository
	studies, err := h.StudyRepo.GetByContactID(contactID)
	if err != nil {
		http.Error(w, "Failed to fetch studies: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"contact_id": contactID,
		"studies":    studies,
	})
}

// GetStudy returns a single study by ID
func (h *StudyHandler) GetStudy(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid study ID", http.StatusBadRequest)
		return
	}
	
	// Fetch study from repository
	study, err := h.StudyRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, study)
}

// StudyRequest represents a request to create or update a study
type StudyRequest struct {
	ContactID       int    `json:"contact_id"`
	LessonID        int    `json:"lesson_id"`
	DateCompleted   string `json:"date_completed"` // Format: YYYY-MM-DD
	Location        string `json:"location,omitempty"`
	DurationMinutes int    `json:"duration_minutes,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

// CreateStudy handles creating a new study
func (h *StudyHandler) CreateStudy(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Parse request
	var req StudyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.ContactID <= 0 {
		http.Error(w, "Contact ID is required", http.StatusBadRequest)
		return
	}
	
	if req.LessonID <= 0 {
		http.Error(w, "Lesson ID is required", http.StatusBadRequest)
		return
	}
	
	if req.DateCompleted == "" {
		http.Error(w, "Date completed is required", http.StatusBadRequest)
		return
	}
	
	// Parse date
	dateCompleted, err := time.Parse("2006-01-02", req.DateCompleted)
	if err != nil {
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	
	// Verify that the lesson exists
	_, err = h.LessonRepo.GetByID(req.LessonID)
	if err != nil {
		http.Error(w, "Invalid lesson ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Check if this study already exists for this contact and lesson
	studies, err := h.StudyRepo.GetByContactID(req.ContactID)
	if err != nil {
		http.Error(w, "Failed to check for duplicates: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	for _, study := range studies {
		if study.LessonID == req.LessonID {
			http.Error(w, fmt.Sprintf("This contact has already completed lesson #%d", req.LessonID), http.StatusBadRequest)
			return
		}
	}
	
	// Create study
	study := &models.Study{
		ContactID:       req.ContactID,
		LessonID:        req.LessonID,
		DateCompleted:   dateCompleted,
		Location:        req.Location,
		DurationMinutes: req.DurationMinutes,
		Notes:           req.Notes,
		TaughtByUserID:  claims.UserID, // Set the current user as the teacher
	}
	
	// Save to database
	if err := h.StudyRepo.Create(study); err != nil {
		http.Error(w, "Failed to create study: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the created study with lesson title
	createdStudy, err := h.StudyRepo.GetByID(study.ID)
	if err != nil {
		http.Error(w, "Study created but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusCreated, createdStudy)
}

// UpdateStudy handles updating an existing study
func (h *StudyHandler) UpdateStudy(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid study ID", http.StatusBadRequest)
		return
	}
	
	// Check if study exists
	existingStudy, err := h.StudyRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Parse request
	var req StudyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.ContactID <= 0 {
		http.Error(w, "Contact ID is required", http.StatusBadRequest)
		return
	}
	
	if req.LessonID <= 0 {
		http.Error(w, "Lesson ID is required", http.StatusBadRequest)
		return
	}
	
	if req.DateCompleted == "" {
		http.Error(w, "Date completed is required", http.StatusBadRequest)
		return
	}
	
	// Parse date
	dateCompleted, err := time.Parse("2006-01-02", req.DateCompleted)
	if err != nil {
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	
	// Verify that the lesson exists
	_, err = h.LessonRepo.GetByID(req.LessonID)
	if err != nil {
		http.Error(w, "Invalid lesson ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// If changing the lesson, check if this contact already completed that lesson
	if req.LessonID != existingStudy.LessonID {
		studies, err := h.StudyRepo.GetByContactID(req.ContactID)
		if err != nil {
			http.Error(w, "Failed to check for duplicates: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		for _, study := range studies {
			if study.ID != id && study.LessonID == req.LessonID {
				http.Error(w, fmt.Sprintf("This contact has already completed lesson #%d", req.LessonID), http.StatusBadRequest)
				return
			}
		}
	}
	
	// Update study
	study := &models.Study{
		ID:              id,
		ContactID:       req.ContactID,
		LessonID:        req.LessonID,
		DateCompleted:   dateCompleted,
		Location:        req.Location,
		DurationMinutes: req.DurationMinutes,
		Notes:           req.Notes,
		TaughtByUserID:  existingStudy.TaughtByUserID, // Preserve the original teacher
	}
	
	// Save to database
	if err := h.StudyRepo.Update(id, study); err != nil {
		http.Error(w, "Failed to update study: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the updated study with lesson title
	updatedStudy, err := h.StudyRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Study updated but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, updatedStudy)
}

// DeleteStudy handles deleting a study
func (h *StudyHandler) DeleteStudy(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid study ID", http.StatusBadRequest)
		return
	}
	
	// Check if study exists
	_, err = h.StudyRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Delete from database
	if err := h.StudyRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete study: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Study deleted successfully",
	})
}

// GetContactStudyStats returns statistics about a contact's Bible study progress
func (h *StudyHandler) GetContactStudyStats(w http.ResponseWriter, r *http.Request) {
	// Get contact ID from URL
	vars := mux.Vars(r)
	contactID, err := strconv.Atoi(vars["contactId"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	
	// Get statistics
	stats, err := h.StudyRepo.GetContactStudyStats(contactID)
	if err != nil {
		http.Error(w, "Failed to get study statistics: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, stats)
}

// GetCompletedLessons returns a list of all lessons and marks which ones have been completed by a contact
func (h *StudyHandler) GetCompletedLessons(w http.ResponseWriter, r *http.Request) {
	// Get contact ID from URL
	vars := mux.Vars(r)
	contactID, err := strconv.Atoi(vars["contactId"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	
	// Get all lessons
	lessons, err := h.LessonRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch lessons: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get completed lessons for this contact
	completedLessons, err := h.StudyRepo.GetCompletedLessonsByContactID(contactID)
	if err != nil {
		http.Error(w, "Failed to fetch completed lessons: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Create response with completed status
	type LessonStatus struct {
		*models.Lesson
		Completed bool `json:"completed"`
	}
	
	var lessonStatuses []LessonStatus
	for _, lesson := range lessons {
		lessonStatuses = append(lessonStatuses, LessonStatus{
			Lesson:    lesson,
			Completed: completedLessons[lesson.ID],
		})
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"contact_id": contactID,
		"lessons":    lessonStatuses,
	})
}