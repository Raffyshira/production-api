package ratelimiter

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
)

var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local member = ARGV[4]

local clearBefore = now - window
redis.call('ZREMRANGEBYSCORE', key, '-inf', clearBefore)

local currentRequests = redis.call('ZCARD', key)
if currentRequests < limit then
    redis.call('ZADD', key, now, member)
    redis.call('PEXPIRE', key, window)
    return {1, 0}
else
    local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
    local retryAfterMs = window
    if #oldest >= 2 then
        local oldestTime = tonumber(oldest[2])
        local remaining = (oldestTime + window) - now
        if remaining > 0 then
            retryAfterMs = remaining
        else
            retryAfterMs = 0
        end
    end
    return {0, retryAfterMs}
end
`)

type RedisSlidingWindowLimiter struct {
	rdb     *redis.Client
	limit   int
	window  time.Duration
	counter uint64
}

func NewRedisSlidingWindowLimiter(rdb *redis.Client, limit int, window time.Duration) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{
		rdb:    rdb,
		limit:  limit,
		window: window,
	}
}

func (rl *RedisSlidingWindowLimiter) Allow(ctx context.Context, ip string) (bool, time.Duration) {
	if rl.rdb == nil {
		return true, 0
	}

	callCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	now := time.Now()
	nowMs := now.UnixMilli()
	windowMs := rl.window.Milliseconds()
	count := atomic.AddUint64(&rl.counter, 1)
	member := fmt.Sprintf("%d-%d", now.UnixNano(), count)
	key := fmt.Sprintf("ratelimit:%s", ip)

	res, err := slidingWindowScript.Run(callCtx, rl.rdb, []string{key}, nowMs, windowMs, rl.limit, member).Slice()
	if err != nil {
		// Fail-open: jika Redis down/timeout, izinkan request lewat agar API tidak lumpuh
		return true, 0
	}

	if len(res) < 2 {
		return true, 0
	}

	allowed, _ := res[0].(int64)
	retryAfterMs, _ := res[1].(int64)

	if allowed == 1 {
		return true, 0
	}

	return false, time.Duration(retryAfterMs) * time.Millisecond
}
