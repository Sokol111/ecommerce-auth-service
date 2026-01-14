package query

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type RoleInfo struct {
	Name        adminuser.Role
	Description string
	Permissions []adminuser.Permission
}

type GetRolesHandler interface {
	Handle(ctx context.Context) []RoleInfo
}

type getRolesHandler struct{}

func NewGetRolesHandler() GetRolesHandler {
	return &getRolesHandler{}
}

func (h *getRolesHandler) Handle(ctx context.Context) []RoleInfo {
	return []RoleInfo{
		{
			Name:        adminuser.RoleSuperAdmin,
			Description: "Full access to all resources",
			Permissions: adminuser.RolePermissions[adminuser.RoleSuperAdmin],
		},
		{
			Name:        adminuser.RoleCatalogManager,
			Description: "Can manage products, categories and attributes",
			Permissions: adminuser.RolePermissions[adminuser.RoleCatalogManager],
		},
		{
			Name:        adminuser.RoleViewer,
			Description: "Read-only access",
			Permissions: adminuser.RolePermissions[adminuser.RoleViewer],
		},
	}
}
