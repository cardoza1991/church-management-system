-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    role ENUM('admin', 'overseer', 'group_leader', 'teacher', 'member') NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Contacts table
CREATE TABLE IF NOT EXISTS contacts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    location VARCHAR(255),
    phone VARCHAR(20),
    email VARCHAR(100),
    date_added TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    current_status_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Contact statuses table
CREATE TABLE IF NOT EXISTS contact_statuses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    order_index INT NOT NULL,
    is_milestone BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Add foreign key after tables are created
ALTER TABLE contacts
ADD CONSTRAINT fk_contact_status
FOREIGN KEY (current_status_id) REFERENCES contact_statuses(id);

-- Insert initial status values
INSERT INTO contact_statuses (name, description, order_index, is_milestone)
VALUES 
('New Contact', 'Initial contact with no studies yet', 1, FALSE),
('In Studies', 'Currently receiving Bible studies', 2, FALSE),
('Baptized', 'Has been baptized', 3, TRUE),
('Gospel Worker', 'Has completed 30 studies and is ready to teach others', 4, TRUE);
