package models

import (
	"database/sql"
	"errors"
	"time"
)

// Lesson represents a Bible study lesson in the curriculum
type Lesson struct {
	ID             int       `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description,omitempty"`
	SequenceNumber int       `json:"sequence_number"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// LessonRepository provides access to the lesson store
type LessonRepository struct {
	DB *sql.DB
}

// NewLessonRepository creates a new LessonRepository
func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{DB: db}
}

// GetAll retrieves all lessons
func (r *LessonRepository) GetAll() ([]*Lesson, error) {
	query := `SELECT id, title, description, sequence_number, created_at, updated_at 
			  FROM lessons ORDER BY sequence_number`
	
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var lessons []*Lesson
	for rows.Next() {
		lesson := &Lesson{}
		err := rows.Scan(
			&lesson.ID, 
			&lesson.Title, 
			&lesson.Description, 
			&lesson.SequenceNumber, 
			&lesson.CreatedAt, 
			&lesson.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}
	
	return lessons, nil
}

// GetByID retrieves a lesson by ID
func (r *LessonRepository) GetByID(id int) (*Lesson, error) {
	query := `SELECT id, title, description, sequence_number, created_at, updated_at 
			  FROM lessons WHERE id = ?`
	
	lesson := &Lesson{}
	err := r.DB.QueryRow(query, id).Scan(
		&lesson.ID, 
		&lesson.Title, 
		&lesson.Description, 
		&lesson.SequenceNumber, 
		&lesson.CreatedAt, 
		&lesson.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("lesson not found")
		}
		return nil, err
	}
	
	return lesson, nil
}

// Create adds a new lesson to the database
func (r *LessonRepository) Create(lesson *Lesson) error {
	query := `INSERT INTO lessons (title, description, sequence_number) VALUES (?, ?, ?)`
	
	result, err := r.DB.Exec(query, lesson.Title, lesson.Description, lesson.SequenceNumber)
	if err != nil {
		return err
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	
	lesson.ID = int(id)
	return nil
}

// Update modifies an existing lesson
func (r *LessonRepository) Update(id int, lesson *Lesson) error {
	query := `UPDATE lessons SET title = ?, description = ?, sequence_number = ? WHERE id = ?`
	
	_, err := r.DB.Exec(query, lesson.Title, lesson.Description, lesson.SequenceNumber, id)
	if err != nil {
		return err
	}
	
	return nil
}

// Delete removes a lesson from the database
func (r *LessonRepository) Delete(id int) error {
	// First check if this lesson is being used by any studies
	checkQuery := `SELECT COUNT(*) FROM studies WHERE lesson_id = ?`
	var count int
	err := r.DB.QueryRow(checkQuery, id).Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return errors.New("cannot delete lesson: it is being used in study records")
	}
	
	// If not used, proceed with deletion
	query := `DELETE FROM lessons WHERE id = ?`
	_, err = r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	return nil
}