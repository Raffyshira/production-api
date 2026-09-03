package cache

import (
	"context"

	"github.com/raffyshira/project-rest-api/internal/store"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
		Delete(context.Context, int64)
	}
	Posts interface {
		Get(context.Context, int64) (*store.Post, error)
		Set(context.Context, *store.Post) error
		Delete(context.Context, int64)
		GetComments(context.Context, int64) ([]store.Comment, error)
		SetComments(context.Context, int64, []store.Comment) error
		DeleteComments(context.Context, int64)
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rbd},
		Posts: &PostStore{rdb: rbd},
	}
}
