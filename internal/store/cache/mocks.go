package cache

import (
	"context"

	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/stretchr/testify/mock"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
		Posts: &MockPostStore{},
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

func (m *MockUserStore) Delete(ctx context.Context, userID int64) {
	m.Called(userID)
}

type MockPostStore struct {
	mock.Mock
}

func (s *MockPostStore) Get(ctx context.Context, postID int64) (*store.Post, error) {
	args := s.Called(postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*store.Post), args.Error(1)
}

func (s *MockPostStore) Set(ctx context.Context, post *store.Post) error {
	args := s.Called(post)
	return args.Error(0)
}

func (s *MockPostStore) Delete(ctx context.Context, postID int64) {
	s.Called(postID)
}

func (s *MockPostStore) GetComments(ctx context.Context, postID int64) ([]store.Comment, error) {
	args := s.Called(postID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]store.Comment), args.Error(1)
}

func (s *MockPostStore) SetComments(ctx context.Context, postID int64, comments []store.Comment) error {
	args := s.Called(postID, comments)
	return args.Error(0)
}

func (s *MockPostStore) DeleteComments(ctx context.Context, postID int64) {
	s.Called(postID)
}
