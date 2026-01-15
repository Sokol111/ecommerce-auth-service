package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/security/token"
)

type RefreshTokenCommand struct {
	RefreshToken string
}

type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type RefreshTokenHandler interface {
	Handle(ctx context.Context, cmd RefreshTokenCommand) (*RefreshTokenResult, error)
}

type refreshTokenHandler struct {
	repo           adminuser.Repository
	tokenGenerator TokenGenerator
	tokenValidator token.TokenValidator
}

func NewRefreshTokenHandler(
	repo adminuser.Repository,
	tokenGenerator TokenGenerator,
	tokenValidator token.TokenValidator,
) RefreshTokenHandler {
	return &refreshTokenHandler{
		repo:           repo,
		tokenGenerator: tokenGenerator,
		tokenValidator: tokenValidator,
	}
}

func (h *refreshTokenHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*RefreshTokenResult, error) {
	claims, err := h.tokenValidator.ValidateToken(cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Ensure it's a refresh token
	if !claims.IsRefresh() {
		return nil, token.ErrInvalidToken
	}

	user, err := h.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, adminuser.ErrAdminUserDisabled
	}

	accessToken, _, expiresIn, err := h.tokenGenerator.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: cmd.RefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
