package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/contact-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/contact-service/internal/models"
)

// StatusHandler handles status-related requests
type StatusHandler struct {
	StatusRepo *models.StatusRepository
}

// GetAllStatuses returns all available statuses
func (h *StatusHandler) GetAllStatuses(w http.ResponseWriter, r *http.Request) {
	// Fetch statuses from repository
	statuses, err := h.StatusRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch statuses: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"statuses": statuses,
	})
}

// GetStatus returns a single status by ID
func (h *StatusHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid status ID", http.StatusBadRequest)
		return
	}
	
	// Fetch status from repository
	status, err := h.StatusRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, status)
}

// StatusRequest represents a request to create or update a status
type StatusRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DisplayOrder int   `json:"display_order"`
}

// CreateStatus handles creating a new status
func (h *StatusHandler) CreateStatus(w http.ResponseWriter, r *http.Request) {
	// Only admins can create statuses
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Parse request
	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Create status
	status := &models.Status{
		Name:        req.Name,
		Description: req.Description,
		DisplayOrder: req.DisplayOrder,
	}
	
	// Save to database
	if err := h.StatusRepo.Create(status); err != nil {
		http.Error(w, "Failed to create status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusCreated, status)
}

// UpdateStatus handles updating an existing status
func (h *StatusHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	// Only admins can update statuses
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid status ID", http.StatusBadRequest)
		return
	}
	
	// Check if status exists
	_, err = h.StatusRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Parse request
	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Update status
	status := &models.Status{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		DisplayOrder: req.DisplayOrder,
	}
	
	// Save to database
	if err := h.StatusRepo.Update(id, status); err != nil {
		http.Error(w, "Failed to update status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, status)
}

// DeleteStatus handles deleting a status
func (h *StatusHandler) DeleteStatus(w http.ResponseWriter, r *http.Request) {
	// Only admins can delete statuses
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid status ID", http.StatusBadRequest)
		return
	}
	
	// Delete from database (this will fail if status is in use)
	if err := h.StatusRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Status deleted successfully",
	})
}