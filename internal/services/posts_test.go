package service

import (
	"context"
	"testing"

	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/raffyshira/project-rest-api/internal/store/cache"
	"github.com/stretchr/testify/mock"
)

func TestPostService_GetByID_CacheHit(t *testing.T) {
	ctx := context.Background()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockStore()

	cachedPost := &store.Post{
		ID:    10,
		Title: "Cached Post Title",
	}

	mockPostCache := mockCache.Posts.(*cache.MockPostStore)
	mockPostCache.On("Get", int64(10)).Return(cachedPost, nil).Once()

	svc := NewPostsService(mockStore, mockCache)
	post, err := svc.GetByID(ctx, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.Title != cachedPost.Title {
		t.Fatalf("expected title %s, got %s", cachedPost.Title, post.Title)
	}

	mockPostCache.AssertExpectations(t)
}

func TestPostService_GetByID_CacheMiss_ThenSet(t *testing.T) {
	ctx := context.Background()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockStore()

	mockPostCache := mockCache.Posts.(*cache.MockPostStore)
	// Cache miss: return nil, nil
	mockPostCache.On("Get", int64(20)).Return(nil, nil).Once()
	// Should then set the fetched post into cache
	mockPostCache.On("Set", mock.AnythingOfType("*store.Post")).Return(nil).Once()

	svc := NewPostsService(mockStore, mockCache)
	post, err := svc.GetByID(ctx, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if post.ID != 20 {
		t.Fatalf("expected post ID 20, got %d", post.ID)
	}

	mockPostCache.AssertExpectations(t)
}

func TestPostService_Update_InvalidatesCache(t *testing.T) {
	ctx := context.Background()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockStore()

	mockPostCache := mockCache.Posts.(*cache.MockPostStore)
	mockPostCache.On("Delete", int64(30)).Return().Once()

	svc := NewPostsService(mockStore, mockCache)
	err := svc.Update(ctx, &store.Post{ID: 30, Title: "Updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mockPostCache.AssertExpectations(t)
}

func TestPostService_Delete_InvalidatesCacheAndComments(t *testing.T) {
	ctx := context.Background()
	mockStore := store.NewMockStore()
	mockCache := cache.NewMockStore()

	mockPostCache := mockCache.Posts.(*cache.MockPostStore)
	mockPostCache.On("Delete", int64(40)).Return().Once()
	mockPostCache.On("DeleteComments", int64(40)).Return().Once()

	svc := NewPostsService(mockStore, mockCache)
	err := svc.Delete(ctx, 40)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	mockPostCache.AssertExpectations(t)
}
