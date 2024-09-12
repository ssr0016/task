
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add a department_id column to the users table if it does not already exist. This column will reference the id of the department:
ALTER TABLE users
ADD COLUMN department_id INTEGER REFERENCES departments(id);