package adminuser

import "errors"

var (
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrAdminUserDisabled       = errors.New("admin user account is disabled")
	ErrInvalidAdminUserData    = errors.New("invalid admin user data")
	ErrCannotDisableSelf       = errors.New("cannot disable your own account")
	ErrCannotDisableSuperAdmin = errors.New("cannot disable super admin")
)
