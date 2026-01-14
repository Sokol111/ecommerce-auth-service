package mongo

import (
	"time"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type adminUserMapper struct{}

func newAdminUserMapper() *adminUserMapper {
	return &adminUserMapper{}
}

func (m *adminUserMapper) ToEntity(u *adminuser.AdminUser) *adminUserEntity {
	return &adminUserEntity{
		ID:           u.ID,
		Version:      u.Version,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Role:         u.Role,
		Enabled:      u.Enabled,
		CreatedAt:    u.CreatedAt,
		ModifiedAt:   u.ModifiedAt,
		LastLoginAt:  u.LastLoginAt,
	}
}

func (m *adminUserMapper) ToDomain(e *adminUserEntity) *adminuser.AdminUser {
	var lastLoginAt *time.Time
	if e.LastLoginAt != nil {
		t := e.LastLoginAt.UTC()
		lastLoginAt = &t
	}

	return adminuser.Reconstruct(
		e.ID,
		e.Version,
		e.Email,
		e.PasswordHash,
		e.FirstName,
		e.LastName,
		e.Role,
		e.Enabled,
		e.CreatedAt.UTC(),
		e.ModifiedAt.UTC(),
		lastLoginAt,
	)
}

func (m *adminUserMapper) GetID(e *adminUserEntity) string {
	return e.ID
}

func (m *adminUserMapper) GetVersion(e *adminUserEntity) int {
	return e.Version
}

func (m *adminUserMapper) SetVersion(e *adminUserEntity, version int) {
	e.Version = version
}
