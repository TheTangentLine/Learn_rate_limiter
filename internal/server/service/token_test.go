package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenToken_ReturnsValidToken(t *testing.T) {
	svc := NewTokenService()
	token, err := svc.GenToken(context.Background())
	if err != nil {
		t.Fatalf("GenToken() error = %v", err)
	}
	if token == nil {
		t.Fatal("GenToken() returned nil token")
	}

	resp := token.Response()
	if _, err := uuid.Parse(resp.Token); err != nil {
		t.Fatalf("token is not a valid UUID: %v", err)
	}
	if resp.Exp.Before(time.Now()) {
		t.Fatal("token expiry is in the past")
	}
	if resp.Exp.After(time.Now().Add(time.Hour + time.Minute)) {
		t.Fatal("token expiry is more than ~1 hour from now")
	}
}

func TestGenToken_RespectsCancelledContext(t *testing.T) {
	svc := NewTokenService()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := svc.GenToken(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestResponse_MapsFields(t *testing.T) {
	exp := time.Now().Add(time.Hour).UTC().Truncate(time.Second)
	token := NewToken("test-token-value", exp)

	resp := token.Response()
	if resp.Token != "test-token-value" {
		t.Fatalf("Token = %q, want %q", resp.Token, "test-token-value")
	}
	if !resp.Exp.Equal(exp) {
		t.Fatalf("Exp = %v, want %v", resp.Exp, exp)
	}
}
