package models

import (
	"database/sql"
	"errors"
	"time"
)

// Reservation represents a room booking
type Reservation struct {
	ID                int       `json:"id"`
	RoomID            int       `json:"room_id"`
	RoomName          string    `json:"room_name,omitempty"`
	UserID            int       `json:"user_id"`
	ContactID         int       `json:"contact_id,omitempty"`
	Title             string    `json:"title"`
	Description       string    `json:"description,omitempty"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	RecurringType     string    `json:"recurring_type"`
	RecurringEndDate  time.Time `json:"recurring_end_date,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// RecurringTypes defines valid recurring types
var RecurringTypes = map[string]bool{
	"none":    true,
	"daily":   true,
	"weekly":  true,
	"monthly": true,
}

// ReservationRepository provides access to the reservation store
type ReservationRepository struct {
	DB *sql.DB
}

// NewReservationRepository creates a new ReservationRepository
func NewReservationRepository(db *sql.DB) *ReservationRepository {
	return &ReservationRepository{DB: db}
}

// GetAll retrieves all reservations
func (r *ReservationRepository) GetAll(limit, offset int) ([]*Reservation, error) {
	query := `
		SELECT r.id, r.room_id, m.name, r.user_id, r.contact_id, r.title, r.description, 
			   r.start_time, r.end_time, r.recurring_type, r.recurring_end_date, 
			   r.created_at, r.updated_at
		FROM reservations r
		JOIN rooms m ON r.room_id = m.id
		ORDER BY r.start_time DESC
		LIMIT ? OFFSET ?
	`
	
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var reservations []*Reservation
	for rows.Next() {
		reservation := &Reservation{}
		var contactID sql.NullInt64
		var recurringEndDate sql.NullTime
		
		err := rows.Scan(
			&reservation.ID, 
			&reservation.RoomID, 
			&reservation.RoomName, 
			&reservation.UserID, 
			&contactID, 
			&reservation.Title, 
			&reservation.Description, 
			&reservation.StartTime, 
			&reservation.EndTime, 
			&reservation.RecurringType, 
			&recurringEndDate, 
			&reservation.CreatedAt, 
			&reservation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if contactID.Valid {
			reservation.ContactID = int(contactID.Int64)
		}
		
		if recurringEndDate.Valid {
			reservation.RecurringEndDate = recurringEndDate.Time
		}
		
		reservations = append(reservations, reservation)
	}
	
	return reservations, nil
}

// GetByID retrieves a reservation by ID
func (r *ReservationRepository) GetByID(id int) (*Reservation, error) {
	query := `
		SELECT r.id, r.room_id, m.name, r.user_id, r.contact_id, r.title, r.description, 
			   r.start_time, r.end_time, r.recurring_type, r.recurring_end_date, 
			   r.created_at, r.updated_at
		FROM reservations r
		JOIN rooms m ON r.room_id = m.id
		WHERE r.id = ?
	`
	
	reservation := &Reservation{}
	var contactID sql.NullInt64
	var recurringEndDate sql.NullTime
	
	err := r.DB.QueryRow(query, id).Scan(
		&reservation.ID, 
		&reservation.RoomID, 
		&reservation.RoomName, 
		&reservation.UserID, 
		&contactID, 
		&reservation.Title, 
		&reservation.Description, 
		&reservation.StartTime, 
		&reservation.EndTime, 
		&reservation.RecurringType, 
		&recurringEndDate, 
		&reservation.CreatedAt, 
		&reservation.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("reservation not found")
		}
		return nil, err
	}
	
	if contactID.Valid {
		reservation.ContactID = int(contactID.Int64)
	}
	
	if recurringEndDate.Valid {
		reservation.RecurringEndDate = recurringEndDate.Time
	}
	
	return reservation, nil
}

// GetByRoomID retrieves reservations for a specific room
func (r *ReservationRepository) GetByRoomID(roomID int, startDate, endDate time.Time) ([]*Reservation, error) {
	query := `
		SELECT r.id, r.room_id, m.name, r.user_id, r.contact_id, r.title, r.description, 
			   r.start_time, r.end_time, r.recurring_type, r.recurring_end_date, 
			   r.created_at, r.updated_at
		FROM reservations r
		JOIN rooms m ON r.room_id = m.id
		WHERE r.room_id = ? AND (
			(r.start_time >= ? AND r.start_time < ?) OR
			(r.end_time > ? AND r.end_time <= ?) OR
			(r.start_time <= ? AND r.end_time >= ?)
		)
		ORDER BY r.start_time
	`
	
	rows, err := r.DB.Query(query, roomID, startDate, endDate, startDate, endDate, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var reservations []*Reservation
	for rows.Next() {
		reservation := &Reservation{}
		var contactID sql.NullInt64
		var recurringEndDate sql.NullTime
		
		err := rows.Scan(
			&reservation.ID, 
			&reservation.RoomID, 
			&reservation.RoomName, 
			&reservation.UserID, 
			&contactID, 
			&reservation.Title, 
			&reservation.Description, 
			&reservation.StartTime, 
			&reservation.EndTime, 
			&reservation.RecurringType, 
			&recurringEndDate, 
			&reservation.CreatedAt, 
			&reservation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if contactID.Valid {
			reservation.ContactID = int(contactID.Int64)
		}
		
		if recurringEndDate.Valid {
			reservation.RecurringEndDate = recurringEndDate.Time
		}
		
		reservations = append(reservations, reservation)
	}
	
	return reservations, nil
}

// GetByDateRange retrieves all reservations within a date range
func (r *ReservationRepository) GetByDateRange(startDate, endDate time.Time) ([]*Reservation, error) {
	query := `
		SELECT r.id, r.room_id, m.name, r.user_id, r.contact_id, r.title, r.description, 
			   r.start_time, r.end_time, r.recurring_type, r.recurring_end_date, 
			   r.created_at, r.updated_at
		FROM reservations r
		JOIN rooms m ON r.room_id = m.id
		WHERE (
			(r.start_time >= ? AND r.start_time < ?) OR
			(r.end_time > ? AND r.end_time <= ?) OR
			(r.start_time <= ? AND r.end_time >= ?)
		)
		ORDER BY r.start_time
	`
	
	rows, err := r.DB.Query(query, startDate, endDate, startDate, endDate, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var reservations []*Reservation
	for rows.Next() {
		reservation := &Reservation{}
		var contactID sql.NullInt64
		var recurringEndDate sql.NullTime
		
		err := rows.Scan(
			&reservation.ID, 
			&reservation.RoomID, 
			&reservation.RoomName, 
			&reservation.UserID, 
			&contactID, 
			&reservation.Title, 
			&reservation.Description, 
			&reservation.StartTime, 
			&reservation.EndTime, 
			&reservation.RecurringType, 
			&recurringEndDate, 
			&reservation.CreatedAt, 
			&reservation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if contactID.Valid {
			reservation.ContactID = int(contactID.Int64)
		}
		
		if recurringEndDate.Valid {
			reservation.RecurringEndDate = recurringEndDate.Time
		}
		
		reservations = append(reservations, reservation)
	}
	
	return reservations, nil
}

// Create adds a new reservation to the database
func (r *ReservationRepository) Create(reservation *Reservation) error {
	query := `
		INSERT INTO reservations 
		(room_id, user_id, contact_id, title, description, start_time, end_time, recurring_type, recurring_end_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	var contactID interface{} = nil
	if reservation.ContactID > 0 {
		contactID = reservation.ContactID
	}
	
	var recurringEndDate interface{} = nil
	if !reservation.RecurringEndDate.IsZero() && reservation.RecurringType != "none" {
		recurringEndDate = reservation.RecurringEndDate
	}
	
	result, err := r.DB.Exec(
		query, 
		reservation.RoomID, 
		reservation.UserID, 
		contactID, 
		reservation.Title, 
		reservation.Description, 
		reservation.StartTime, 
		reservation.EndTime, 
		reservation.RecurringType, 
		recurringEndDate,
	)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	reservation.ID = int(id)
	return nil
}

// Update modifies an existing reservation
func (r *ReservationRepository) Update(id int, reservation *Reservation) error {
	query := `
		UPDATE reservations 
		SET room_id = ?, contact_id = ?, title = ?, description = ?, 
			start_time = ?, end_time = ?, recurring_type = ?, recurring_end_date = ?
		WHERE id = ?
	`
	
	var contactID interface{} = nil
	if reservation.ContactID > 0 {
		contactID = reservation.ContactID
	}
	
	var recurringEndDate interface{} = nil
	if !reservation.RecurringEndDate.IsZero() && reservation.RecurringType != "none" {
		recurringEndDate = reservation.RecurringEndDate
	}
	
	_, err := r.DB.Exec(
		query, 
		reservation.RoomID, 
		contactID, 
		reservation.Title, 
		reservation.Description, 
		reservation.StartTime, 
		reservation.EndTime, 
		reservation.RecurringType, 
		recurringEndDate,
		id,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

// Delete removes a reservation from the database
func (r *ReservationRepository) Delete(id int) error {
	query := `DELETE FROM reservations WHERE id = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	return nil
}