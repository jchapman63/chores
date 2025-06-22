package rotation

import (
	"context"

	"github.com/jackc/pgx/v5"
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
func nextChore(current ChoreType) ChoreType {
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

// Starts the chore rotation based on the current state of roommates in the real world
func (s *RotationService) InitializeChores(ctx context.Context) error {
	roommates := []struct {
		name  string
		email string
		chore string
	}{
		{"Jason", "jasonschmitz@icloud.com", "BATHROOM"},
		{"Daniel", "danielhiser17@gmail.com", "FLOOR"},
		{"Jordan", "jordanchap20@gmail.com", "COUNTER"},
	}
	for _, rm := range roommates {
		err := s.onboardRoommate(ctx, rm.name, rm.email, rm.chore)
		if err != nil && pgx.ErrNoRows.Error() != err.Error() {
			return err
		}
	}
	return nil
}

func (s *RotationService) onboardRoommate(ctx context.Context, name, email, chore string) error {
	_, err := s.dbQueries.UpsertRoommate(ctx, db.UpsertRoommateParams{
		Name:  name,
		Email: email,
		Chore: chore,
	})

	return err
}

// RotateChores rotates all roommates' chores to the next in the sequence
func (s *RotationService) RotateChores(ctx context.Context) ([]db.Roommate, error) {
	// Get all roommates
	roommates, err := s.dbQueries.GetRoommates(ctx)
	if err != nil {
		return nil, err
	}

	// Update each roommate's chore to the next one
	for _, roommate := range roommates {
		currentChore := ChoreType(roommate.Chore)
		nextChore := nextChore(currentChore)
		_, err := s.dbQueries.UpdateRoommateChore(ctx, db.UpdateRoommateChoreParams{
			ID:    roommate.ID,
			Chore: string(nextChore),
		})
		if err != nil {
			return nil, err
		}
		roommate.Chore = string(nextChore)
	}
	return roommates, nil
}

func (s *RotationService) GetRoommates(ctx context.Context) ([]db.Roommate, error) {
	// Fetch all roommates from the database
	roommates, err := s.dbQueries.GetRoommates(ctx)
	if err != nil {
		return nil, err
	}
	return roommates, nil
}

func (s *RotationService) CreateChoreDigest(rms []db.Roommate) *string {
	// Create a digest message for the chores
	digest := "Chore Rotation Digest:\n"
	for _, rm := range rms {
		digest += rm.Name + " is responsible for " + rm.Chore + "\n"
	}
	return &digest
}
