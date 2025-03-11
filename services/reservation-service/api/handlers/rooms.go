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

// RoomHandler handles room-related requests
type RoomHandler struct {
	RoomRepo *models.RoomRepository
}

// GetAllRooms returns a list of all rooms
func (h *RoomHandler) GetAllRooms(w http.ResponseWriter, r *http.Request) {
	// Fetch rooms from repository
	rooms, err := h.RoomRepo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"rooms": rooms,
	})
}

// GetRoom returns a single room by ID
func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	
	// Fetch room from repository
	room, err := h.RoomRepo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, room)
}

// RoomRequest represents a request to create or update a room
type RoomRequest struct {
	Name             string `json:"name"`
	Capacity         int    `json:"capacity"`
	Location         string `json:"location,omitempty"`
	Description      string `json:"description,omitempty"`
	AvailabilityStart string `json:"availability_start,omitempty"`
	AvailabilityEnd  string `json:"availability_end,omitempty"`
	IsAvailable      bool   `json:"is_available"`
}

// CreateRoom handles creating a new room
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	// Only admins can create rooms
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Parse request
	var req RoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	if req.Capacity <= 0 {
		http.Error(w, "Capacity must be greater than zero", http.StatusBadRequest)
		return
	}
	
	// Validate time formats if provided
	if req.AvailabilityStart != "" {
		_, err := time.Parse("15:04:05", req.AvailabilityStart)
		if err != nil {
			http.Error(w, "Invalid availability start time format (HH:MM:SS)", http.StatusBadRequest)
			return
		}
	}
	
	if req.AvailabilityEnd != "" {
		_, err := time.Parse("15:04:05", req.AvailabilityEnd)
		if err != nil {
			http.Error(w, "Invalid availability end time format (HH:MM:SS)", http.StatusBadRequest)
			return
		}
	}
	
	// Create room
	room := &models.Room{
		Name:             req.Name,
		Capacity:         req.Capacity,
		Location:         req.Location,
		Description:      req.Description,
		AvailabilityStart: req.AvailabilityStart,
		AvailabilityEnd:  req.AvailabilityEnd,
		IsAvailable:      req.IsAvailable,
	}
	
	// Save to database
	if err := h.RoomRepo.Create(room); err != nil {
		http.Error(w, "Failed to create room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusCreated, room)
}

// UpdateRoom handles updating an existing room
func (h *RoomHandler) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	// Only admins can update rooms
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	
	
	// Parse request
	var req RoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate input
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	
	if req.Capacity <= 0 {
		http.Error(w, "Capacity must be greater than zero", http.StatusBadRequest)
		return
	}
	
	// Validate time formats if provided
	if req.AvailabilityStart != "" {
		_, err := time.Parse("15:04:05", req.AvailabilityStart)
		if err != nil {
			http.Error(w, "Invalid availability start time format (HH:MM:SS)", http.StatusBadRequest)
			return
		}
	}
	
	if req.AvailabilityEnd != "" {
		_, err := time.Parse("15:04:05", req.AvailabilityEnd)
		if err != nil {
			http.Error(w, "Invalid availability end time format (HH:MM:SS)", http.StatusBadRequest)
			return
		}
	}
	
	// Update room
	room := &models.Room{
		ID:               id,
		Name:             req.Name,
		Capacity:         req.Capacity,
		Location:         req.Location,
		Description:      req.Description,
		AvailabilityStart: req.AvailabilityStart,
		AvailabilityEnd:  req.AvailabilityEnd,
		IsAvailable:      req.IsAvailable,
	}
	
	// Save to database
	if err := h.RoomRepo.Update(id, room); err != nil {
		http.Error(w, "Failed to update room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the updated room to return
	updatedRoom, err := h.RoomRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Room updated but failed to retrieve: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, updatedRoom)
}

// DeleteRoom handles deleting a room
func (h *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	// Only admins can delete rooms
	claims, ok := r.Context().Value("user").(*middleware.Claims)
	if !ok || claims.Role != "admin" {
		http.Error(w, "Admin access required", http.StatusForbidden)
		return
	}
	
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	
	// Delete room
	if err := h.RoomRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete room: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Room deleted successfully",
	})
}

// GetAvailableRooms returns rooms available for a given time slot
func (h *RoomHandler) GetAvailableRooms(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	capacityStr := r.URL.Query().Get("capacity")
	
	if startStr == "" || endStr == "" {
		http.Error(w, "Start and end times are required", http.StatusBadRequest)
		return
	}
	
	// Parse start and end times
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "Invalid start time format (ISO 8601)", http.StatusBadRequest)
		return
	}
	
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "Invalid end time format (ISO 8601)", http.StatusBadRequest)
		return
	}
	
	// Validate time range
	if end.Before(start) || end.Equal(start) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}
	
	// Parse minimum capacity
	minCapacity := 1 // Default minimum capacity
	if capacityStr != "" {
		parsedCapacity, err := strconv.Atoi(capacityStr)
		if err != nil || parsedCapacity < 1 {
			http.Error(w, "Invalid capacity", http.StatusBadRequest)
			return
		}
		minCapacity = parsedCapacity
	}
	
	// Get available rooms
	rooms, err := h.RoomRepo.GetAvailableRooms(start, end, minCapacity)
	if err != nil {
		http.Error(w, "Failed to get available rooms: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"start":  start,
		"end":    end,
		"rooms":  rooms,
	})
}

// CheckRoomAvailability checks if a specific room is available for a given time slot
func (h *RoomHandler) CheckRoomAvailability(w http.ResponseWriter, r *http.Request) {
	// Get room ID from URL
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	
	// Parse query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")
	
	if startStr == "" || endStr == "" {
		http.Error(w, "Start and end times are required", http.StatusBadRequest)
		return
	}
	
	// Parse start and end times
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		http.Error(w, "Invalid start time format (ISO 8601)", http.StatusBadRequest)
		return
	}
	
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		http.Error(w, "Invalid end time format (ISO 8601)", http.StatusBadRequest)
		return
	}
	
	// Validate time range
	if end.Before(start) || end.Equal(start) {
		http.Error(w, "End time must be after start time", http.StatusBadRequest)
		return
	}
	
	// Check availability
	isAvailable, err := h.RoomRepo.CheckAvailability(roomID, start, end)
	if err != nil {
		http.Error(w, "Failed to check availability: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Return response
	middleware.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"room_id":      roomID,
		"start":        start,
		"end":          end,
		"is_available": isAvailable,
	})
}