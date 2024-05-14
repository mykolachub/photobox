CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    storage_used BIGINT DEFAULT 0,
    max_storage BIGINT DEFAULT 100000000,
    -- 100 MB in bytes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Create the Metadata table
CREATE TABLE IF NOT EXISTS metadata (
    file_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    file_location TEXT NOT NULL,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT,
    file_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);
-- Create the Keywords table
CREATE TABLE IF NOT EXISTS keywords (
    keyword_id SERIAL PRIMARY KEY,
    keyword TEXT NOT NULL,
    metadata_id INT NOT NULL,
    FOREIGN KEY (metadata_id) REFERENCES metadata(file_id) ON DELETE CASCADE
);