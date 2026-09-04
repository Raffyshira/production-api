package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStore struct {
	rdb *redis.Client
}

func (s *TokenStore) Blacklist(ctx context.Context, jti string, ttl time.Duration) error {
	if s.rdb == nil || jti == "" || ttl <= 0 {
		return nil
	}

	key := fmt.Sprintf("blacklist:jti:%s", jti)
	return s.rdb.Set(ctx, key, "1", ttl).Err()
}

func (s *TokenStore) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	if s.rdb == nil || jti == "" {
		return false, nil
	}

	key := fmt.Sprintf("blacklist:jti:%s", jti)
	exists, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
