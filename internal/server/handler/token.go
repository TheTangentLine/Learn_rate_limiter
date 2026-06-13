package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/thetangentline/rlimit/internal/server/service"
)

type TokenHandler struct {
	svc *service.TokenService
}

func NewTokenHandler(svc *service.TokenService) *TokenHandler {
	return &TokenHandler{svc: svc}
}

func (th *TokenHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	token, err := th.svc.GenToken(r.Context())
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(token.Response()); err != nil {
		log.Printf("failed to encode token response: %v", err)
	}
}
