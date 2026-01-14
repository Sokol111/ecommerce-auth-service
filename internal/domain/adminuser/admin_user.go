package adminuser

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Role represents an admin user role
type Role string

const (
	RoleSuperAdmin     Role = "super_admin"
	RoleCatalogManager Role = "catalog_manager"
	RoleViewer         Role = "viewer"
)

// ValidRoles returns all valid roles
func ValidRoles() []Role {
	return []Role{
		RoleSuperAdmin,
		RoleCatalogManager,
		RoleViewer,
	}
}

// IsValid checks if the role is valid
func (r Role) IsValid() bool {
	for _, valid := range ValidRoles() {
		if r == valid {
			return true
		}
	}
	return false
}

// Permission represents a system permission
type Permission string

const (
	// User management
	PermUsersRead   Permission = "users:read"
	PermUsersWrite  Permission = "users:write"
	PermUsersDelete Permission = "users:delete"

	// Products
	PermProductsRead   Permission = "products:read"
	PermProductsWrite  Permission = "products:write"
	PermProductsDelete Permission = "products:delete"

	// Categories
	PermCategoriesRead   Permission = "categories:read"
	PermCategoriesWrite  Permission = "categories:write"
	PermCategoriesDelete Permission = "categories:delete"

	// Attributes
	PermAttributesRead   Permission = "attributes:read"
	PermAttributesWrite  Permission = "attributes:write"
	PermAttributesDelete Permission = "attributes:delete"
)

// RolePermissions maps roles to their permissions
var RolePermissions = map[Role][]Permission{
	RoleSuperAdmin: {
		PermUsersRead, PermUsersWrite, PermUsersDelete,
		PermProductsRead, PermProductsWrite, PermProductsDelete,
		PermCategoriesRead, PermCategoriesWrite, PermCategoriesDelete,
		PermAttributesRead, PermAttributesWrite, PermAttributesDelete,
	},
	RoleCatalogManager: {
		PermProductsRead, PermProductsWrite, PermProductsDelete,
		PermCategoriesRead, PermCategoriesWrite, PermCategoriesDelete,
		PermAttributesRead, PermAttributesWrite, PermAttributesDelete,
	},
	RoleViewer: {
		PermUsersRead,
		PermProductsRead,
		PermCategoriesRead,
		PermAttributesRead,
	},
}

// AdminUser - domain aggregate root for admin panel users
type AdminUser struct {
	ID           string
	Version      int
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         Role
	Enabled      bool
	CreatedAt    time.Time
	ModifiedAt   time.Time
	LastLoginAt  *time.Time
}

// NewAdminUser creates a new admin user
func NewAdminUser(email, passwordHash, firstName, lastName string, role Role) *AdminUser {
	now := time.Now().UTC()
	return &AdminUser{
		ID:           uuid.New().String(),
		Version:      1,
		Email:        normalizeEmail(email),
		PasswordHash: passwordHash,
		FirstName:    strings.TrimSpace(firstName),
		LastName:     strings.TrimSpace(lastName),
		Role:         role,
		Enabled:      true,
		CreatedAt:    now,
		ModifiedAt:   now,
		LastLoginAt:  nil,
	}
}

// Reconstruct rebuilds an admin user from persistence (no validation)
func Reconstruct(
	id string,
	version int,
	email, passwordHash, firstName, lastName string,
	role Role,
	enabled bool,
	createdAt, modifiedAt time.Time,
	lastLoginAt *time.Time,
) *AdminUser {
	return &AdminUser{
		ID:           id,
		Version:      version,
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         role,
		Enabled:      enabled,
		CreatedAt:    createdAt,
		ModifiedAt:   modifiedAt,
		LastLoginAt:  lastLoginAt,
	}
}

// Update modifies admin user data
func (u *AdminUser) Update(firstName, lastName string, role Role) {
	u.FirstName = strings.TrimSpace(firstName)
	u.LastName = strings.TrimSpace(lastName)
	u.Role = role
	u.ModifiedAt = time.Now().UTC()
}

// ChangePassword updates the password hash
func (u *AdminUser) ChangePassword(newPasswordHash string) {
	u.PasswordHash = newPasswordHash
	u.ModifiedAt = time.Now().UTC()
}

// RecordLogin updates the last login time
func (u *AdminUser) RecordLogin() {
	now := time.Now().UTC()
	u.LastLoginAt = &now
	u.ModifiedAt = now
}

// Disable disables the admin user
func (u *AdminUser) Disable() {
	u.Enabled = false
	u.ModifiedAt = time.Now().UTC()
}

// Enable enables the admin user
func (u *AdminUser) Enable() {
	u.Enabled = true
	u.ModifiedAt = time.Now().UTC()
}

// GetPermissions returns the permissions for this user based on their role
func (u *AdminUser) GetPermissions() []Permission {
	return RolePermissions[u.Role]
}

// HasPermission checks if the user has a specific permission
func (u *AdminUser) HasPermission(perm Permission) bool {
	for _, p := range u.GetPermissions() {
		if p == perm {
			return true
		}
	}
	return false
}

// IsSuperAdmin checks if the user is a super admin
func (u *AdminUser) IsSuperAdmin() bool {
	return u.Role == RoleSuperAdmin
}

// FullName returns the user's full name
func (u *AdminUser) FullName() string {
	return u.FirstName + " " + u.LastName
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
