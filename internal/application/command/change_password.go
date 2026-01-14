package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type ChangePasswordCommand struct {
	UserID          string
	CurrentPassword string
	NewPassword     string
}

type ChangePasswordHandler interface {
	Handle(ctx context.Context, cmd ChangePasswordCommand) error
}

type changePasswordHandler struct {
	repo           adminuser.Repository
	passwordHasher PasswordHasher
}

func NewChangePasswordHandler(
	repo adminuser.Repository,
	passwordHasher PasswordHasher,
) ChangePasswordHandler {
	return &changePasswordHandler{
		repo:           repo,
		passwordHasher: passwordHasher,
	}
}

func (h *changePasswordHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) error {
	user, err := h.repo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	if !h.passwordHasher.Compare(user.PasswordHash, cmd.CurrentPassword) {
		return adminuser.ErrInvalidCredentials
	}

	newHash, err := h.passwordHasher.Hash(cmd.NewPassword)
	if err != nil {
		return err
	}

	user.ChangePassword(newHash)

	_, err = h.repo.Update(ctx, user)
	return err
}
