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
