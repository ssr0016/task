CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    status INT NOT NULL DEFAULT 1, -- Using the TaskStatus enum: 1=Pending, 2=In Progress, 3=Done
    user_id INT NOT NULL, -- Foreign key to the users table (assuming users table exists)
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
        REFERENCES users(id) ON DELETE CASCADE
);

-- -- Step 1: Define ENUM types for priority and difficulty
-- CREATE TYPE task_priority AS ENUM ('low', 'medium', 'high');
-- CREATE TYPE task_difficulty AS ENUM ('easy', 'medium', 'hard');


-- -- Step 2: Modify the tasks table to include these fields
-- ALTER TABLE tasks
-- ADD COLUMN priority task_priority NOT NULL DEFAULT 'medium',
-- ADD COLUMN difficulty task_difficulty NOT NULL DEFAULT 'medium';