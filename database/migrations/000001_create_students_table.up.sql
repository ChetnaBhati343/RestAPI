
-- Create the students table
CREATE TABLE students (
    id INTEGER GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(100) NOT NULL,
    age INTEGER NOT NULL,
    email VARCHAR(255) NOT NULL,

    -- Define table constraints
    CONSTRAINT students_pkey PRIMARY KEY (id),
    CONSTRAINT students_email_unique UNIQUE (email),
    CONSTRAINT students_age_positive CHECK (age > 0)
);
