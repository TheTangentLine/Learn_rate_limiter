package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"time"

	"github.com/thetangentline/rlimit/internal/gateway/config"
	"github.com/thetangentline/rlimit/internal/gateway/identity"
	"github.com/thetangentline/rlimit/internal/gateway/ratelimit"
)

type ProxyHandler struct {
	limiter ratelimit.Limiter
	proxy   *httputil.ReverseProxy
	cfg     config.Config
}

func NewProxyHandler(cfg config.Config, limiter ratelimit.Limiter) (*ProxyHandler, error) {
	backendURL, err := url.Parse(cfg.BackendURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}

	return &ProxyHandler{
		limiter: limiter,
		proxy:   proxy,
		cfg:     cfg,
	}, nil
}

func (h *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientID := identity.ClientID(r)
	key := identity.RateLimitKey(h.cfg.KeyPrefix, clientID, h.cfg.WindowSec, time.Now())

	result, err := h.limiter.Allow(r.Context(), key)
	if err != nil {
		log.Printf("rate limit check failed: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if !result.Allowed {
		w.Header().Set("Retry-After", strconv.Itoa(h.cfg.WindowSec))
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}

	h.proxy.ServeHTTP(w, r)
}
