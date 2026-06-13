package ratelimit

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestFixedWindowLimiter_AllowsUpToLimit(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	limiter := NewFixedWindowLimiter(client, 3, 60)

	ctx := context.Background()
	key := "test:client:1"

	for i := 0; i < 3; i++ {
		result, err := limiter.Allow(ctx, key)
		if err != nil {
			t.Fatalf("request %d: Allow() error = %v", i+1, err)
		}
		if !result.Allowed {
			t.Fatalf("request %d: expected allowed", i+1)
		}
	}

	result, err := limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("request 4: Allow() error = %v", err)
	}
	if result.Allowed {
		t.Fatal("request 4: expected denied")
	}
	if result.Remaining != 0 {
		t.Fatalf("Remaining = %d, want 0", result.Remaining)
	}
}

func TestFixedWindowLimiter_KeyExpiresAfterWindow(t *testing.T) {
	mr := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	limiter := NewFixedWindowLimiter(client, 1, 1)

	ctx := context.Background()
	key := "test:client:expire"

	result, err := limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("first Allow() error = %v", err)
	}
	if !result.Allowed {
		t.Fatal("first request should be allowed")
	}

	result, err = limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("second Allow() error = %v", err)
	}
	if result.Allowed {
		t.Fatal("second request should be denied")
	}

	mr.FastForward(2 * 1e9)

	result, err = limiter.Allow(ctx, key)
	if err != nil {
		t.Fatalf("third Allow() after expiry error = %v", err)
	}
	if !result.Allowed {
		t.Fatal("request after window expiry should be allowed")
	}
}
