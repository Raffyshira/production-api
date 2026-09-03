package ratelimiter

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, ip string) (bool, time.Duration)
}

type Config struct {
	RequestsPerTimeFrame int
	TimeFrame            time.Duration
	Enabled              bool
}
