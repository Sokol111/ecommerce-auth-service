package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
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
	repo         adminuser.Repository
	tokenService TokenService
}

func NewRefreshTokenHandler(
	repo adminuser.Repository,
	tokenService TokenService,
) RefreshTokenHandler {
	return &refreshTokenHandler{
		repo:         repo,
		tokenService: tokenService,
	}
}

func (h *refreshTokenHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*RefreshTokenResult, error) {
	claims, err := h.tokenService.ValidateRefreshToken(cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	user, err := h.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		return nil, adminuser.ErrAdminUserDisabled
	}

	accessToken, refreshToken, expiresIn, err := h.tokenService.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
