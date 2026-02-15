package http //nolint:revive // package name intentional

import (
	"net/http"

	"go.uber.org/fx"

	"github.com/Sokol111/ecommerce-auth-service-api/gen/httpapi"
)

func NewHTTPHandlerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			newRateLimitConfig,
			newAuthHandler,
			newSecurityHandler,
			httpapi.ProvideServer,
			newLoginRateLimiter,
		),
		fx.Invoke(registerOgenRoutes),
	)
}

// newLoginRateLimiter creates rate limiter from config
func newLoginRateLimiter(cfg RateLimitConfig) *LoginRateLimiter {
	return NewLoginRateLimiter(cfg.LoginTokens, cfg.LoginInterval)
}

func registerOgenRoutes(mux *http.ServeMux, server *httpapi.Server) {
	mux.Handle("/", server)
}
