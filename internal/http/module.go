package http

import (
	"net/http"

	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/Sokol111/ecommerce-auth-service-api/gen/httpapi"
)

func NewHttpHandlerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			newRateLimitConfig,
			newAuthHandler,
			newSecurityHandler,
			newOgenServer,
			newLoginRateLimiter,
		),
		fx.Invoke(registerOgenRoutes),
	)
}

// newLoginRateLimiter creates rate limiter from config
func newLoginRateLimiter(cfg RateLimitConfig) *LoginRateLimiter {
	return NewLoginRateLimiter(cfg.LoginTokens, cfg.LoginInterval)
}

func newOgenServer(
	handler httpapi.Handler,
	securityHandler httpapi.SecurityHandler,
	tracerProvider trace.TracerProvider,
	meterProvider metric.MeterProvider,
	middlewares []middleware.Middleware,
	errorHandler ogenerrors.ErrorHandler,
) (*httpapi.Server, error) {
	return httpapi.NewServer(
		handler,
		securityHandler,
		httpapi.WithTracerProvider(tracerProvider),
		httpapi.WithMeterProvider(meterProvider),
		httpapi.WithErrorHandler(errorHandler),
		httpapi.WithMiddleware(middlewares...),
	)
}

func registerOgenRoutes(mux *http.ServeMux, server *httpapi.Server) {
	mux.Handle("/", server)
}
