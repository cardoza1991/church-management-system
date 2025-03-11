package models

import (
	"database/sql"
	"errors"
	"time"
)

// Room represents a bookable room in the facility
type Room struct {
	ID               int        `json:"id"`
	Name             string     `json:"name"`
	Capacity         int        `json:"capacity"`
	Location         string     `json:"location,omitempty"`
	Description      string     `json:"description,omitempty"`
	AvailabilityStart string    `json:"availability_start,omitempty"`
	AvailabilityEnd  string     `json:"availability_end,omitempty"`
	IsAvailable      bool       `json:"is_available"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// RoomRepository provides access to the room store
type RoomRepository struct {
	DB *sql.DB
}

// NewRoomRepository creates a new RoomRepository
func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{DB: db}
}

// GetAll retrieves all rooms
func (r *RoomRepository) GetAll() ([]*Room, error) {
	query := `SELECT id, name, capacity, location, description, 
	          availability_start, availability_end, is_available, created_at, updated_at 
			  FROM rooms ORDER BY name`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rooms []*Room
	for rows.Next() {
		room := &Room{}
		var availabilityStart, availabilityEnd sql.NullString
		
		err := rows.Scan(
			&room.ID, 
			&room.Name, 
			&room.Capacity, 
			&room.Location, 
			&room.Description, 
			&availabilityStart,
			&availabilityEnd,
			&room.IsAvailable,
			&room.CreatedAt, 
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if availabilityStart.Valid {
			room.AvailabilityStart = availabilityStart.String
		}
		if availabilityEnd.Valid {
			room.AvailabilityEnd = availabilityEnd.String
		}
		
		rooms = append(rooms, room)
	}
	
	return rooms, nil
}

// GetByID retrieves a room by ID
func (r *RoomRepository) GetByID(id int) (*Room, error) {
	query := `SELECT id, name, capacity, location, description, 
	          availability_start, availability_end, is_available, created_at, updated_at 
			  FROM rooms WHERE id = ?`
	
	room := &Room{}
	var availabilityStart, availabilityEnd sql.NullString
	
	err := r.DB.QueryRow(query, id).Scan(
		&room.ID, 
		&room.Name, 
		&room.Capacity, 
		&room.Location, 
		&room.Description, 
		&availabilityStart,
		&availabilityEnd,
		&room.IsAvailable,
		&room.CreatedAt, 
		&room.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("room not found")
		}
		return nil, err
	}
	
	if availabilityStart.Valid {
		room.AvailabilityStart = availabilityStart.String
	}
	if availabilityEnd.Valid {
		room.AvailabilityEnd = availabilityEnd.String
	}
	
	return room, nil
}

// Create adds a new room to the database
func (r *RoomRepository) Create(room *Room) error {
	query := `INSERT INTO rooms (name, capacity, location, description, 
			  availability_start, availability_end, is_available) 
			  VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		room.Name, 
		room.Capacity, 
		room.Location, 
		room.Description, 
		room.AvailabilityStart,
		room.AvailabilityEnd,
		room.IsAvailable,
	)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	room.ID = int(id)
	return nil
}

// Update modifies an existing room
func (r *RoomRepository) Update(id int, room *Room) error {
	query := `UPDATE rooms SET name = ?, capacity = ?, location = ?, description = ?,
			  availability_start = ?, availability_end = ?, is_available = ?
			  WHERE id = ?`
	
	_, err := r.DB.Exec(
		query, 
		room.Name, 
		room.Capacity, 
		room.Location, 
		room.Description, 
		room.AvailabilityStart,
		room.AvailabilityEnd,
		room.IsAvailable,
		id,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

// Delete removes a room from the database
func (r *RoomRepository) Delete(id int) error {
	// First check if this room has any reservations
	checkQuery := `SELECT COUNT(*) FROM reservations WHERE room_id = ?`
	var count int
	err := r.DB.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return errors.New("cannot delete room: it has existing reservations")
	}
	
	// If no reservations, proceed with deletion
	query := `DELETE FROM rooms WHERE id = ?`
	_, err = r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	return nil
}

// CheckAvailability checks if a room is available for a given time slot
func (r *RoomRepository) CheckAvailability(roomID int, start, end time.Time) (bool, error) {
	// First check if the room exists and is available for booking
	roomQuery := `SELECT is_available FROM rooms WHERE id = ?`
	var isAvailable bool
	err := r.DB.QueryRow(roomQuery, roomID).Scan(&isAvailable)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("room not found")
		}
		return false, err
	}
	
	if !isAvailable {
		return false, nil
	}
	
	// Check if there are any overlapping reservations
	query := `
		SELECT COUNT(*) FROM reservations 
		WHERE room_id = ? AND (
			(start_time <= ? AND end_time > ?) OR  -- Reservation starts before and ends during/after
			(start_time < ? AND end_time >= ?) OR  -- Reservation starts during and ends after
			(start_time >= ? AND end_time <= ?)    -- Reservation is completely within
		)
	`
	
	var count int
	err = r.DB.QueryRow(query, roomID, end, start, end, start, start, end).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count == 0, nil
}

// GetAvailableRooms returns rooms available for a given time slot
func (r *RoomRepository) GetAvailableRooms(start, end time.Time, minCapacity int) ([]*Room, error) {
	// Get all rooms that are marked as available
	query := `
		SELECT id, name, capacity, location, description, 
		availability_start, availability_end, is_available, created_at, updated_at 
		FROM rooms 
		WHERE is_available = TRUE AND capacity >= ?
	`
	
	rows, err := r.DB.Query(query, minCapacity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var availableRooms []*Room
	
	for rows.Next() {
		room := &Room{}
		var availabilityStart, availabilityEnd sql.NullString
		
		err := rows.Scan(
			&room.ID, 
			&room.Name, 
			&room.Capacity, 
			&room.Location, 
			&room.Description, 
			&availabilityStart,
			&availabilityEnd,
			&room.IsAvailable,
			&room.CreatedAt, 
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if availabilityStart.Valid {
			room.AvailabilityStart = availabilityStart.String
		}
		if availabilityEnd.Valid {
			room.AvailabilityEnd = availabilityEnd.String
		}
		
		// Check if room has conflicting reservations for the requested time
		isAvailable, err := r.CheckAvailability(room.ID, start, end)
		if err != nil {
			return nil, err
		}
		
		if isAvailable {
			availableRooms = append(availableRooms, room)
		}
	}
	
	return availableRooms, nil
}