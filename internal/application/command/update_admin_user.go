package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type UpdateAdminUserCommand struct {
	ID        string
	FirstName *string
	LastName  *string
	Role      *adminuser.Role
}

type UpdateAdminUserHandler interface {
	Handle(ctx context.Context, cmd UpdateAdminUserCommand) (*adminuser.AdminUser, error)
}

type updateAdminUserHandler struct {
	repo adminuser.Repository
}

func NewUpdateAdminUserHandler(repo adminuser.Repository) UpdateAdminUserHandler {
	return &updateAdminUserHandler{repo: repo}
}

func (h *updateAdminUserHandler) Handle(ctx context.Context, cmd UpdateAdminUserCommand) (*adminuser.AdminUser, error) {
	user, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	firstName := user.FirstName
	if cmd.FirstName != nil {
		firstName = *cmd.FirstName
	}

	lastName := user.LastName
	if cmd.LastName != nil {
		lastName = *cmd.LastName
	}

	role := user.Role
	if cmd.Role != nil {
		role = *cmd.Role
	}

	user.Update(firstName, lastName, role)

	return h.repo.Update(ctx, user)
}
