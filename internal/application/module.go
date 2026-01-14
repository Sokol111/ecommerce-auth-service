package application

import (
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/query"
	"go.uber.org/fx"
)

// Module provides application layer dependencies
func Module() fx.Option {
	return fx.Options(
		// Command handlers
		fx.Provide(
			command.NewLoginHandler,
			command.NewChangePasswordHandler,
			command.NewCreateAdminUserHandler,
			command.NewUpdateAdminUserHandler,
			command.NewDisableAdminUserHandler,
			command.NewEnableAdminUserHandler,
			command.NewRefreshTokenHandler,
		),
		// Query handlers
		fx.Provide(
			query.NewGetAdminUserByIDHandler,
			query.NewListAdminUsersHandler,
			query.NewGetRolesHandler,
			query.NewGetPermissionsHandler,
		),
	)
}
