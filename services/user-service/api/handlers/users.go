package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/cardoza1991/church-management-system/services/user-service/internal/models"
)

// UserHandler handles user-related requests
type UserHandler struct {
	UserRepo *models.UserRepository
}

// GetUser returns a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	// Example response for demonstration
	user := &models.User{
		ID:       id,
		Username: "user" + vars["id"],
		Email:    "user" + vars["id"] + "@example.com",
		Role:     "member",
		FullName: "User " + vars["id"],
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetSelf returns the current user
func (h *UserHandler) GetSelf(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by AuthMiddleware)
	claims, ok := r.Context().Value("user").(*auth.Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	
	// Find user in database
	user, err := h.UserRepo.GetByUsername(claims.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	
	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
