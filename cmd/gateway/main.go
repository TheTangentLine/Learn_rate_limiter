package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thetangentline/rlimit/internal/gateway/config"
	gatewayhandler "github.com/thetangentline/rlimit/internal/gateway/handler"
	"github.com/thetangentline/rlimit/internal/gateway/ratelimit"
)

func main() {
	cfg := config.Load()

	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})
	defer redisClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("warning: redis ping failed: %v", err)
	}
	cancel()

	limiter := ratelimit.NewFixedWindowLimiter(redisClient, cfg.RateLimit, cfg.WindowSec)

	proxyHandler, err := gatewayhandler.NewProxyHandler(cfg, limiter)
	if err != nil {
		log.Fatalf("failed to create proxy handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /token", proxyHandler)

	server := &http.Server{
		Addr:    cfg.GatewayAddr,
		Handler: mux,
	}

	go func() {
		log.Printf("Gateway starting on %s (backend: %s, redis: %s)", cfg.GatewayAddr, cfg.BackendURL, cfg.RedisAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gateway server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("gateway shutdown error: %v", err)
	}
}
