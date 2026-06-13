package identity

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientID_FromXForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.1, 198.51.100.2")

	got := ClientID(req)
	if got != "203.0.113.1" {
		t.Fatalf("ClientID() = %q, want %q", got, "203.0.113.1")
	}
}

func TestClientID_FromRemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	req.RemoteAddr = "192.168.1.10:54321"

	got := ClientID(req)
	if got != "192.168.1.10" {
		t.Fatalf("ClientID() = %q, want %q", got, "192.168.1.10")
	}
}

func TestRateLimitKey_WindowBucketing(t *testing.T) {
	windowSec := 60
	now := time.Unix(120, 0)

	key1 := RateLimitKey("rlimit", "1.2.3.4", windowSec, now)
	key2 := RateLimitKey("rlimit", "1.2.3.4", windowSec, now.Add(30*time.Second))
	key3 := RateLimitKey("rlimit", "1.2.3.4", windowSec, now.Add(60*time.Second))

	if key1 != key2 {
		t.Fatalf("keys in same window should match: %q vs %q", key1, key2)
	}
	if key1 == key3 {
		t.Fatalf("keys in different windows should differ: %q vs %q", key1, key3)
	}
	if key1 != "rlimit:1.2.3.4:2" {
		t.Fatalf("RateLimitKey() = %q, want %q", key1, "rlimit:1.2.3.4:2")
	}
}
