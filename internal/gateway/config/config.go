package config

import (
	"os"
	"strconv"
)

type Config struct {
	GatewayAddr string
	BackendURL  string
	RedisAddr   string
	RateLimit   int
	WindowSec   int
	KeyPrefix   string
}

func Load() Config {
	return Config{
		GatewayAddr: envOrDefault("GATEWAY_ADDR", ":8081"),
		BackendURL:  envOrDefault("BACKEND_URL", "http://localhost:8080"),
		RedisAddr:   envOrDefault("REDIS_ADDR", "localhost:6379"),
		RateLimit:   envIntOrDefault("RATE_LIMIT", 10),
		WindowSec:   envIntOrDefault("WINDOW_SEC", 60),
		KeyPrefix:   envOrDefault("KEY_PREFIX", "rlimit"),
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envIntOrDefault(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}
