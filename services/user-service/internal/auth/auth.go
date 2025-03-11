package auth

import (
	"errors"
	"time"
	"fmt"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	ExpiresAt int64 `json:"exp"`
}

// For a real app, use a proper JWT library and store the secret securely
var hmacSecret = []byte("your-secret-key")

// GenerateToken creates a new JWT token for a user
func GenerateToken(userID int, username, role string) (string, error) {
	// Set expiration time - 24 hours from now
	expirationTime := time.Now().Add(24 * time.Hour).Unix()
	
	// Create claims
	claims := Claims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		ExpiresAt: expirationTime,
	}
	
	// Create token
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	
	// Base64 encode the header and payload
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload64 := base64.RawURLEncoding.EncodeToString(payload)
	
	// Create signature
	h := hmac.New(sha256.New, hmacSecret)
	h.Write([]byte(header + "." + payload64))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	
	// Combine to create token
	token := header + "." + payload64 + "." + signature
	
	return token, nil
}

// VerifyToken validates a JWT token
func VerifyToken(tokenString string) (*Claims, error) {
	// Split the token
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	
	// Get the header and payload
	header, payload, signature := parts[0], parts[1], parts[2]
	
	// Verify signature
	h := hmac.New(sha256.New, hmacSecret)
	h.Write([]byte(header + "." + payload))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	
	if signature != expectedSignature {
		return nil, errors.New("invalid token signature")
	}
	
	// Decode payload
	decodedPayload, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, errors.New("failed to decode payload")
	}
	
	// Parse claims
	var claims Claims
	if err := json.Unmarshal(decodedPayload, &claims); err != nil {
		return nil, errors.New("failed to parse claims")
	}
	
	// Check if token is expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token is expired")
	}
	
	return &claims, nil
}
