package store

import (
	"context"
	"database/sql"
	"errors"
)

type RoleStore struct {
	db *sql.DB
}
type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Level       int64  `json:"level"`
	Description string `json:"description"`
}
type RolesRepository interface {
	GetByName(context.Context, string) (*Role, error)
}

func (s *RoleStore) GetByName(ctx context.Context, slug string) (*Role, error) {
	query := `SELECT id, name, level, description FROM roles WHERE name = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	role := &Role{}
	err := s.db.QueryRowContext(ctx, query, slug).Scan(&role.ID, &role.Name, &role.Level, &role.Description)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return role, nil
}
