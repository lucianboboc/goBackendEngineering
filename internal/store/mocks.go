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
}

func (s *MockUserStore) Create(context.Context, *sql.Tx, *User) error {
	return nil
}
func (s *MockUserStore) GetUserByID(context.Context, int64) (*User, error) {
	return nil, nil
}
func (s *MockUserStore) UpdateUser(ctx context.Context, user *User) error {
	return nil
}
func (s *MockUserStore) DeleteUser(ctx context.Context, userID int64) error {
	return nil
}
func (s *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return nil
}
func (s *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}
func (s *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}
