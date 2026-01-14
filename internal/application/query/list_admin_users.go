package query

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

type ListAdminUsersQuery struct {
	Page    int
	Size    int
	Role    *adminuser.Role
	Enabled *bool
	Search  *string
}

type ListAdminUsersHandler interface {
	Handle(ctx context.Context, query ListAdminUsersQuery) (*commonsmongo.PageResult[adminuser.AdminUser], error)
}

type listAdminUsersHandler struct {
	repo adminuser.Repository
}

func NewListAdminUsersHandler(repo adminuser.Repository) ListAdminUsersHandler {
	return &listAdminUsersHandler{repo: repo}
}

func (h *listAdminUsersHandler) Handle(ctx context.Context, query ListAdminUsersQuery) (*commonsmongo.PageResult[adminuser.AdminUser], error) {
	return h.repo.FindList(ctx, adminuser.ListQuery{
		Page:    query.Page,
		Size:    query.Size,
		Role:    query.Role,
		Enabled: query.Enabled,
		Search:  query.Search,
	})
}
