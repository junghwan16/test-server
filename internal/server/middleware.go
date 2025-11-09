package server

import (
	"log/slog"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/junghwan16/test-server/internal/identity/application"
	"github.com/junghwan16/test-server/internal/identity/handler"
)

// RequireAuth checks if the user is authenticated via session cookie
func RequireAuth(authSvc *application.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			_, u, err := authSvc.ValidateSession(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !u.Active() {
				http.Error(w, "Account deactivated", http.StatusForbidden)
				return
			}

			ctx := handler.SetUserInContext(r.Context(), u)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAdmin checks if the authenticated user has admin role
func RequireAdmin(authSvc *application.AuthService) func(http.Handler) http.Handler {
	authMiddleware := RequireAuth(authSvc)
	return func(next http.Handler) http.Handler {
		return authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			u := handler.GetUserFromContext(r.Context())
			if u == nil || !u.IsAdmin() {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}))
	}
}

// RateLimit implements simple in-memory rate limiting per IP
func RateLimit(rps float64, burst int) func(http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)

			mu.Lock()
			c, exists := clients[ip]
			if !exists {
				c = &client{limiter: rate.NewLimiter(rate.Limit(rps), burst)}
				clients[ip] = c
			}
			c.lastSeen = time.Now()
			mu.Unlock()

			if !c.limiter.Allow() {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Logging logs all requests
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"ip", getIP(r),
			)
			next.ServeHTTP(w, r)
		})
	}
}

func getIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
