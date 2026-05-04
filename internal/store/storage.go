package store

import "database/sql"

type Storage struct {
	Posts PostsRepository
	Users UsersRepository
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts: &PostsStore{db},
		Users: &UsersStore{db},
	}
}
