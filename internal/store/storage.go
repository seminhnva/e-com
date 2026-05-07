package store

import (
	"database/sql"
	"errors"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrEditConflict = errors.New("edit conflict")
	ErrConflict     = errors.New("resource already exists")
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
