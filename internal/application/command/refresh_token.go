package command

import (
	"context"
	"errors"
	"time"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"github.com/Sokol111/ecommerce-commons/pkg/security/token"
	"go.uber.org/zap"
)

// ErrRefreshTokenReused indicates that a refresh token was used after it was already rotated
var ErrRefreshTokenReused = errors.New("refresh token reused: possible token theft detected")

type RefreshTokenCommand struct {
	RefreshToken string
}

type RefreshTokenResult struct {
	AccessToken      string
	RefreshToken     string
	ExpiresIn        int
	ExpiresAt        time.Time
	RefreshExpiresIn int
	RefreshExpiresAt time.Time
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
	log := logger.Get(ctx)

	claims, err := h.tokenValidator.ValidateToken(cmd.RefreshToken)
	if err != nil {
		return nil, err
	}

	if !claims.IsRefresh() {
		return nil, token.ErrInvalidToken
	}

	log = log.With(zap.String("user_id", claims.UserID))

	user, err := h.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.Enabled {
		log.Warn("token refresh failed: account disabled")
		return nil, adminuser.ErrAdminUserDisabled
	}

	tokenID, _ := claims.GetString("jti")
	if tokenID == "" || tokenID != user.RefreshTokenID {
		log.Error("SECURITY: refresh token reuse detected - invalidating all sessions",
			zap.String("presented_token_id", tokenID),
			zap.String("expected_token_id", user.RefreshTokenID),
		)
		user.SetRefreshTokenID("")
		_, _ = h.repo.Update(ctx, user)
		return nil, ErrRefreshTokenReused
	}

	tokens, err := h.tokenGenerator.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	user.SetRefreshTokenID(tokens.RefreshTokenID)
	if _, err := h.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	log.Debug("token refreshed successfully")

	return &RefreshTokenResult{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		ExpiresIn:        tokens.ExpiresIn,
		ExpiresAt:        tokens.ExpiresAt,
		RefreshExpiresIn: tokens.RefreshExpiresIn,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
	}, nil
}
