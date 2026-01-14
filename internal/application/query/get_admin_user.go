package query

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type GetAdminUserByIDQuery struct {
	ID string
}

type GetAdminUserByIDHandler interface {
	Handle(ctx context.Context, query GetAdminUserByIDQuery) (*adminuser.AdminUser, error)
}

type getAdminUserByIDHandler struct {
	repo adminuser.Repository
}

func NewGetAdminUserByIDHandler(repo adminuser.Repository) GetAdminUserByIDHandler {
	return &getAdminUserByIDHandler{repo: repo}
}

func (h *getAdminUserByIDHandler) Handle(ctx context.Context, query GetAdminUserByIDQuery) (*adminuser.AdminUser, error) {
	return h.repo.FindByID(ctx, query.ID)
}
