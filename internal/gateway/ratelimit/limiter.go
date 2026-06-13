package ratelimit

import "context"

type Result struct {
	Allowed   bool
	Remaining int
}

type Limiter interface {
	Allow(ctx context.Context, key string) (Result, error)
}
