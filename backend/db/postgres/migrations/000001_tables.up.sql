CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    google_id VARCHAR(100) UNIQUE DEFAULT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) DEFAULT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    picture TEXT DEFAULT NULL,
    storage_used BIGINT DEFAULT 0,
    max_storage BIGINT DEFAULT 100000000,
    -- 100 MB in bytes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Create the Metadata table
CREATE TABLE IF NOT EXISTS metadata (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    file_location TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT,
    file_type VARCHAR(50),
    file_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- Create the Keywords table
CREATE TABLE IF NOT EXISTS keywords (
    id SERIAL PRIMARY KEY,
    keyword TEXT NOT NULL,
    metadata_id INT NOT NULL,
    FOREIGN KEY (metadata_id) REFERENCES metadata(id) ON DELETE CASCADE
);