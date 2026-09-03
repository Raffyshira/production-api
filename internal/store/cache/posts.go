package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/redis/go-redis/v9"
)

type PostStore struct {
	rdb *redis.Client
}

const PostExpTime = 10 * time.Minute

func (s *PostStore) Get(ctx context.Context, postID int64) (*store.Post, error) {
	if s.rdb == nil {
		return nil, nil
	}

	cacheKey := fmt.Sprintf("post-%d", postID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var post store.Post
	if data != "" {
		if err := json.Unmarshal([]byte(data), &post); err != nil {
			return nil, err
		}
	}

	return &post, nil
}

func (s *PostStore) Set(ctx context.Context, post *store.Post) error {
	if s.rdb == nil || post == nil {
		return nil
	}

	cacheKey := fmt.Sprintf("post-%d", post.ID)
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, data, PostExpTime).Err()
}

func (s *PostStore) Delete(ctx context.Context, postID int64) {
	if s.rdb == nil {
		return
	}

	cacheKey := fmt.Sprintf("post-%d", postID)
	s.rdb.Del(ctx, cacheKey)
}

func (s *PostStore) GetComments(ctx context.Context, postID int64) ([]store.Comment, error) {
	if s.rdb == nil {
		return nil, nil
	}

	cacheKey := fmt.Sprintf("post-%d-comments", postID)
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var comments []store.Comment
	if data != "" {
		if err := json.Unmarshal([]byte(data), &comments); err != nil {
			return nil, err
		}
	}

	return comments, nil
}

func (s *PostStore) SetComments(ctx context.Context, postID int64, comments []store.Comment) error {
	if s.rdb == nil {
		return nil
	}

	cacheKey := fmt.Sprintf("post-%d-comments", postID)
	data, err := json.Marshal(comments)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, data, PostExpTime).Err()
}

func (s *PostStore) DeleteComments(ctx context.Context, postID int64) {
	if s.rdb == nil {
		return
	}

	cacheKey := fmt.Sprintf("post-%d-comments", postID)
	s.rdb.Del(ctx, cacheKey)
}
