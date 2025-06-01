-- Schema for the chore rotation app
-- This creates a Roommate table with name, phone number, and chore fields
CREATE TABLE IF NOT EXISTS roommates (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    chore TEXT NOT NULL CHECK (chore IN ('BATHROOM', 'FLOOR', 'COUNTER')),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
