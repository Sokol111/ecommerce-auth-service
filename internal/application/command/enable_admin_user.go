package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type EnableAdminUserCommand struct {
	ID string
}

type EnableAdminUserHandler interface {
	Handle(ctx context.Context, cmd EnableAdminUserCommand) error
}

type enableAdminUserHandler struct {
	repo adminuser.Repository
}

func NewEnableAdminUserHandler(repo adminuser.Repository) EnableAdminUserHandler {
	return &enableAdminUserHandler{repo: repo}
}

func (h *enableAdminUserHandler) Handle(ctx context.Context, cmd EnableAdminUserCommand) error {
	user, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	user.Enable()

	_, err = h.repo.Update(ctx, user)
	return err
}
