package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Connect establishes a connection to the database
func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to database successfully")
	return db, nil
}

// EnsureTablesExist creates the necessary tables if they don't exist
func EnsureTablesExist(db *sql.DB) error {
	// Create rooms table if it doesn't exist
	roomsTable := `
		CREATE TABLE IF NOT EXISTS rooms (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			capacity INT NOT NULL,
			location VARCHAR(255),
			description TEXT,
			availability_start TIME,
			availability_end TIME,
			is_available BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY (name)
		) ENGINE=InnoDB;
	`
	_, err := db.Exec(roomsTable)
	if err != nil {
		return err
	}

	// Create reservations table if it doesn't exist
	reservationsTable := `
		CREATE TABLE IF NOT EXISTS reservations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			room_id INT NOT NULL,
			user_id INT NOT NULL,
			contact_id INT,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			start_time DATETIME NOT NULL,
			end_time DATETIME NOT NULL,
			recurring_type ENUM('none', 'daily', 'weekly', 'monthly') DEFAULT 'none',
			recurring_end_date DATE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (room_id) REFERENCES rooms(id),
			UNIQUE KEY (room_id, start_time, end_time)
		) ENGINE=InnoDB;
	`
	_, err = db.Exec(reservationsTable)
	if err != nil {
		return err
	}

	// Insert default rooms if none exist
	var roomCount int
	err = db.QueryRow("SELECT COUNT(*) FROM rooms").Scan(&roomCount)
	if err != nil {
		return err
	}

	if roomCount == 0 {
		defaultRooms := `
			INSERT INTO rooms (name, capacity, location, description, availability_start, availability_end) VALUES
			('Main Sanctuary', 200, 'Main Building', 'The main church sanctuary for worship services', '08:00:00', '21:00:00'),
			('Fellowship Hall', 100, 'Basement', 'Large open space for events and gatherings', '08:00:00', '22:00:00'),
			('Classroom A', 30, 'Education Wing', 'Classroom with tables and chairs', '08:00:00', '21:00:00'),
			('Classroom B', 30, 'Education Wing', 'Classroom with tables and chairs', '08:00:00', '21:00:00'),
			('Conference Room', 15, 'Office Area', 'Conference room with large table', '08:00:00', '20:00:00'),
			('Youth Room', 50, 'West Wing', 'Activity space designed for youth ministry', '08:00:00', '22:00:00'),
			('Prayer Chapel', 20, 'East Wing', 'Small chapel for prayer gatherings', '07:00:00', '23:00:00'),
			('Kitchen', 10, 'Near Fellowship Hall', 'Fully equipped kitchen for event preparation', '08:00:00', '21:00:00'),
			('Nursery', 15, 'Near Main Sanctuary', 'Childcare area for infants and toddlers', '08:00:00', '13:00:00'),
			('Choir Room', 35, 'Near Main Sanctuary', 'Practice space for the choir', '16:00:00', '21:00:00');
		`
		_, err = db.Exec(defaultRooms)
		if err != nil {
			log.Println("Warning: Failed to insert some or all default rooms: ", err)
			log.Println("This may be normal if some rooms already exist")
		} else {
			log.Println("Default rooms created")
		}
	}

	log.Println("Database tables ready")
	return nil
}