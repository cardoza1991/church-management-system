package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/reservation-service/api/middleware"
	"github.com/cardoza1991/church-management-system/services/reservation-service/internal/models"
)

// ReservationHandler handles reservation-related requests
type ReservationHandler struct {
	ReservationRepo *models.ReservationRepository
	RoomRepo        *models.RoomRepository
}

// GetAllReservations returns a list of all reservations
func (h *ReservationHandler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
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
	
	// Fetch reservations from repository
	reservations, err := h.ReservationRepo.GetAll(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch reservations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"reservations": reservations,
		"limit":        limit,
		"offset":       offset,
	})
}

// GetReservationsByDate returns reservations for a specific date range
func (h *ReservationHandler) GetReservationsByDate(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	roomIDStr := r.URL.Query().Get("room_id")
	
	// Default to today if no date range provided
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.Add(24 * time.Hour)
	
	if startStr != "" {
		parsedStart, err := time.Parse("2006-01-02", startStr)
		if err == nil {
			startDate = parsedStart
		} else {
			http.Error(w, "Invalid start date format (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	}
	
	if endStr != "" {
		parsedEnd, err := time.Parse("2006-01-02", endStr)
		if err == nil {
			endDate = parsedEnd.Add(24 * time.Hour) // Include the whole end day
		} else {
			http.Error(w, "Invalid end date format (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
	}
	
	// If room ID provided, get reservations for that room
	if roomIDStr != "" {
		roomID, err := strconv.Atoi(roomIDStr)
		if err != nil {
			http.Error(w, "Invalid room ID", http.StatusBadRequest)
			return
		}
		
		// Check if room exists
		_, err = h.RoomRepo.GetByID(roomID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		
		reservations, err := h.ReservationRepo.GetByRoomID(roomID, startDate, endDate)
		if err != nil {
			http.Error(w, "Failed to fetch reservations: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"room_id":      roomID,
			"start_date":   startDate.Format("2006-01-02"),
			"end_date":     endDate.Add(-time.Second).Format("2006-01-02"),
			"reservations": reservations,
		})
		return
	}
	
	// Get all reservations within the date range
	reservations, err := h.ReservationRepo.GetByDateRange(startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to fetch reservations: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"start_date":   startDate.Format("2006-01-02"),
		"end_date":     endDate.Add(-time.Second).Format("2006-01-02"),
		"reservations": reservations,
	})
}

// GetReservation returns a single reservation by ID
func (h *ReservationHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}
	
	// Fetch reservation from repository
	reservation, err := h.ReservationRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, reservation)
}

// ReservationRequest represents a request to create or update a reservation
type ReservationRequest struct {
	RoomID           int    `json:"room_id"`
	ContactID        int    `json:"contact_id,omitempty"`
	Title            string `json:"title"`
	Description      string `json:"description,omitempty"`
	StartTime        string `json:"start_time"` // ISO 8601 format
	EndTime          string `json:"end_time"`   // ISO 8601 format
	RecurringType    string `json:"recurring_type,omitempty"`
	RecurringEndDate string `json:"recurring_end_date,omitempty"` // YYYY-MM-DD format
}

// CreateReservation handles creating a new reservation
func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Parse request
	var req ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.RoomID <= 0 {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	
	if req.StartTime == "" || req.EndTime == "" {
		http.Error(w, "Start and end times are required", http.StatusBadRequest)
		return
	}
	
	// Parse times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format (ISO 8601)", http.StatusBadRequest)
		return
	}
	
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format (ISO 8601)", http.StatusBadRequest)
		return
	}
	
	// Validate time range
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}
	
	// Check if the room exists
	_, err = h.RoomRepo.GetByID(req.RoomID)
	if err != nil {
		http.Error(w, "Invalid room ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Check recurring type
	recurringType := "none"
	if req.RecurringType != "" {
		if !models.RecurringTypes[req.RecurringType] {
			http.Error(w, "Invalid recurring type", http.StatusBadRequest)
			return
		}
		recurringType = req.RecurringType
	}
	
	// Parse recurring end date if provided
	var recurringEndDate time.Time
	if req.RecurringEndDate != "" && recurringType != "none" {
		recurringEndDate, err = time.Parse("2006-01-02", req.RecurringEndDate)
		if err != nil {
			http.Error(w, "Invalid recurring end date format (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		
		// Ensure recurring end date is after the start date
		if recurringEndDate.Before(startTime) {
			http.Error(w, "Recurring end date must be after the start date", http.StatusBadRequest)
			return
		}
	}
	
	// Check if the room is available for the requested time
	isAvailable, err := h.RoomRepo.CheckAvailability(req.RoomID, startTime, endTime)
	if err != nil {
		http.Error(w, "Failed to check availability: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	if !isAvailable {
		http.Error(w, "Room is not available for the requested time", http.StatusConflict)
		return
	}
	
	// Create reservation
	reservation := &models.Reservation{
		RoomID:          req.RoomID,
		UserID:          claims.UserID,
		ContactID:       req.ContactID,
		Title:           req.Title,
		Description:     req.Description,
		StartTime:       startTime,
		EndTime:         endTime,
		RecurringType:   recurringType,
		RecurringEndDate: recurringEndDate,
	}
	
	// Save to database
	if err := h.ReservationRepo.Create(reservation); err != nil {
		http.Error(w, "Failed to create reservation: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the created reservation with room name
	createdReservation, err := h.ReservationRepo.GetByID(reservation.ID)
	if err != nil {
		http.Error(w, "Reservation created but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusCreated, createdReservation)
}

// UpdateReservation handles updating an existing reservation
func (h *ReservationHandler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}
	
	// Check if reservation exists
	existingReservation, err := h.ReservationRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Only the user who created the reservation or an admin can update it
	if existingReservation.UserID != claims.UserID && claims.Role != "admin" {
		http.Error(w, "You can only update your own reservations", http.StatusForbidden)
		return
	}
	
	// Parse request
	var req ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.RoomID <= 0 {
		http.Error(w, "Room ID is required", http.StatusBadRequest)
		return
	}
	
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	
	if req.StartTime == "" || req.EndTime == "" {
		http.Error(w, "Start and end times are required", http.StatusBadRequest)
		return
	}
	
	// Parse times
	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		http.Error(w, "Invalid start time format (ISO 8601)", http.StatusBadRequest)
		return
	}
	
	endTime, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		http.Error(w, "Invalid end time format (ISO 8601)", http.StatusBadRequest)
		return
	}
	
	// Validate time range
	if endTime.Before(startTime) || endTime.Equal(startTime) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}
	
	// Check if the room exists
	_, err = h.RoomRepo.GetByID(req.RoomID)
	if err != nil {
		http.Error(w, "Invalid room ID: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Check recurring type
	recurringType := "none"
	if req.RecurringType != "" {
		if !models.RecurringTypes[req.RecurringType] {
			http.Error(w, "Invalid recurring type", http.StatusBadRequest)
			return
		}
		recurringType = req.RecurringType
	}
	
	// Parse recurring end date if provided
	var recurringEndDate time.Time
	if req.RecurringEndDate != "" && recurringType != "none" {
		recurringEndDate, err = time.Parse("2006-01-02", req.RecurringEndDate)
		if err != nil {
			http.Error(w, "Invalid recurring end date format (YYYY-MM-DD)", http.StatusBadRequest)
			return
		}
		
		// Ensure recurring end date is after the start date
		if recurringEndDate.Before(startTime) {
			http.Error(w, "Recurring end date must be after the start date", http.StatusBadRequest)
			return
		}
	}
	
	// Check if the room is available for the requested time
	// Skip availability check if room and time are unchanged
	if req.RoomID != existingReservation.RoomID || 
	   !startTime.Equal(existingReservation.StartTime) || 
	   !endTime.Equal(existingReservation.EndTime) {
		
		isAvailable, err := h.RoomRepo.CheckAvailability(req.RoomID, startTime, endTime)
		if err != nil {
			http.Error(w, "Failed to check availability: "+err.Error(), http.StatusInternalServerError)
			return
		}
		
		if !isAvailable {
			http.Error(w, "Room is not available for the requested time", http.StatusConflict)
			return
		}
	}
	
	// Update reservation
	reservation := &models.Reservation{
		ID:              id,
		RoomID:          req.RoomID,
		UserID:          existingReservation.UserID, // Preserve original user
		ContactID:       req.ContactID,
		Title:           req.Title,
		Description:     req.Description,
		StartTime:       startTime,
		EndTime:         endTime,
		RecurringType:   recurringType,
		RecurringEndDate: recurringEndDate,
	}
	
	// Save to database
	if err := h.ReservationRepo.Update(id, reservation); err != nil {
		http.Error(w, "Failed to update reservation: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the updated reservation with room name
	updatedReservation, err := h.ReservationRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Reservation updated but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, updatedReservation)
}

// DeleteReservation handles deleting a reservation
func (h *ReservationHandler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}
	
	// Check if reservation exists
	existingReservation, err := h.ReservationRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Only the user who created the reservation or an admin can delete it
	if existingReservation.UserID != claims.UserID && claims.Role != "admin" {
		http.Error(w, "You can only delete your own reservations", http.StatusForbidden)
		return
	}
	
	// Delete reservation
	if err := h.ReservationRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete reservation: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Reservation deleted successfully",
	})
}