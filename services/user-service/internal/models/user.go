package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	FullName     string    `json:"full_name"`
	Phone        string    `json:"phone,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserRepository provides access to the user store
type UserRepository struct {
	DB *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create adds a new user to the database
func (r *UserRepository) Create(user *User, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	
	query := `INSERT INTO users (username, password_hash, email, role, full_name, phone)
		VALUES (?, ?, ?, ?, ?, ?)`
	
	result, err := r.DB.Exec(query, user.Username, string(hashedPassword), 
		user.Email, user.Role, user.FullName, user.Phone)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	user.ID = int(id)
	return nil
}

// GetByUsername finds a user by username
func (r *UserRepository) GetByUsername(username string) (*User, error) {
	user := &User{}
	
	query := `SELECT id, username, password_hash, email, role, full_name, phone, created_at, updated_at 
		FROM users WHERE username = ?`
	
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Email, 
		&user.Role, &user.FullName, &user.Phone, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	
	return user, nil
}

// CheckPassword verifies a user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
