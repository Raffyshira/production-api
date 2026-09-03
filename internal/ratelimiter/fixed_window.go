package ratelimiter

import (
	"context"
	"sync"
	"time"
)

type clientWindow struct {
	count     int
	lastReset time.Time
}

type FixedWindowRateLimiter struct {
	sync.Mutex
	clients map[string]*clientWindow
	limit   int
	window  time.Duration
}

func NewFixedWindowRateLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]*clientWindow),
		limit:   limit,
		window:  window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ctx context.Context, ip string) (bool, time.Duration) {
	rl.Lock()
	defer rl.Unlock()

	now := time.Now()
	client, exist := rl.clients[ip]
	if !exist || now.Sub(client.lastReset) >= rl.window {
		rl.clients[ip] = &clientWindow{
			count:     1,
			lastReset: now,
		}
		return true, 0
	}

	if client.count < rl.limit {
		client.count++
		return true, 0
	}

	retryAfter := rl.window - now.Sub(client.lastReset)
	if retryAfter < 0 {
		retryAfter = 0
	}

	return false, retryAfter
}
