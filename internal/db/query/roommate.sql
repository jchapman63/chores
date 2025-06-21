-- name: GetRoommates :many
-- Get all roommates sorted by name and email
SELECT * FROM roommates
ORDER BY name, email;

-- name: GetRoommateByID :one
-- Get a roommate by ID
SELECT * FROM roommates
WHERE id = $1;

-- name: UpdateRoommateChore :one
-- Update a roommate's chore
UPDATE roommates
SET
    chore = $2,
    last_updated = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpsertRoommate :one
-- Insert or update a roommate record
INSERT INTO roommates (name, email, chore)
VALUES ($1, $2, $3)
ON CONFLICT (name, email) DO NOTHING
RETURNING *;