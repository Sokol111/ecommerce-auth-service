package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type DisableAdminUserCommand struct {
	ID            string
	RequestUserID string // ID of the user making the request
}

type DisableAdminUserHandler interface {
	Handle(ctx context.Context, cmd DisableAdminUserCommand) error
}

type disableAdminUserHandler struct {
	repo adminuser.Repository
}

func NewDisableAdminUserHandler(repo adminuser.Repository) DisableAdminUserHandler {
	return &disableAdminUserHandler{repo: repo}
}

func (h *disableAdminUserHandler) Handle(ctx context.Context, cmd DisableAdminUserCommand) error {
	// Cannot disable yourself
	if cmd.ID == cmd.RequestUserID {
		return adminuser.ErrCannotDisableSelf
	}

	user, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// Cannot disable super admin
	if user.IsSuperAdmin() {
		return adminuser.ErrCannotDisableSuperAdmin
	}

	user.Disable()

	_, err = h.repo.Update(ctx, user)
	return err
}
