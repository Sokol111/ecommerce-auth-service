package command

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/core/logger"
	"go.uber.org/zap"
)

type LogoutCommand struct {
	UserID string
}

type LogoutHandler interface {
	Handle(ctx context.Context, cmd LogoutCommand) error
}

type logoutHandler struct {
	repo adminuser.Repository
}

func NewLogoutHandler(repo adminuser.Repository) LogoutHandler {
	return &logoutHandler{repo: repo}
}

func (h *logoutHandler) Handle(ctx context.Context, cmd LogoutCommand) error {
	user, err := h.repo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	user.SetRefreshTokenID("")

	if _, err := h.repo.Update(ctx, user); err != nil {
		return err
	}

	logger.Get(ctx).With(zap.String("user_id", cmd.UserID)).Info("user logged out")
	return nil
}
