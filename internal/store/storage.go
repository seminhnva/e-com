package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const QueryTimeoutDuration = 5 * time.Second

var (
	ErrNotFound          = errors.New("resource not found")
	ErrEditConflict      = errors.New("edit conflict")
	ErrConflict          = errors.New("resource already exists")
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type Storage struct {
	Posts    PostsRepository
	Users    UsersRepository
	Follower FollowerRepository
	Comments CommentsRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Follower: &FollowerStore{db},
		Comments: &CommentsStore{db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()

		return err
	}
	return tx.Commit()
}
