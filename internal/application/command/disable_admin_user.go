package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"go.uber.org/zap"
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
	log := logger.Get(ctx).With(
		zap.String("target_user_id", cmd.ID),
		zap.String("requested_by", cmd.RequestUserID),
	)

	if cmd.ID == cmd.RequestUserID {
		log.Warn("disable user failed: cannot disable self")
		return adminuser.ErrCannotDisableSelf
	}

	user, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if user.IsSuperAdmin() {
		log.Warn("disable user failed: cannot disable super admin")
		return adminuser.ErrCannotDisableSuperAdmin
	}

	user.Disable()

	_, err = h.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	log.Info("admin user disabled", zap.String("email", user.Email))
	return nil
}
