package seeder

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Module provides the seeder for initial admin creation
func Module() fx.Option {
	return fx.Options(
		fx.Provide(newConfig),
		fx.Provide(NewSeeder),
		fx.Invoke(runSeeder),
	)
}

func runSeeder(lc fx.Lifecycle, seeder *Seeder, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := seeder.EnsureInitialAdmin(ctx); err != nil {
				log.Error("Failed to seed initial admin", zap.Error(err))
				// Don't fail startup - just log the error
				// The admin can be created manually later
			}
			return nil
		},
	})
}
