package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestTokenStore_NilClient(t *testing.T) {
	ctx := context.Background()
	ts := &TokenStore{rdb: nil}

	err := ts.Blacklist(ctx, "test-jti", time.Minute)
	if err != nil {
		t.Fatalf("expected nil err when rdb is nil, got %v", err)
	}

	blacklisted, err := ts.IsBlacklisted(ctx, "test-jti")
	if err != nil || blacklisted {
		t.Fatalf("expected false, nil when rdb is nil, got %v, %v", blacklisted, err)
	}
}

func TestTokenStore_Integration(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("skipping integration test: Redis not reachable")
	}

	ts := &TokenStore{rdb: rdb}
	testJTI := "unit-test-jti-12345"

	// 1. Awalnya belum di-blacklist
	blacklisted, err := ts.IsBlacklisted(ctx, testJTI)
	if err != nil || blacklisted {
		t.Fatalf("expected false, got %v, %v", blacklisted, err)
	}

	// 2. Blacklist dengan TTL 300ms
	err = ts.Blacklist(ctx, testJTI, 300*time.Millisecond)
	if err != nil {
		t.Fatalf("failed to blacklist token: %v", err)
	}

	// 3. Sekarang harus blacklisted
	blacklisted, err = ts.IsBlacklisted(ctx, testJTI)
	if err != nil || !blacklisted {
		t.Fatalf("expected true, got %v, %v", blacklisted, err)
	}

	// 4. Setelah TTL lewat, harus otomatis un-blacklisted (expired)
	time.Sleep(350 * time.Millisecond)
	blacklisted, err = ts.IsBlacklisted(ctx, testJTI)
	if err != nil || blacklisted {
		t.Fatalf("expected false after TTL expiry, got %v, %v", blacklisted, err)
	}
}
