package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/contact-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/contact-service/internal/models"
)

// ContactHandler handles contact-related requests
type ContactHandler struct {
	ContactRepo *models.ContactRepository
	StatusRepo  *models.StatusRepository
}

// ListContacts returns a list of contacts
func (h *ContactHandler) ListContacts(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	
	limit := 20 // Default limit
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}
	
	offset := 0 // Default offset
	if offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}
	
	// Fetch contacts from repository
	contacts, err := h.ContactRepo.GetAll(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch contacts: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"contacts": contacts,
		"limit":    limit,
		"offset":   offset,
	})
}

// GetContact returns a single contact by ID
func (h *ContactHandler) GetContact(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	
	// Fetch contact from repository
	contact, err := h.ContactRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, contact)
}

// ContactRequest represents a request to create or update a contact
type ContactRequest struct {
	Name           string `json:"name"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Location       string `json:"location,omitempty"`
	Notes          string `json:"notes,omitempty"`
	CurrentStatusID int   `json:"current_status_id"`
}

// CreateContact handles creating a new contact
func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Validate status ID
	if req.CurrentStatusID <= 0 {
		// Use default "New Contact" status if not provided
		statuses, err := h.StatusRepo.GetAll()
		if err != nil || len(statuses) == 0 {
			http.Error(w, "Failed to get default status", http.StatusInternalServerError)
			return
		}
		req.CurrentStatusID = statuses[0].ID // Assuming first status is "New Contact"
	} else {
		// Verify that the status exists
		_, err := h.StatusRepo.GetByID(req.CurrentStatusID)
		if err != nil {
			http.Error(w, "Invalid status ID", http.StatusBadRequest)
			return
		}
	}
	
	// Create contact
	contact := &models.Contact{
		Name:            req.Name,
		Email:           req.Email,
		Phone:           req.Phone,
		Location:        req.Location,
		Notes:           req.Notes,
		CurrentStatusID: req.CurrentStatusID,
		DateAdded:       time.Now(),
		LastUpdated:     time.Now(),
	}
	
	// Save to database
	if err := h.ContactRepo.Create(contact); err != nil {
		http.Error(w, "Failed to create contact: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Also add an entry in the status history
	err := h.ContactRepo.UpdateStatus(contact.ID, contact.CurrentStatusID, "Initial status")
	if err != nil {
		// Log the error but don't fail the request
		// In a real app, you might want to use proper logging
		println("Failed to create status history: " + err.Error())
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusCreated, contact)
}

// UpdateContact handles updating an existing contact
func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	
	// Check if contact exists
	existingContact, err := h.ContactRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Parse request
	var req ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	// Verify that the status exists if provided
	if req.CurrentStatusID > 0 {
		_, err := h.StatusRepo.GetByID(req.CurrentStatusID)
		if err != nil {
			http.Error(w, "Invalid status ID", http.StatusBadRequest)
			return
		}
	} else {
		// Keep existing status if not provided
		req.CurrentStatusID = existingContact.CurrentStatusID
	}
	
	// Check if status is being changed
	statusChanged := existingContact.CurrentStatusID != req.CurrentStatusID
	
	// Update contact
	contact := &models.Contact{
		ID:              id,
		Name:            req.Name,
		Email:           req.Email,
		Phone:           req.Phone,
		Location:        req.Location,
		Notes:           req.Notes,
		CurrentStatusID: req.CurrentStatusID,
		LastUpdated:     time.Now(),
	}
	
	// Save to database
	if err := h.ContactRepo.Update(id, contact); err != nil {
		http.Error(w, "Failed to update contact: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// If status changed, add an entry to the status history
	if statusChanged {
		err := h.ContactRepo.UpdateStatus(id, req.CurrentStatusID, "Status updated via contact edit")
		if err != nil {
			// Log the error but don't fail the request
			println("Failed to update status history: " + err.Error())
		}
	}
	
	// Get the updated contact to return
	updatedContact, err := h.ContactRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Contact updated but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, updatedContact)
}

// DeleteContact handles deleting a contact
func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	
	// Check if contact exists
	_, err = h.ContactRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Delete from database
	if err := h.ContactRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete contact: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Contact deleted successfully",
	})
}

// StatusUpdateRequest represents a request to change a contact's status
type StatusUpdateRequest struct {
	StatusID int    `json:"status_id"`
	Notes    string `json:"notes,omitempty"`
}

// UpdateContactStatus handles changing a contact's status
func (h *ContactHandler) UpdateContactStatus(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	
	// Check if contact exists
	contact, err := h.ContactRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Parse request
	var req StatusUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.StatusID <= 0 {
		http.Error(w, "Status ID is required", http.StatusBadRequest)
		return
	}
	
	// Check if status exists
	_, err = h.StatusRepo.GetByID(req.StatusID)
	if err != nil {
		http.Error(w, "Invalid status ID", http.StatusBadRequest)
		return
	}
	
	// Don't update if status hasn't changed
	if contact.CurrentStatusID == req.StatusID && req.Notes == "" {
		middleware.RespondJSON(w, http.StatusOK, map[string]string{
			"message": "Status unchanged",
		})
		return
	}
	
	// Update status
	if err := h.ContactRepo.UpdateStatus(id, req.StatusID, req.Notes); err != nil {
		http.Error(w, "Failed to update status: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the updated contact to return
	updatedContact, err := h.ContactRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Status updated but failed to retrieve contact: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, updatedContact)
}

// GetContactStatusHistory returns the status history for a contact
func (h *ContactHandler) GetContactStatusHistory(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}
	
	// Check if contact exists
	_, err = h.ContactRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Get status history
	history, err := h.ContactRepo.GetStatusHistory(id)
	if err != nil {
		http.Error(w, "Failed to get status history: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"contact_id": id,
		"history":    history,
	})
}