package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockStore() Storage {
	return Storage{
		Users:    &MockUserStore{},
		Posts:    &MockPostStore{},
		Comments: &MockCommentsStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	return nil
}

func (m *MockUserStore) GetByID(ctx context.Context, userID int64) (*User, error) {
	return &User{ID: userID}, nil
}

func (m *MockUserStore) GetByEmail(context.Context, string) (*User, error) {
	return &User{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, t string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id int64) error {
	return nil
}

type MockPostStore struct{}

func (m *MockPostStore) GetByID(ctx context.Context, id int64) (*Post, error) {
	return &Post{ID: id, Title: "Mock Post"}, nil
}

func (m *MockPostStore) Create(ctx context.Context, post *Post) error {
	return nil
}

func (m *MockPostStore) Delete(ctx context.Context, id int64) error {
	return nil
}

func (m *MockPostStore) Update(ctx context.Context, post *Post) error {
	return nil
}

func (m *MockPostStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	return nil, nil
}

func (m *MockPostStore) GetExploreFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	return nil, nil
}

type MockCommentsStore struct{}

func (m *MockCommentsStore) Create(ctx context.Context, comment *Comment) error {
	return nil
}

func (m *MockCommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {
	return []Comment{{ID: 1, PostID: postID, Content: "Mock Comment"}}, nil
}
