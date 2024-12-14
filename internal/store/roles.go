package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int64  `json:"level"`
	Description string `json:"description"`
}

type RolesStore struct {
	db *sql.DB
}

func (s *RolesStore) GetByName(ctx context.Context, roleName string) (*Role, error) {
	query := `SELECT id, name, level, description FROM roles WHERE name = $1`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var role Role
	err := s.db.QueryRowContext(ctx, query, roleName).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &role, nil
}
