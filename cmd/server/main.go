package main

import (
	"log"
	"net/http"

	"github.com/thetangentline/rlimit/internal/server/handler"
	"github.com/thetangentline/rlimit/internal/server/service"
)

func main() {
	tokenService := service.NewTokenService()
	tokenHandler := handler.NewTokenHandler(tokenService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /token", tokenHandler.GetToken)
	log.Println("Standard server starting on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
