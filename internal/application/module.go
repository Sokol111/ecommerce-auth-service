package application

import (
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"github.com/Sokol111/ecommerce-auth-service/internal/application/query"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			command.NewLoginHandler,
			command.NewLogoutHandler,
			command.NewCreateAdminUserHandler,
			command.NewDisableAdminUserHandler,
			command.NewEnableAdminUserHandler,
			command.NewRefreshTokenHandler,
		),
		fx.Provide(
			query.NewGetAdminUserByIDHandler,
			query.NewListAdminUsersHandler,
			query.NewGetRolesHandler,
		),
	)
}
