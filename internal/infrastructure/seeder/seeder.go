package seeder

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"github.com/Sokol111/ecommerce-auth-service/internal/domain/adminuser"
	"github.com/Sokol111/ecommerce-commons/pkg/persistence/mongo"
)

const superAdminRole = adminuser.Role("super_admin")

// Seeder handles initial admin user creation
type Seeder struct {
	cfg            Config
	repo           adminuser.Repository
	passwordHasher command.PasswordHasher
	log            *zap.Logger
}

// NewSeeder creates a new Seeder instance
func NewSeeder(
	cfg Config,
	repo adminuser.Repository,
	passwordHasher command.PasswordHasher,
	log *zap.Logger,
) *Seeder {
	return &Seeder{
		cfg:            cfg,
		repo:           repo,
		passwordHasher: passwordHasher,
		log:            log.Named("seeder"),
	}
}

// EnsureInitialAdmin creates the initial super admin if none exists.
// This method is idempotent - it will not create duplicates.
func (s *Seeder) EnsureInitialAdmin(ctx context.Context) error {
	if !s.cfg.Enabled {
		s.log.Debug("Initial admin seeding is disabled")
		return nil
	}

	if err := s.validateConfig(); err != nil {
		s.log.Warn("Invalid initial admin configuration, skipping seed", zap.Error(err))
		return nil
	}

	// Check if admin with this email already exists
	existing, err := s.repo.FindByEmail(ctx, s.cfg.Email)
	if err != nil && !errors.Is(err, mongo.ErrEntityNotFound) {
		return fmt.Errorf("failed to check existing admin: %w", err)
	}

	if existing != nil {
		s.log.Info("Initial admin already exists, skipping seed",
			zap.String("email", s.cfg.Email),
			zap.String("role", string(existing.Role)),
		)
		return nil
	}

	// Create the initial super admin
	passwordHash, err := s.passwordHasher.Hash(s.cfg.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := adminuser.NewAdminUser(
		s.cfg.Email,
		passwordHash,
		s.cfg.FirstName,
		s.cfg.LastName,
		superAdminRole,
	)

	if err := s.repo.Insert(ctx, user); err != nil {
		if errors.Is(err, adminuser.ErrEmailAlreadyExists) {
			// Race condition - another instance created the user
			s.log.Info("Initial admin was created by another instance")
			return nil
		}
		return fmt.Errorf("failed to create initial admin: %w", err)
	}

	s.log.Info("Initial super admin created successfully",
		zap.String("email", s.cfg.Email),
		zap.String("firstName", s.cfg.FirstName),
		zap.String("lastName", s.cfg.LastName),
	)

	return nil
}

func (s *Seeder) validateConfig() error {
	if s.cfg.Email == "" {
		return errors.New("initial admin email is required")
	}
	if s.cfg.Password == "" {
		return errors.New("initial admin password is required")
	}
	return nil
}
