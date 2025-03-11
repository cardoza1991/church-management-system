package models

import (
	"database/sql"
	"errors"
)

// Status represents a stage in the contact's spiritual journey
type Status struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DisplayOrder int   `json:"display_order"`
}

// StatusRepository provides access to the status store
type StatusRepository struct {
	DB *sql.DB
}

// NewStatusRepository creates a new StatusRepository
func NewStatusRepository(db *sql.DB) *StatusRepository {
	return &StatusRepository{DB: db}
}

// GetAll retrieves all statuses ordered by display_order
func (r *StatusRepository) GetAll() ([]*Status, error) {
	query := `SELECT id, name, description, display_order FROM statuses ORDER BY display_order`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var statuses []*Status
	for rows.Next() {
		status := &Status{}
		err := rows.Scan(
			&status.ID, 
			&status.Name, 
			&status.Description, 
			&status.DisplayOrder,
		)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	
	return statuses, nil
}

// GetByID retrieves a status by ID
func (r *StatusRepository) GetByID(id int) (*Status, error) {
	query := `SELECT id, name, description, display_order FROM statuses WHERE id = ?`
	
	status := &Status{}
	err := r.DB.QueryRow(query, id).Scan(
		&status.ID, 
		&status.Name, 
		&status.Description, 
		&status.DisplayOrder,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("status not found")
		}
		return nil, err
	}
	
	return status, nil
}

// Create adds a new status to the database
func (r *StatusRepository) Create(status *Status) error {
	query := `INSERT INTO statuses (name, description, display_order) VALUES (?, ?, ?)`
	
	result, err := r.DB.Exec(query, status.Name, status.Description, status.DisplayOrder)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	status.ID = int(id)
	return nil
}

// Update modifies an existing status
func (r *StatusRepository) Update(id int, status *Status) error {
	query := `UPDATE statuses SET name = ?, description = ?, display_order = ? WHERE id = ?`
	
	_, err := r.DB.Exec(query, status.Name, status.Description, status.DisplayOrder, id)
	if err != nil {
		return err
	}
	
	return nil
}

// Delete removes a status from the database
func (r *StatusRepository) Delete(id int) error {
	// First check if this status is being used by any contacts
	checkQuery := `SELECT COUNT(*) FROM contacts WHERE current_status_id = ?`
	var count int
	err := r.DB.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return errors.New("cannot delete status: it is being used by contacts")
	}
	
	// Also check status history
	checkHistoryQuery := `SELECT COUNT(*) FROM contact_status_history WHERE status_id = ?`
	err = r.DB.QueryRow(checkHistoryQuery, id).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return errors.New("cannot delete status: it is used in contact history")
	}
	
	// If not used, proceed with deletion
	query := `DELETE FROM statuses WHERE id = ?`
	_, err = r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	return nil
}