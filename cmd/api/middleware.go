package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

type Client struct {
	mu       sync.Mutex
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	Clients map[string]*Client
}

func NewRateLimiter() (*RateLimiter, func(*RateLimiter)) {
	rl := &RateLimiter{
		Clients: make(map[string]*Client),
	}
	cleanUp := func(rl *RateLimiter) {
		for {
			time.Sleep(time.Minute)
			rl.mu.Lock()

			for ip, client := range rl.Clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(rl.Clients, ip)
				}
			}

			rl.mu.Unlock()
		}
	}

	return rl, cleanUp
}

func (rl *RateLimiter) RateLimit(app *application, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		if _, found := rl.Clients[ip]; !found {
			rl.Clients[ip] = &Client{limiter: rate.NewLimiter(2, 4)}
		}

		rl.Clients[ip].mu.Lock()
		rl.Clients[ip].lastSeen = time.Now()
		rl.Clients[ip].mu.Unlock()

		if !rl.Clients[ip].limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
