package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// Claims represents the JWT claims we expect from the auth service
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// AuthMiddleware checks for a valid JWT token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}
		
		// Extract the token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		
		// Note: We would normally verify the token here
		// but for simplicity, we're just mocking the verification
		
		// Parse token claims (hardcoded example for demonstration)
		claims := &Claims{
			UserID:   1,
			Username: "demo_user",
			Role:     "admin",
		}
		
		// Add claims to request context
		ctx := context.WithValue(r.Context(), "user", claims)
		
		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminRequired ensures the user has admin role
func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user from context (set by AuthMiddleware)
		claims, ok := r.Context().Value("user").(*Claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		
		// Check if user has admin role
		if claims.Role != "admin" {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// RespondJSON is a helper function to respond with JSON
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}