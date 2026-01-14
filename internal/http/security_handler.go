package http

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
)

type bearerTokenKey struct{}

type securityHandler struct {
	tokenService command.TokenService
}

func newSecurityHandler(tokenService command.TokenService) httpapi.SecurityHandler {
	return &securityHandler{tokenService: tokenService}
}

// HandleBearerAuth handles BearerAuth security.
func (s *securityHandler) HandleBearerAuth(ctx context.Context, operationName httpapi.OperationName, t httpapi.BearerAuth) (context.Context, error) {
	_, err := s.tokenService.ValidateAccessToken(t.Token)
	if err != nil {
		return ctx, err
	}

	// Store token in context for handler to use
	return context.WithValue(ctx, bearerTokenKey{}, t.Token), nil
}
