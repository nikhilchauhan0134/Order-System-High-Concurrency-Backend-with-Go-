package ratelimit

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// IPRateLimiter limits requests per IP
type IPRateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     int           // max requests
	interval time.Duration // per interval
}

// Visitor stores request count and last seen time
type Visitor struct {
	lastSeen time.Time
	count    int
}

// NewIPRateLimiter creates a new limiter
func NewIPRateLimiter(rate int, interval time.Duration) *IPRateLimiter {
	limiter := &IPRateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		interval: interval,
	}

	// cleanup stale visitors periodically
	go limiter.cleanup()
	return limiter
}

// Middleware wraps http.Handler
func (l *IPRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "invalid IP", http.StatusBadRequest)
			return
		}

		if !l.Allow(ip) {
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Allow checks if request is allowed
func (l *IPRateLimiter) Allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, exists := l.visitors[ip]
	if !exists || time.Since(v.lastSeen) > l.interval {
		l.visitors[ip] = &Visitor{
			lastSeen: time.Now(),
			count:    1,
		}
		return true
	}

	if v.count >= l.rate {
		return false
	}

	v.count++
	v.lastSeen = time.Now()
	return true
}

// cleanup removes stale visitors every minute
func (l *IPRateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > l.interval {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}
