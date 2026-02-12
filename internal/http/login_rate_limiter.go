package http //nolint:revive // package name intentional

import (
	"context"
	"time"

	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/memorystore"
)

// LoginRateLimiter provides rate limiting per email for login attempts.
type LoginRateLimiter struct {
	store limiter.Store
}

// NewLoginRateLimiter creates a new rate limiter.
// tokens: allowed requests per interval
// interval: time window for the limit
// Example: NewLoginRateLimiter(5, time.Minute) = 5 requests per minute
func NewLoginRateLimiter(tokens uint64, interval time.Duration) *LoginRateLimiter {
	store, _ := memorystore.New(&memorystore.Config{ //nolint:errcheck // config is static
		Tokens:        tokens,
		Interval:      interval,
		SweepInterval: time.Minute,
		SweepMinTTL:   10 * time.Minute,
	})
	return &LoginRateLimiter{store: store}
}

// Allow checks if a login attempt for the given email is allowed.
func (l *LoginRateLimiter) Allow(ctx context.Context, email string) bool {
	_, _, _, ok, _ := l.store.Take(ctx, email) //nolint:errcheck // we only care about ok
	return ok
}
