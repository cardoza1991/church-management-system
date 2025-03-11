package models

import (
	"database/sql"
	"errors"
	"time"
)

// Contact represents a contact in the system
type Contact struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email,omitempty"`
	Phone           string    `json:"phone,omitempty"`
	Location        string    `json:"location,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	DateAdded       time.Time `json:"date_added"`
	LastUpdated     time.Time `json:"last_updated"`
	CurrentStatusID int       `json:"current_status_id"`
}

// ContactRepository provides access to the contact store
type ContactRepository struct {
	DB *sql.DB
}

// NewContactRepository creates a new ContactRepository
func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{DB: db}
}

// GetAll retrieves all contacts
func (r *ContactRepository) GetAll(limit, offset int) ([]*Contact, error) {
	query := `SELECT id, name, email, phone, location, notes, date_added, 
	          last_updated, current_status_id 
	          FROM contacts ORDER BY name LIMIT ? OFFSET ?`
	
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var contacts []*Contact
	for rows.Next() {
		contact := &Contact{}
		err := rows.Scan(
			&contact.ID, 
			&contact.Name, 
			&contact.Email, 
			&contact.Phone, 
			&contact.Location, 
			&contact.Notes, 
			&contact.DateAdded, 
			&contact.LastUpdated, 
			&contact.CurrentStatusID,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	
	return contacts, nil
}

// GetByID retrieves a contact by ID
func (r *ContactRepository) GetByID(id int) (*Contact, error) {
	query := `SELECT id, name, email, phone, location, notes, date_added, 
	          last_updated, current_status_id 
	          FROM contacts WHERE id = ?`
	
	contact := &Contact{}
	err := r.DB.QueryRow(query, id).Scan(
		&contact.ID, 
		&contact.Name, 
		&contact.Email, 
		&contact.Phone, 
		&contact.Location, 
		&contact.Notes, 
		&contact.DateAdded, 
		&contact.LastUpdated, 
		&contact.CurrentStatusID,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("contact not found")
		}
		return nil, err
	}
	
	return contact, nil
}

// Create adds a new contact to the database
func (r *ContactRepository) Create(contact *Contact) error {
	query := `INSERT INTO contacts (name, email, phone, location, notes, current_status_id)
	          VALUES (?, ?, ?, ?, ?, ?)`
	
	result, err := r.DB.Exec(
		query, 
		contact.Name, 
		contact.Email, 
		contact.Phone, 
		contact.Location, 
		contact.Notes, 
		contact.CurrentStatusID,
	)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	contact.ID = int(id)
	return nil
}

// Update modifies an existing contact
func (r *ContactRepository) Update(id int, contact *Contact) error {
	query := `UPDATE contacts 
	          SET name = ?, email = ?, phone = ?, location = ?, notes = ?, 
	          current_status_id = ?, last_updated = NOW()
	          WHERE id = ?`
	
	_, err := r.DB.Exec(
		query, 
		contact.Name, 
		contact.Email, 
		contact.Phone, 
		contact.Location, 
		contact.Notes, 
		contact.CurrentStatusID,
		id,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

// Delete removes a contact from the database
func (r *ContactRepository) Delete(id int) error {
	query := `DELETE FROM contacts WHERE id = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	return nil
}

// UpdateStatus changes a contact's status and logs the change
func (r *ContactRepository) UpdateStatus(contactID, statusID int, notes string) error {
	// Start a transaction to ensure both operations succeed or fail together
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	
	// Update the contact's current status
	updateQuery := `UPDATE contacts SET current_status_id = ?, last_updated = NOW() WHERE id = ?`
	_, err = tx.Exec(updateQuery, statusID, contactID)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// Log the status change in the history table
	historyQuery := `INSERT INTO contact_status_history (contact_id, status_id, notes)
	                 VALUES (?, ?, ?)`
	_, err = tx.Exec(historyQuery, contactID, statusID, notes)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	// Commit the transaction
	return tx.Commit()
}

// StatusHistoryEntry represents a status change for a contact
type StatusHistoryEntry struct {
	ID          int       `json:"id"`
	ContactID   int       `json:"contact_id"`
	StatusID    int       `json:"status_id"`
	StatusName  string    `json:"status_name"`
	Notes       string    `json:"notes,omitempty"`
	DateChanged time.Time `json:"date_changed"`
}

// GetStatusHistory retrieves the status history for a contact
func (r *ContactRepository) GetStatusHistory(contactID int) ([]*StatusHistoryEntry, error) {
	query := `SELECT h.id, h.contact_id, h.status_id, s.name as status_name, h.notes, h.date_changed
	          FROM contact_status_history h
	          JOIN statuses s ON h.status_id = s.id
	          WHERE h.contact_id = ?
	          ORDER BY h.date_changed DESC`
	
	rows, err := r.DB.Query(query, contactID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var history []*StatusHistoryEntry
	for rows.Next() {
		entry := &StatusHistoryEntry{}
		err := rows.Scan(
			&entry.ID,
			&entry.ContactID,
			&entry.StatusID,
			&entry.StatusName,
			&entry.Notes,
			&entry.DateChanged,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, entry)
	}
	
	return history, nil
}