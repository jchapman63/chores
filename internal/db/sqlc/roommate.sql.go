// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: roommate.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getRoommateByID = `-- name: GetRoommateByID :one
SELECT id, name, email, chore, last_updated FROM roommates
WHERE id = $1
`

// Get a roommate by ID
func (q *Queries) GetRoommateByID(ctx context.Context, id pgtype.UUID) (Roommate, error) {
	row := q.db.QueryRow(ctx, getRoommateByID, id)
	var i Roommate
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Chore,
		&i.LastUpdated,
	)
	return i, err
}

const getRoommates = `-- name: GetRoommates :many
SELECT id, name, email, chore, last_updated FROM roommates
ORDER BY name, email
`

// Get all roommates sorted by name and email
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
			&i.Email,
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
RETURNING id, name, email, chore, last_updated
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
		&i.Email,
		&i.Chore,
		&i.LastUpdated,
	)
	return i, err
}

const upsertRoommate = `-- name: UpsertRoommate :one
INSERT INTO roommates (name, email, chore)
VALUES ($1, $2, $3)
ON CONFLICT (name, email) DO NOTHING
RETURNING id, name, email, chore, last_updated
`

type UpsertRoommateParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Chore string `json:"chore"`
}

// Insert or update a roommate record
func (q *Queries) UpsertRoommate(ctx context.Context, arg UpsertRoommateParams) (Roommate, error) {
	row := q.db.QueryRow(ctx, upsertRoommate, arg.Name, arg.Email, arg.Chore)
	var i Roommate
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Chore,
		&i.LastUpdated,
	)
	return i, err
}
