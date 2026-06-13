package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Token struct {
	token string
	exp   time.Time
}

type TokenResponse struct {
	Token string    `json:"token"`
	Exp   time.Time `json:"exp"`
}

type TokenService struct {
}

func NewToken(token string, exp time.Time) *Token {
	res := Token{token: token, exp: exp}
	return &res
}

func NewTokenService() *TokenService {
	return &TokenService{}
}

func (t *Token) Response() TokenResponse {
	return TokenResponse{
		Token: t.token,
		Exp:   t.exp,
	}
}

func (tsv *TokenService) GenToken(ctx context.Context) (*Token, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	token, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}
	exp := time.Now().Add(time.Hour)
	res := NewToken(token.String(), exp)
	return res, nil
}
