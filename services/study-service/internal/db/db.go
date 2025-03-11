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
	// Create predefined lessons table if it doesn't exist
	lessonsTable := `
		CREATE TABLE IF NOT EXISTS lessons (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			sequence_number INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY (title),
			UNIQUE KEY (sequence_number)
		) ENGINE=InnoDB;
	`
	_, err := db.Exec(lessonsTable)
	if err != nil {
		return err
	}

	// Create study sessions table if it doesn't exist
	studiesTable := `
		CREATE TABLE IF NOT EXISTS studies (
			id INT AUTO_INCREMENT PRIMARY KEY,
			contact_id INT NOT NULL,
			lesson_id INT NOT NULL,
			date_completed DATE NOT NULL,
			location VARCHAR(255),
			duration_minutes INT,
			notes TEXT,
			taught_by_user_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY (contact_id, lesson_id, date_completed),
			FOREIGN KEY (lesson_id) REFERENCES lessons(id)
		) ENGINE=InnoDB;
	`
	_, err = db.Exec(studiesTable)
	if err != nil {
		return err
	}

	// Insert standard Bible study lessons if none exist
	var lessonCount int
	err = db.QueryRow("SELECT COUNT(*) FROM lessons").Scan(&lessonCount)
	if err != nil {
		return err
	}

	if lessonCount == 0 {
		defaultLessons := `
			INSERT INTO lessons (title, description, sequence_number) VALUES
			('The Word of God', 'Introduction to the Bible as God\'s inspired word', 1),
			('The Gospel', 'The good news of salvation through Jesus Christ', 2),
			('Conversion', 'The process of becoming a disciple of Jesus', 3),
			('Discipleship', 'What it means to follow Jesus daily', 4),
			('The Church', 'God\'s family and kingdom on earth', 5),
			('Prayer', 'Communicating with God through prayer', 6),
			('Baptism', 'The meaning and importance of baptism', 7),
			('The Holy Spirit', 'The gift and work of the Holy Spirit', 8),
			('Spiritual Disciplines', 'Practices that help us grow spiritually', 9),
			('Evangelism', 'Sharing your faith with others', 10),
			('Spiritual Gifts', 'Using your gifts to serve the church', 11),
			('Christian Character', 'Growing in Christ-like character', 12),
			('The Kingdom of God', 'Understanding God\'s kingdom', 13),
			('Spiritual Warfare', 'Standing firm against spiritual opposition', 14),
			('Biblical Leadership', 'Principles of godly leadership', 15),
			('God\'s Plan for Marriage', 'Biblical view of marriage and family', 16),
			('Financial Stewardship', 'Managing finances God\'s way', 17),
			('Living in Community', 'The importance of Christian fellowship', 18),
			('Sharing Your Testimony', 'How to effectively share your story', 19),
			('The Great Commission', 'Our call to make disciples of all nations', 20),
			('The Return of Christ', 'The second coming and end times', 21),
			('Spiritual Multiplication', 'Discipling others who disciple others', 22),
			('World Missions', 'God\'s heart for all nations', 23),
			('Biblical Worldview', 'Seeing all of life through God\'s perspective', 24),
			('Apologetics', 'Defending the Christian faith', 25),
			('God\'s Purpose for Work', 'Integrating faith and work', 26),
			('Living by Faith', 'Trusting God in all circumstances', 27),
			('The Life of Christ', 'Key events and teachings from Jesus\' life', 28),
			('The Cross', 'The centrality and significance of the crucifixion', 29),
			('Servanthood', 'Following Christ\'s example of serving others', 30);
		`
		_, err = db.Exec(defaultLessons)
		if err != nil {
			// It's possible that the INSERT fails because of duplicate keys
			// For example, if there's a partial insert from a previous attempt
			log.Println("Warning: Failed to insert some or all default lessons: ", err)
			log.Println("This may be normal if some lessons already exist")
		} else {
			log.Println("Default Bible study lessons created")
		}
	}

	log.Println("Database tables ready")
	return nil
}