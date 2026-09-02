package service

import (
	"context"

	"github.com/raffyshira/project-rest-api/internal/store"
)

type PostService interface {
	GetByID(context.Context, int64) (*store.Post, error)
	Create(context.Context, *store.Post) error
	Delete(context.Context, int64) error
	Update(context.Context, *store.Post) error
	GetCommentsByID(context.Context, int64) ([]store.Comment, error)
	GetUserFeed(context.Context, int64, store.PaginatedFeedQuery) ([]store.PostWithMetadata, error)
	GetExploreFeed(context.Context, int64, store.PaginatedFeedQuery) ([]store.PostWithMetadata, error)
}

type postService struct {
	store store.Storage
}

func NewPostsService(store store.Storage) PostService {
	return &postService{
		store: store,
	}
}

// Create implements [PostService].
func (p *postService) Create(ctx context.Context, post *store.Post) error {
	return p.store.Posts.Create(ctx, post)
}

// Delete implements [PostService].
func (p *postService) Delete(ctx context.Context, postID int64) error {
	return p.store.Posts.Delete(ctx, postID)
}

// GetByID implements [PostService].
func (p *postService) GetByID(ctx context.Context, postID int64) (*store.Post, error) {
	return p.store.Posts.GetByID(ctx, postID)
}

// GetCommentsByID implements [PostService].
func (p *postService) GetCommentsByID(ctx context.Context, postID int64) ([]store.Comment, error) {
	return p.store.Comments.GetByPostID(ctx, postID)
}

// GetUserFeed implements [PostService].
func (p *postService) GetUserFeed(ctx context.Context, userID int64, query store.PaginatedFeedQuery) ([]store.PostWithMetadata, error) {
	return p.store.Posts.GetUserFeed(ctx, userID, query)
}

// GetExploreFeed implements [PostService].
func (p *postService) GetExploreFeed(ctx context.Context, userID int64, query store.PaginatedFeedQuery) ([]store.PostWithMetadata, error) {
	return p.store.Posts.GetExploreFeed(ctx, userID, query)
}

// Update implements [PostService].
func (p *postService) Update(ctx context.Context, post *store.Post) error {
	return p.store.Posts.Update(ctx, post)
}
