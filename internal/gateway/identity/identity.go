package identity

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ClientID(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if ip := strings.TrimSpace(parts[0]); ip != "" {
			return ip
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func RateLimitKey(prefix, clientID string, windowSec int, now time.Time) string {
	windowStart := now.Unix() / int64(windowSec)
	return prefix + ":" + clientID + ":" + strconv.FormatInt(windowStart, 10)
}
