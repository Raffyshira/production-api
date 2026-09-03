package service

import (
	"context"

	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/raffyshira/project-rest-api/internal/store/cache"
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
	cache cache.Storage
}

func NewPostsService(store store.Storage, cache cache.Storage) PostService {
	return &postService{
		store: store,
		cache: cache,
	}
}

// Create implements [PostService].
func (p *postService) Create(ctx context.Context, post *store.Post) error {
	return p.store.Posts.Create(ctx, post)
}

// Delete implements [PostService].
func (p *postService) Delete(ctx context.Context, postID int64) error {
	if err := p.store.Posts.Delete(ctx, postID); err != nil {
		return err
	}

	if p.cache.Posts != nil {
		p.cache.Posts.Delete(ctx, postID)
		p.cache.Posts.DeleteComments(ctx, postID)
	}

	return nil
}

// GetByID implements [PostService].
func (p *postService) GetByID(ctx context.Context, postID int64) (*store.Post, error) {
	if p.cache.Posts != nil {
		cachedPost, err := p.cache.Posts.Get(ctx, postID)
		if err == nil && cachedPost != nil && cachedPost.ID != 0 {
			println("[CACHE HIT] Mengambil post dari Redis:", postID)
			return cachedPost, nil
		}
	}

	println("[CACHE MISS] Mengambil post dari PostgreSQL:", postID)

	post, err := p.store.Posts.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if p.cache.Posts != nil {
		_ = p.cache.Posts.Set(ctx, post)
	}

	return post, nil
}

// GetCommentsByID implements [PostService].
func (p *postService) GetCommentsByID(ctx context.Context, postID int64) ([]store.Comment, error) {
	if p.cache.Posts != nil {
		cachedComments, err := p.cache.Posts.GetComments(ctx, postID)
		if err == nil && cachedComments != nil {
			return cachedComments, nil
		}
	}

	comments, err := p.store.Comments.GetByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if p.cache.Posts != nil {
		_ = p.cache.Posts.SetComments(ctx, postID, comments)
	}

	return comments, nil
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
	if err := p.store.Posts.Update(ctx, post); err != nil {
		return err
	}

	if p.cache.Posts != nil {
		p.cache.Posts.Delete(ctx, post.ID)
	}

	return nil
}
