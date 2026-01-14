package query

import (
	"context"
	"strings"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type PermissionInfo struct {
	Name        adminuser.Permission
	Description string
	Resource    string
}

type GetPermissionsHandler interface {
	Handle(ctx context.Context) []PermissionInfo
}

type getPermissionsHandler struct{}

func NewGetPermissionsHandler() GetPermissionsHandler {
	return &getPermissionsHandler{}
}

func (h *getPermissionsHandler) Handle(ctx context.Context) []PermissionInfo {
	permissions := []adminuser.Permission{
		adminuser.PermUsersRead, adminuser.PermUsersWrite, adminuser.PermUsersDelete,
		adminuser.PermProductsRead, adminuser.PermProductsWrite, adminuser.PermProductsDelete,
		adminuser.PermCategoriesRead, adminuser.PermCategoriesWrite, adminuser.PermCategoriesDelete,
		adminuser.PermAttributesRead, adminuser.PermAttributesWrite, adminuser.PermAttributesDelete,
	}

	result := make([]PermissionInfo, len(permissions))
	for i, p := range permissions {
		parts := strings.Split(string(p), ":")
		resource := parts[0]
		action := parts[1]

		result[i] = PermissionInfo{
			Name:        p,
			Description: strings.Title(action) + " " + resource,
			Resource:    resource,
		}
	}
	return result
}
