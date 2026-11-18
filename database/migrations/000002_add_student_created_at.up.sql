
-- Add a creation timestamp to the students table
ALTER TABLE students
ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
