package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
	once    sync.Once
)

// getClient retrieves or creates a rate limiter for a given IP
func getClient(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	c, exists := clients[ip]
	if !exists {
		limiter := rate.NewLimiter(rate.Every(time.Second), 2) // 1 request per second, burst of 3
		clients[ip] = &client{limiter, time.Now()}
		return limiter
	}

	c.lastSeen = time.Now()
	return c.limiter
}

// cleanupClients removes old clients from the map
func cleanupClients() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, c := range clients {
			if time.Since(c.lastSeen) > 3*time.Minute {
				delete(clients, ip)
			}
		}
		mu.Unlock()
	}
}

// RateLimitMiddleware limits the number of requests per client IP
func RateLimitMiddleware(next http.Handler) http.Handler {
	// start cleanup goroutine only once
	once.Do(func() {
		go cleanupClients()
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := getClient(ip)

		if !limiter.Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
