package service

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	token string
	exp   time.Time
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

func (tsv *TokenService) GenToken() (*Token, error) {
	token, err := uuid.NewV6()
	if err != nil {
		return nil, err
	}
	exp := time.Now().Add(time.Hour)
	res := NewToken(token.String(), exp)
	return res, nil
}
