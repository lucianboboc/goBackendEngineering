package store

import (
	"context"
	"database/sql"
	"time"
)

type Follower struct {
	UserID     int64      `json:"user_id"`
	FollowerID int64      `json:"follower_id"`
	CreatedAt  *time.Time `json:"created_at"`
}

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) Follow(ctx context.Context, followerID, userID int64) error {
	query := `INSERT INTO followers(user_id, follower_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *FollowersStore) Unfollow(ctx context.Context, followerID, userID int64) error {
	query := `DELETE FROM followers WHERE user_id = $1 AND follower_id = $2`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
