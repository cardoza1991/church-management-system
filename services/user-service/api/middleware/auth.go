package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cardoza1991/church-management-system/services/user-service/internal/auth"
)

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
		
		// Verify the token
		claims, err := auth.VerifyToken(tokenParts[1])
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		
		// Add claims to request context
		ctx := context.WithValue(r.Context(), "user", claims)
		
		// Call the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
