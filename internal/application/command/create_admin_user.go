package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type CreateAdminUserCommand struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Role      adminuser.Role
}

type CreateAdminUserHandler interface {
	Handle(ctx context.Context, cmd CreateAdminUserCommand) (*adminuser.AdminUser, error)
}

type createAdminUserHandler struct {
	repo           adminuser.Repository
	passwordHasher PasswordHasher
}

func NewCreateAdminUserHandler(
	repo adminuser.Repository,
	passwordHasher PasswordHasher,
) CreateAdminUserHandler {
	return &createAdminUserHandler{
		repo:           repo,
		passwordHasher: passwordHasher,
	}
}

func (h *createAdminUserHandler) Handle(ctx context.Context, cmd CreateAdminUserCommand) (*adminuser.AdminUser, error) {
	exists, err := h.repo.ExistsByEmail(ctx, cmd.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, adminuser.ErrEmailAlreadyExists
	}

	passwordHash, err := h.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return nil, err
	}

	user := adminuser.NewAdminUser(cmd.Email, passwordHash, cmd.FirstName, cmd.LastName, cmd.Role)

	if err := h.repo.Insert(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
