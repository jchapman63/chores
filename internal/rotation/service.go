package rotation

import (
	"context"

	db "github.com/jchapman63/chores/internal/db/sqlc"
)

// ChoreType represents the type of chore
type ChoreType string

const (
	// ChoreTypeBathroom represents bathroom cleaning duty
	ChoreTypeBathroom ChoreType = "BATHROOM"
	// ChoreTypeFloor represents floor cleaning duty
	ChoreTypeFloor ChoreType = "FLOOR"
	// ChoreTypeCounter represents counter cleaning duty
	ChoreTypeCounter ChoreType = "COUNTER"
)

type RotationService struct {
	dbQueries db.Querier
}

func NewService(dbQueries db.Querier) *RotationService {
	return &RotationService{
		dbQueries: dbQueries,
	}
}

// NextChore returns the next chore in the rotation
func NextChore(current ChoreType) ChoreType {
	switch current {
	case ChoreTypeBathroom:
		return ChoreTypeFloor
	case ChoreTypeFloor:
		return ChoreTypeCounter
	case ChoreTypeCounter:
		return ChoreTypeBathroom
	default:
		// Default to bathroom if an invalid chore is provided
		return ChoreTypeBathroom
	}
}

// RotateChores rotates all roommates' chores to the next in the sequence
func (s *RotationService) RotateChores(ctx context.Context) error {
	// Get all roommates
	roommates, err := s.dbQueries.GetRoommates(ctx)
	if err != nil {
		return err
	}

	// Update each roommate's chore to the next one
	for _, roommate := range roommates {
		currentChore := ChoreType(roommate.Chore)
		nextChore := NextChore(currentChore)

		_, err := s.dbQueries.UpdateRoommateChore(ctx, db.UpdateRoommateChoreParams{
			ID:    roommate.ID,
			Chore: string(nextChore),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
