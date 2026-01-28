package command

import (
	"context"
	"time"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"go.uber.org/zap"
)

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	User             *adminuser.AdminUser
	AccessToken      string
	RefreshToken     string
	ExpiresIn        int
	ExpiresAt        time.Time
	RefreshExpiresIn int
	RefreshExpiresAt time.Time
}

type LoginHandler interface {
	Handle(ctx context.Context, cmd LoginCommand) (*LoginResult, error)
}

type loginHandler struct {
	repo           adminuser.Repository
	passwordHasher PasswordHasher
	tokenGenerator TokenGenerator
}

func NewLoginHandler(
	repo adminuser.Repository,
	passwordHasher PasswordHasher,
	tokenGenerator TokenGenerator,
) LoginHandler {
	return &loginHandler{
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

func (h *loginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	log := logger.Get(ctx).With(zap.String("email", cmd.Email))

	user, err := h.repo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		log.Warn("login failed: user not found")
		return nil, adminuser.ErrInvalidCredentials
	}

	if !h.passwordHasher.Compare(user.PasswordHash, cmd.Password) {
		log.Warn("login failed: invalid password", zap.String("user_id", user.ID))
		return nil, adminuser.ErrInvalidCredentials
	}

	if !user.Enabled {
		log.Warn("login failed: account disabled", zap.String("user_id", user.ID))
		return nil, adminuser.ErrAdminUserDisabled
	}

	tokens, err := h.tokenGenerator.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Store refresh token ID for rotation validation
	user.SetRefreshTokenID(tokens.RefreshTokenID)
	user.RecordLogin()
	if _, err := h.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	log.Info("login successful", zap.String("user_id", user.ID), zap.String("role", string(user.Role)))

	return &LoginResult{
		User:             user,
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		ExpiresIn:        tokens.ExpiresIn,
		ExpiresAt:        tokens.ExpiresAt,
		RefreshExpiresIn: tokens.RefreshExpiresIn,
		RefreshExpiresAt: tokens.RefreshExpiresAt,
	}, nil
}
