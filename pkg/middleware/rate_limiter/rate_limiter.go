package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex

	cleanupVisitor = func(ip string) {
		time.AfterFunc(10*time.Minute, func() {
			mu.Lock()
			delete(visitors, ip)
			mu.Unlock()
		})
	}
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, 3) // 1 req/segundo, burst 3
		visitors[ip] = limiter
		cleanupVisitor(ip)
	}
	return limiter
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := getVisitor(r.RemoteAddr)
		if !limiter.Allow() {
			http.Error(w, "Muitas requisições. Tente novamente mais tarde.", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
