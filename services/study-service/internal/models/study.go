package models

import (
	"database/sql"
	"errors"
	"time"
)

// Study represents a completed Bible study session with a contact
type Study struct {
	ID              int       `json:"id"`
	ContactID       int       `json:"contact_id"`
	LessonID        int       `json:"lesson_id"`
	LessonTitle     string    `json:"lesson_title,omitempty"`
	DateCompleted   time.Time `json:"date_completed"`
	Location        string    `json:"location,omitempty"`
	DurationMinutes int       `json:"duration_minutes,omitempty"`
	Notes           string    `json:"notes,omitempty"`
	TaughtByUserID  int       `json:"taught_by_user_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// StudyRepository provides access to the study store
type StudyRepository struct {
	DB *sql.DB
}

// NewStudyRepository creates a new StudyRepository
func NewStudyRepository(db *sql.DB) *StudyRepository {
	return &StudyRepository{DB: db}
}

// GetByContactID retrieves all studies for a specific contact
func (r *StudyRepository) GetByContactID(contactID int) ([]*Study, error) {
	query := `
		SELECT s.id, s.contact_id, s.lesson_id, l.title, s.date_completed, 
			   s.location, s.duration_minutes, s.notes, s.taught_by_user_id, 
			   s.created_at, s.updated_at
		FROM studies s
		JOIN lessons l ON s.lesson_id = l.id
		WHERE s.contact_id = ?
		ORDER BY s.date_completed DESC
	`
	
	rows, err := r.DB.Query(query, contactID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var studies []*Study
	for rows.Next() {
		study := &Study{}
		err := rows.Scan(
			&study.ID, 
			&study.ContactID, 
			&study.LessonID, 
			&study.LessonTitle, 
			&study.DateCompleted, 
			&study.Location, 
			&study.DurationMinutes, 
			&study.Notes, 
			&study.TaughtByUserID, 
			&study.CreatedAt, 
			&study.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		studies = append(studies, study)
	}
	
	return studies, nil
}

// GetByID retrieves a study by ID
func (r *StudyRepository) GetByID(id int) (*Study, error) {
	query := `
		SELECT s.id, s.contact_id, s.lesson_id, l.title, s.date_completed, 
			   s.location, s.duration_minutes, s.notes, s.taught_by_user_id, 
			   s.created_at, s.updated_at
		FROM studies s
		JOIN lessons l ON s.lesson_id = l.id
		WHERE s.id = ?
	`
	
	study := &Study{}
	err := r.DB.QueryRow(query, id).Scan(
		&study.ID, 
		&study.ContactID, 
		&study.LessonID, 
		&study.LessonTitle, 
		&study.DateCompleted, 
		&study.Location, 
		&study.DurationMinutes, 
		&study.Notes, 
		&study.TaughtByUserID, 
		&study.CreatedAt, 
		&study.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("study not found")
		}
		return nil, err
	}
	
	return study, nil
}

// Create adds a new study to the database
func (r *StudyRepository) Create(study *Study) error {
	query := `
		INSERT INTO studies 
		(contact_id, lesson_id, date_completed, location, duration_minutes, notes, taught_by_user_id)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.DB.Exec(
		query, 
		study.ContactID, 
		study.LessonID, 
		study.DateCompleted, 
		study.Location, 
		study.DurationMinutes, 
		study.Notes, 
		study.TaughtByUserID,
	)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	study.ID = int(id)
	return nil
}

// Update modifies an existing study
func (r *StudyRepository) Update(id int, study *Study) error {
	query := `
		UPDATE studies 
		SET contact_id = ?, lesson_id = ?, date_completed = ?, 
			location = ?, duration_minutes = ?, notes = ?, taught_by_user_id = ?
		WHERE id = ?
	`
	
	_, err := r.DB.Exec(
		query, 
		study.ContactID, 
		study.LessonID, 
		study.DateCompleted, 
		study.Location, 
		study.DurationMinutes, 
		study.Notes, 
		study.TaughtByUserID,
		id,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

// Delete removes a study from the database
func (r *StudyRepository) Delete(id int) error {
	query := `DELETE FROM studies WHERE id = ?`
	
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	return nil
}

// GetCompletedLessonsByContactID returns a list of lesson IDs completed by a contact
func (r *StudyRepository) GetCompletedLessonsByContactID(contactID int) (map[int]bool, error) {
	query := `SELECT DISTINCT lesson_id FROM studies WHERE contact_id = ?`
	
	rows, err := r.DB.Query(query, contactID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	completedLessons := make(map[int]bool)
	for rows.Next() {
		var lessonID int
		if err := rows.Scan(&lessonID); err != nil {
			return nil, err
		}
		completedLessons[lessonID] = true
	}
	
	return completedLessons, nil
}

// GetStudyStats returns statistics about studies for a contact
type StudyStats struct {
	TotalLessons        int       `json:"total_lessons"`
	CompletedLessons    int       `json:"completed_lessons"`
	ProgressPercentage  float64   `json:"progress_percentage"`
	LastStudyDate       time.Time `json:"last_study_date,omitempty"`
	TotalStudyTimeMinutes int     `json:"total_study_time_minutes"`
}

// GetContactStudyStats returns study statistics for a given contact
func (r *StudyRepository) GetContactStudyStats(contactID int) (*StudyStats, error) {
	// Get total lessons count
	var totalLessons int
	err := r.DB.QueryRow("SELECT COUNT(*) FROM lessons").Scan(&totalLessons)
	if err != nil {
		return nil, err
	}
	
	// Get completed lessons count
	var completedLessons int
	err = r.DB.QueryRow("SELECT COUNT(DISTINCT lesson_id) FROM studies WHERE contact_id = ?", contactID).Scan(&completedLessons)
	if err != nil {
		return nil, err
	}
	
	// Calculate progress percentage
	var progressPercentage float64
	if totalLessons > 0 {
		progressPercentage = float64(completedLessons) / float64(totalLessons) * 100
	}
	
	// Get last study date
	var lastStudyDate sql.NullTime
	err = r.DB.QueryRow("SELECT MAX(date_completed) FROM studies WHERE contact_id = ?", contactID).Scan(&lastStudyDate)
	if err != nil {
		return nil, err
	}
	
	// Get total study time in minutes
	var totalStudyTimeMinutes int
	err = r.DB.QueryRow("SELECT COALESCE(SUM(duration_minutes), 0) FROM studies WHERE contact_id = ? AND duration_minutes IS NOT NULL", contactID).Scan(&totalStudyTimeMinutes)
	if err != nil {
		return nil, err
	}
	
	stats := &StudyStats{
		TotalLessons:         totalLessons,
		CompletedLessons:     completedLessons,
		ProgressPercentage:   progressPercentage,
		TotalStudyTimeMinutes: totalStudyTimeMinutes,
	}
	
	if lastStudyDate.Valid {
		stats.LastStudyDate = lastStudyDate.Time
	}
	
	return stats, nil
}