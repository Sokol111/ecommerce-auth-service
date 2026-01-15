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

type getRolesHandler struct {
	provider adminuser.RolePermissionProvider
}

func NewGetRolesHandler(provider adminuser.RolePermissionProvider) GetRolesHandler {
	return &getRolesHandler{provider: provider}
}

func (h *getRolesHandler) Handle(ctx context.Context) []RoleInfo {
	roles := h.provider.GetValidRoles()

	result := make([]RoleInfo, 0, len(roles))
	for _, role := range roles {
		description := h.provider.GetRoleDescription(role)
		if description == "" {
			description = string(role)
		}

		result = append(result, RoleInfo{
			Name:        role,
			Description: description,
			Permissions: h.provider.GetPermissionsForRole(role),
		})
	}

	return result
}
