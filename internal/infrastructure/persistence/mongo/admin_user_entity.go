package mongo

import (
	"time"

	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
)

type adminUserEntity struct {
	ID           string         `bson:"_id"`
	Version      int            `bson:"version"`
	Email        string         `bson:"email"`
	PasswordHash string         `bson:"passwordHash"`
	FirstName    string         `bson:"firstName"`
	LastName     string         `bson:"lastName"`
	Role         adminuser.Role `bson:"role"`
	Enabled      bool           `bson:"enabled"`
	CreatedAt    time.Time      `bson:"createdAt"`
	ModifiedAt   time.Time      `bson:"modifiedAt"`
	LastLoginAt  *time.Time     `bson:"lastLoginAt,omitempty"`
}
