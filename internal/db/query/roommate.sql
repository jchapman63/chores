-- name: GetRoommates :many
-- Get all roommates sorted by ID
SELECT * FROM roommates
ORDER BY id;

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
