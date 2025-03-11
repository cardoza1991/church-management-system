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
	// Create statuses table if it doesn't exist
	statusesTable := `
		CREATE TABLE IF NOT EXISTS statuses (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT,
			display_order INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY (name)
		) ENGINE=InnoDB;
	`
	_, err := db.Exec(statusesTable)
	if err != nil {
		return err
	}

	// Create contacts table if it doesn't exist
	contactsTable := `
		CREATE TABLE IF NOT EXISTS contacts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			phone VARCHAR(50),
			location VARCHAR(255),
			notes TEXT,
			current_status_id INT NOT NULL,
			date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (current_status_id) REFERENCES statuses(id)
		) ENGINE=InnoDB;
	`
	_, err = db.Exec(contactsTable)
	if err != nil {
		return err
	}

	// Create contact status history table if it doesn't exist
	historyTable := `
		CREATE TABLE IF NOT EXISTS contact_status_history (
			id INT AUTO_INCREMENT PRIMARY KEY,
			contact_id INT NOT NULL,
			status_id INT NOT NULL,
			notes TEXT,
			date_changed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE,
			FOREIGN KEY (status_id) REFERENCES statuses(id)
		) ENGINE=InnoDB;
	`
	_, err = db.Exec(historyTable)
	if err != nil {
		return err
	}

	// Insert default statuses if none exist
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM statuses").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		defaultStatuses := `
			INSERT INTO statuses (name, description, display_order) VALUES
			('New Contact', 'Initial contact with the church', 1),
			('In Studies', 'Actively participating in Bible studies', 2),
			('Baptized', 'Has been baptized', 3),
			('Gospel Worker', 'Actively sharing faith with others', 4);
		`
		_, err = db.Exec(defaultStatuses)
		if err != nil {
			return err
		}
		log.Println("Default statuses created")
	}

	log.Println("Database tables ready")
	return nil
}