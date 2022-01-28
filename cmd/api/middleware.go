package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/khatibomar/gogive/internal/data"
	"github.com/khatibomar/gogive/internal/validator"
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

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Authorization" header to the response. This indicates to any
		// caches that the response may vary based on the value of the Authorization
		// header in the request.
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// we expect the value of the Authorization header to be in the format
		// "Bearer <token>". We try to split this into its constituent parts, and if the
		// header isn't in the expected format we return a 401 Unauthorized response
		// using the invalidAuthenticationTokenResponse() helper (which we will create
		// in a moment).
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})
}

// Create a new requireAuthenticatedUser() middleware to check that a user is not
// anonymous.
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Checks that a user is both authenticated and activated.
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	// Rather than returning this http.HandlerFunc we assign it to the variable fn.
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		// Check that a user is activated.
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	// Wrap fn with the requireAuthenticatedUser() middleware before returning it.
	return app.requireAuthenticatedUser(fn)
}

// requireOneOfRole check if the user is one of the roles/labels specified
// match the user
func (app *application) requireOneOfRole(roles []string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		valid := false
		for _, role := range roles {
			switch role {
			case data.ADMIN_ROLE:
				if user.Role == data.ADMIN_ROLE {
					valid = true
					break
				}
			case data.ANALYTIC_ROLE:
				if user.Role == data.ANALYTIC_ROLE {
					valid = true
					break
				}
			case data.ITEM_CREATOR_ROLE:
				itemID, err := app.readIDParam(r)
				if err != nil {
					app.badRequestResponse(w, r, err)
					return
				}
				item, err := app.models.Items.Get(itemID)
				if err != nil {
					app.badRequestResponse(w, r, err)
					return
				}
				app.contextSetItem(r, item)
				if user.ID == item.CreatedBy {
					valid = true
					break
				}
			}
			if valid {
				break
			}
		}
		if valid == false {
			app.errorRequireAtLeastRole(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
