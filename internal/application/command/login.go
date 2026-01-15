package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type LoginCommand struct {
	Email    string
	Password string
}

type LoginResult struct {
	User         *adminuser.AdminUser
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
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
	user, err := h.repo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, adminuser.ErrInvalidCredentials
	}

	if !h.passwordHasher.Compare(user.PasswordHash, cmd.Password) {
		return nil, adminuser.ErrInvalidCredentials
	}

	if !user.Enabled {
		return nil, adminuser.ErrAdminUserDisabled
	}

	user.RecordLogin()
	if _, err := h.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	accessToken, refreshToken, expiresIn, err := h.tokenGenerator.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
