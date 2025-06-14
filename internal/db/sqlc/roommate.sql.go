// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: roommate.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getRoommateByID = `-- name: GetRoommateByID :one
SELECT id, name, phone_number, chore, last_updated FROM roommates
WHERE id = $1
`

// Get a roommate by ID
func (q *Queries) GetRoommateByID(ctx context.Context, id pgtype.UUID) (Roommate, error) {
	row := q.db.QueryRow(ctx, getRoommateByID, id)
	var i Roommate
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PhoneNumber,
		&i.Chore,
		&i.LastUpdated,
	)
	return i, err
}

const getRoommates = `-- name: GetRoommates :many
SELECT id, name, phone_number, chore, last_updated FROM roommates
ORDER BY id
`

// Get all roommates sorted by ID
func (q *Queries) GetRoommates(ctx context.Context) ([]Roommate, error) {
	rows, err := q.db.Query(ctx, getRoommates)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Roommate
	for rows.Next() {
		var i Roommate
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PhoneNumber,
			&i.Chore,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateRoommateChore = `-- name: UpdateRoommateChore :one
UPDATE roommates
SET
    chore = $2,
    last_updated = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, name, phone_number, chore, last_updated
`

type UpdateRoommateChoreParams struct {
	ID    pgtype.UUID `json:"id"`
	Chore string      `json:"chore"`
}

// Update a roommate's chore
func (q *Queries) UpdateRoommateChore(ctx context.Context, arg UpdateRoommateChoreParams) (Roommate, error) {
	row := q.db.QueryRow(ctx, updateRoommateChore, arg.ID, arg.Chore)
	var i Roommate
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PhoneNumber,
		&i.Chore,
		&i.LastUpdated,
	)
	return i, err
}
