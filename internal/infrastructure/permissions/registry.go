package permissions

import (
	"fmt"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

// registry implements adminuser.RolePermissionProvider
type registry struct {
	rolePermissions  map[adminuser.Role][]adminuser.Permission
	roleDescriptions map[adminuser.Role]string
}

func newRegistry(cfg Config) (adminuser.RolePermissionProvider, error) {
	if len(cfg.Roles) == 0 {
		return nil, fmt.Errorf("at least one role must be configured")
	}

	r := &registry{
		rolePermissions:  make(map[adminuser.Role][]adminuser.Permission),
		roleDescriptions: make(map[adminuser.Role]string),
	}

	for _, role := range cfg.Roles {
		perms := make([]adminuser.Permission, len(role.Permissions))
		for i, p := range role.Permissions {
			perms[i] = adminuser.Permission(p)
		}
		r.rolePermissions[adminuser.Role(role.Name)] = perms
		r.roleDescriptions[adminuser.Role(role.Name)] = role.Description
	}

	return r, nil
}

func (r *registry) GetPermissionsForRole(role adminuser.Role) []adminuser.Permission {
	return r.rolePermissions[role]
}

func (r *registry) GetValidRoles() []adminuser.Role {
	roles := make([]adminuser.Role, 0, len(r.rolePermissions))
	for role := range r.rolePermissions {
		roles = append(roles, role)
	}
	return roles
}

func (r *registry) IsValidRole(role adminuser.Role) bool {
	_, exists := r.rolePermissions[role]
	return exists
}

func (r *registry) GetRoleDescription(role adminuser.Role) string {
	return r.roleDescriptions[role]
}
