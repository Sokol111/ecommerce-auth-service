package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"go.uber.org/zap"
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
	log := logger.Get(ctx).With(zap.String("target_user_id", cmd.ID))

	user, err := h.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	user.Enable()

	_, err = h.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	log.Info("admin user enabled", zap.String("email", user.Email))
	return nil
}
