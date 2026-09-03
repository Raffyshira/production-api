package cache

import (
	"context"
	"testing"

	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/redis/go-redis/v9"
)

func TestPostStore_NilClient(t *testing.T) {
	ctx := context.Background()
	ps := &PostStore{rdb: nil}

	post, err := ps.Get(ctx, 1)
	if err != nil || post != nil {
		t.Fatalf("expected nil, nil with nil rdb, got %v, %v", post, err)
	}

	if err := ps.Set(ctx, &store.Post{ID: 1}); err != nil {
		t.Fatalf("expected nil err with nil rdb, got %v", err)
	}

	ps.Delete(ctx, 1)

	comments, err := ps.GetComments(ctx, 1)
	if err != nil || comments != nil {
		t.Fatalf("expected nil, nil comments with nil rdb, got %v, %v", comments, err)
	}

	if err := ps.SetComments(ctx, 1, []store.Comment{{ID: 1}}); err != nil {
		t.Fatalf("expected nil err with nil rdb, got %v", err)
	}

	ps.DeleteComments(ctx, 1)
}

func TestPostStore_Integration(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("skipping integration test: Redis not reachable")
	}

	ps := &PostStore{rdb: rdb}
	postID := int64(99999)

	ps.Delete(ctx, postID)
	ps.DeleteComments(ctx, postID)

	// 1. Get saat kosong -> nil, nil
	p, err := ps.Get(ctx, postID)
	if err != nil || p != nil {
		t.Fatalf("expected nil post on cache miss, got %v, %v", p, err)
	}

	// 2. Set post
	samplePost := &store.Post{
		ID:      postID,
		Title:   "Cache Test Title",
		Content: "Cache Test Content",
		Tags:    []string{"go", "redis"},
	}
	if err := ps.Set(ctx, samplePost); err != nil {
		t.Fatalf("failed to set post in cache: %v", err)
	}

	// 3. Get post -> harus hit
	cached, err := ps.Get(ctx, postID)
	if err != nil || cached == nil {
		t.Fatalf("expected cached post, got %v, %v", cached, err)
	}
	if cached.Title != samplePost.Title {
		t.Fatalf("expected title %s, got %s", samplePost.Title, cached.Title)
	}

	// 4. Comments Set & Get
	sampleComments := []store.Comment{
		{ID: 1, PostID: postID, Content: "Nice post!"},
	}
	if err := ps.SetComments(ctx, postID, sampleComments); err != nil {
		t.Fatalf("failed to set comments in cache: %v", err)
	}

	cachedComments, err := ps.GetComments(ctx, postID)
	if err != nil || len(cachedComments) != 1 {
		t.Fatalf("expected 1 cached comment, got %v, %v", cachedComments, err)
	}

	// 5. Delete post and comments
	ps.Delete(ctx, postID)
	ps.DeleteComments(ctx, postID)

	deletedPost, _ := ps.Get(ctx, postID)
	if deletedPost != nil {
		t.Fatal("expected post to be deleted from cache")
	}

	deletedComments, _ := ps.GetComments(ctx, postID)
	if deletedComments != nil {
		t.Fatal("expected comments to be deleted from cache")
	}
}
