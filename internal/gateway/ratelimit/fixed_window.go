package ratelimit

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const fixedWindowScript = `
local count = redis.call('INCR', KEYS[1])
if count == 1 then
    redis.call('EXPIRE', KEYS[1], ARGV[2])
end
return count
`

type FixedWindowLimiter struct {
	client    *redis.Client
	limit     int
	windowSec int
}

func NewFixedWindowLimiter(client *redis.Client, limit, windowSec int) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		client:    client,
		limit:     limit,
		windowSec: windowSec,
	}
}

func (l *FixedWindowLimiter) Allow(ctx context.Context, key string) (Result, error) {
	count, err := l.client.Eval(ctx, fixedWindowScript, []string{key}, l.limit, l.windowSec).Int()
	if err != nil {
		return Result{}, fmt.Errorf("redis eval: %w", err)
	}

	allowed := count <= l.limit
	remaining := l.limit - count
	if remaining < 0 {
		remaining = 0
	}

	return Result{
		Allowed:   allowed,
		Remaining: remaining,
	}, nil
}
