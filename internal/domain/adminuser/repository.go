package adminuser

import (
	"context"

	commonsmongo "github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

// ListQuery represents the query parameters for listing admin users
type ListQuery struct {
	Page    int
	Size    int
	Role    *Role
	Enabled *bool
	Search  *string // Search by name or email
}

// Repository defines the admin user persistence interface
type Repository interface {
	// Insert creates a new admin user
	Insert(ctx context.Context, user *AdminUser) error

	// FindByID retrieves an admin user by ID
	FindByID(ctx context.Context, id string) (*AdminUser, error)

	// FindByEmail retrieves an admin user by email
	FindByEmail(ctx context.Context, email string) (*AdminUser, error)

	// FindList retrieves a paginated list of admin users
	FindList(ctx context.Context, query ListQuery) (*commonsmongo.PageResult[AdminUser], error)

	// Update updates an existing admin user with optimistic locking
	Update(ctx context.Context, user *AdminUser) (*AdminUser, error)

	// ExistsByEmail checks if an admin user with given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}
