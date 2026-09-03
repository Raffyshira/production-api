package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestFixedWindowRateLimiter(t *testing.T) {
	ctx := context.Background()
	limiter := NewFixedWindowRateLimiter(3, 100*time.Millisecond)

	ip := "192.168.1.1"

	// 3 request pertama harus diizinkan
	for i := 0; i < 3; i++ {
		allow, _ := limiter.Allow(ctx, ip)
		if !allow {
			t.Fatalf("request %d expected to be allowed", i+1)
		}
	}

	// Request ke-4 harus ditolak
	allow, retryAfter := limiter.Allow(ctx, ip)
	if allow {
		t.Fatal("request 4 expected to be blocked")
	}
	if retryAfter <= 0 {
		t.Fatalf("expected positive retryAfter, got %v", retryAfter)
	}

	// Setelah window lewat, harus diizinkan lagi
	time.Sleep(110 * time.Millisecond)
	allow, _ = limiter.Allow(ctx, ip)
	if !allow {
		t.Fatal("request expected to be allowed after window reset")
	}
}

func TestRedisSlidingWindowLimiter_FailOpenOnNil(t *testing.T) {
	ctx := context.Background()
	limiter := NewRedisSlidingWindowLimiter(nil, 3, time.Second)

	// Jika rdb nil, harus fail-open (allow = true)
	allow, _ := limiter.Allow(ctx, "127.0.0.1")
	if !allow {
		t.Fatal("expected allow=true when redis client is nil (fail-open)")
	}
}

func TestRedisSlidingWindowLimiter_Integration(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	defer rdb.Close()

	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("skipping Redis integration test: Redis not reachable on 127.0.0.1:6379")
	}

	testIP := "10.0.0.99"
	rdb.Del(ctx, "ratelimit:"+testIP)

	limiter := NewRedisSlidingWindowLimiter(rdb, 3, 300*time.Millisecond)

	for i := 0; i < 3; i++ {
		allow, _ := limiter.Allow(ctx, testIP)
		if !allow {
			t.Fatalf("request %d expected to be allowed", i+1)
		}
	}

	// Request ke-4 harus ditolak
	allow, retryAfter := limiter.Allow(ctx, testIP)
	if allow {
		t.Fatal("request 4 expected to be blocked by sliding window")
	}
	if retryAfter <= 0 {
		t.Fatalf("expected retryAfter > 0, got %v", retryAfter)
	}

	// Tunggu sliding window lewat
	time.Sleep(350 * time.Millisecond)
	allow, _ = limiter.Allow(ctx, testIP)
	if !allow {
		t.Fatal("request expected to be allowed after sliding window expired")
	}

	rdb.Del(ctx, "ratelimit:"+testIP)
}
