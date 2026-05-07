package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type UsersRepository interface {
	GetById(context.Context, int64) (*User, error)
	Create(context.Context, *User) error
}

type User struct {
	ID        int64     `json:"id,omitempty"`
	Username  string    `json:"username"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type UsersStore struct {
	db *sql.DB
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
	INSERT INTO users(username, email, password)
	VALUES ($1, $2, $3)			
	RETURNING id, created_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) GetById(ctx context.Context, id int64) (*User, error) {
	query := `
	SELECT id, username, email, created_at
	FROM users
	WHERE id = $1
	`
	row := s.db.QueryRowContext(ctx, query, id)
	user := &User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}
