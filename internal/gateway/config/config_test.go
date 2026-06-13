package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()

	cfg := Load()

	if cfg.GatewayAddr != ":8081" {
		t.Fatalf("GatewayAddr = %q, want :8081", cfg.GatewayAddr)
	}
	if cfg.BackendURL != "http://localhost:8080" {
		t.Fatalf("BackendURL = %q, want http://localhost:8080", cfg.BackendURL)
	}
	if cfg.RedisAddr != "localhost:6379" {
		t.Fatalf("RedisAddr = %q, want localhost:6379", cfg.RedisAddr)
	}
	if cfg.RateLimit != 10 {
		t.Fatalf("RateLimit = %d, want 10", cfg.RateLimit)
	}
	if cfg.WindowSec != 60 {
		t.Fatalf("WindowSec = %d, want 60", cfg.WindowSec)
	}
	if cfg.KeyPrefix != "rlimit" {
		t.Fatalf("KeyPrefix = %q, want rlimit", cfg.KeyPrefix)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	t.Setenv("GATEWAY_ADDR", ":9090")
	t.Setenv("BACKEND_URL", "http://backend:8080")
	t.Setenv("REDIS_ADDR", "redis:6379")
	t.Setenv("RATE_LIMIT", "25")
	t.Setenv("WINDOW_SEC", "30")
	t.Setenv("KEY_PREFIX", "test")

	cfg := Load()

	if cfg.GatewayAddr != ":9090" {
		t.Fatalf("GatewayAddr = %q, want :9090", cfg.GatewayAddr)
	}
	if cfg.BackendURL != "http://backend:8080" {
		t.Fatalf("BackendURL = %q, want http://backend:8080", cfg.BackendURL)
	}
	if cfg.RedisAddr != "redis:6379" {
		t.Fatalf("RedisAddr = %q, want redis:6379", cfg.RedisAddr)
	}
	if cfg.RateLimit != 25 {
		t.Fatalf("RateLimit = %d, want 25", cfg.RateLimit)
	}
	if cfg.WindowSec != 30 {
		t.Fatalf("WindowSec = %d, want 30", cfg.WindowSec)
	}
	if cfg.KeyPrefix != "test" {
		t.Fatalf("KeyPrefix = %q, want test", cfg.KeyPrefix)
	}
}
