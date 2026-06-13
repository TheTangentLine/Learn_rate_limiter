package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thetangentline/rlimit/internal/server/service"
)

func TestGetToken_ReturnsJSON(t *testing.T) {
	svc := service.NewTokenService()
	h := NewTokenHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/token", nil)
	rec := httptest.NewRecorder()

	h.GetToken(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if _, ok := body["token"]; !ok {
		t.Fatal("response missing token field")
	}
	if _, ok := body["exp"]; !ok {
		t.Fatal("response missing exp field")
	}
}
