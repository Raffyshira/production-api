package cache

import (
	"context"

	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/stretchr/testify/mock"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	mock.Mock
}

func (s *MockUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	args := s.Called(userID)
	return nil, args.Error(1)
}

func (s *MockUserStore) Set(ctx context.Context, user *store.User) error {
	args := s.Called(user)
	return args.Error(0)
}
