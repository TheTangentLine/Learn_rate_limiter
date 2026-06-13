package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thetangentline/rlimit/internal/gateway/config"
	"github.com/thetangentline/rlimit/internal/gateway/ratelimit"
)

type mockLimiter struct {
	result ratelimit.Result
	err    error
}

func (m *mockLimiter) Allow(_ context.Context, _ string) (ratelimit.Result, error) {
	return m.result, m.err
}

func TestProxyHandler_Allowed_ProxiesToBackend(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token":"abc"}`))
	}))
	defer backend.Close()

	cfg := config.Config{
		BackendURL: backend.URL,
		WindowSec:  60,
		KeyPrefix:  "test",
	}
	h, err := NewProxyHandler(cfg, &mockLimiter{result: ratelimit.Result{Allowed: true}})
	if err != nil {
		t.Fatalf("NewProxyHandler() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if rec.Body.String() != `{"token":"abc"}` {
		t.Fatalf("body = %q, want %q", rec.Body.String(), `{"token":"abc"}`)
	}
}

func TestProxyHandler_Denied_Returns429(t *testing.T) {
	cfg := config.Config{
		BackendURL: "http://localhost:1",
		WindowSec:  60,
		KeyPrefix:  "test",
	}
	h, err := NewProxyHandler(cfg, &mockLimiter{result: ratelimit.Result{Allowed: false}})
	if err != nil {
		t.Fatalf("NewProxyHandler() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
	if got := rec.Header().Get("Retry-After"); got != "60" {
		t.Fatalf("Retry-After = %q, want %q", got, "60")
	}
}

func TestProxyHandler_RedisError_Returns503(t *testing.T) {
	cfg := config.Config{
		BackendURL: "http://localhost:1",
		WindowSec:  60,
		KeyPrefix:  "test",
	}
	h, err := NewProxyHandler(cfg, &mockLimiter{err: errors.New("redis unavailable")})
	if err != nil {
		t.Fatalf("NewProxyHandler() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
	}
}
