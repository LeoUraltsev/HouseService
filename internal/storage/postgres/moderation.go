package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (s *Storage) SelectModeratorByFlatID(ctx context.Context, flatID int) (uuid.UUID, error) {
	var moderationID uuid.UUID
	q := `select moderator_id from moderation where flat_id = $1`

	err := s.Pool.QueryRow(ctx, q, flatID).Scan(&moderationID)
	if err != nil {
		return uuid.Nil, err
	}
	return moderationID, nil
}

func (s *Storage) InsertModeratorForFlat(ctx context.Context, moderatorID uuid.UUID, flatID int) error {

	q := `insert into moderation (moderator_id, flat_id) values ($1, $2)`

	_, err := s.Pool.Exec(ctx, q, moderatorID, flatID)
	if err != nil {
		return err
	}

	return nil
}
