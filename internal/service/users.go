package service

import (
	"context"

	"github.com/raffyshira/project-rest-api/internal/store"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID int64) (*store.User, error)
	FollowUser(ctx context.Context, userID, followerID int64) error
	UnfollowUser(ctx context.Context, followerID, userID int64) error
	ActivateUser(ctx context.Context, token string) error
}

type userService struct {
	store store.Storage
}

func NewUserService(store store.Storage) UserService {
	return &userService{
		store: store,
	}
}

func (s *userService) GetUserByID(ctx context.Context, userID int64) (*store.User, error) {
	return s.store.Users.GetByID(ctx, userID)
}

func (s *userService) FollowUser(ctx context.Context, userID, followerID int64) error {
	return s.store.Followers.Follow(ctx, userID, followerID)
}

func (s *userService) UnfollowUser(ctx context.Context, followerID, userID int64) error {
	return s.store.Followers.Unfollow(ctx, followerID, userID)
}

func (s *userService) ActivateUser(ctx context.Context, token string) error {
	return s.store.Users.Activate(ctx, token)
}
