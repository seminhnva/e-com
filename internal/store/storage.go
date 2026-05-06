package store

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound    = errors.New("resource not found")
	ErrEditConflict = errors.New("edit conflict")
)

type Storage struct {
	Posts    PostsRepository
	Users    UsersRepository
	Comments CommentsRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
