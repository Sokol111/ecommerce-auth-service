package adminuser

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
)

// Role represents an admin user role
type Role string

// Permission represents a system permission
type Permission string

// RolePermissionProvider provides role and permission information
type RolePermissionProvider interface {
	GetPermissionsForRole(role Role) []Permission
	GetValidRoles() []Role
	IsValidRole(role Role) bool
	GetRoleDescription(role Role) string
}

// AdminUser - domain aggregate root for admin panel users
type AdminUser struct {
	ID             string
	Version        int
	Email          string
	PasswordHash   string
	FirstName      string
	LastName       string
	Role           Role
	Enabled        bool
	RefreshTokenID string
	CreatedAt      time.Time
	ModifiedAt     time.Time
	LastLoginAt    *time.Time
}

// NewAdminUser creates a new admin user
func NewAdminUser(email, passwordHash, firstName, lastName string, role Role) *AdminUser {
	now := time.Now().UTC()
	return &AdminUser{
		ID:             uuid.New().String(),
		Version:        1,
		Email:          normalizeEmail(email),
		PasswordHash:   passwordHash,
		FirstName:      strings.TrimSpace(firstName),
		LastName:       strings.TrimSpace(lastName),
		Role:           role,
		Enabled:        true,
		RefreshTokenID: "",
		CreatedAt:      now,
		ModifiedAt:     now,
		LastLoginAt:    nil,
	}
}

// Reconstruct rebuilds an admin user from persistence (no validation)
func Reconstruct(
	id string,
	version int,
	email, passwordHash, firstName, lastName string,
	role Role,
	enabled bool,
	refreshTokenID string,
	createdAt, modifiedAt time.Time,
	lastLoginAt *time.Time,
) *AdminUser {
	return &AdminUser{
		ID:             id,
		Version:        version,
		Email:          email,
		PasswordHash:   passwordHash,
		FirstName:      firstName,
		LastName:       lastName,
		Role:           role,
		Enabled:        enabled,
		RefreshTokenID: refreshTokenID,
		CreatedAt:      createdAt,
		ModifiedAt:     modifiedAt,
		LastLoginAt:    lastLoginAt,
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

// SetRefreshTokenID updates the refresh token ID for token rotation
func (u *AdminUser) SetRefreshTokenID(tokenID string) {
	u.RefreshTokenID = tokenID
	u.ModifiedAt = time.Now().UTC()
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

// HasPermission checks if the user has a specific permission
func (u *AdminUser) HasPermission(permissions []Permission, perm Permission) bool {
	return lo.Contains(permissions, perm)
}

// IsSuperAdmin checks if the user is a super admin
func (u *AdminUser) IsSuperAdmin() bool {
	return u.Role == "super_admin"
}

// FullName returns the user's full name
func (u *AdminUser) FullName() string {
	return u.FirstName + " " + u.LastName
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
