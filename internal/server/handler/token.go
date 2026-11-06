package handler

import (
	"fmt"
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
	token, err := th.svc.GenToken()
	if err != nil {
		log.Printf("failed to generate token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, token)
}
