package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	GetByIdFn   func(ctx context.Context, id int64) (*User, error)
	GetByEmailFn func(ctx context.Context, email string) (*User, error)
}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	return nil
}
func (m *MockUserStore) GetById(ctx context.Context, id int64) (*User, error) {
	if m.GetByIdFn != nil {
		return m.GetByIdFn(ctx, id)
	}
	return &User{}, nil
}
func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return nil
}
func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}
func (m *MockUserStore) Delete(ctx context.Context, userID int64) error {
	return nil
}
func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return &User{}, nil
}
