package security

import (
	"github.com/Sokol111/ecommerce-auth-service/internal/application/command"
	"go.uber.org/fx"
)

// Module provides security infrastructure dependencies
func Module() fx.Option {
	return fx.Provide(
		NewBcryptHasher,
		fx.Annotate(
			func(h *BcryptHasher) command.PasswordHasher { return h },
			fx.As(new(command.PasswordHasher)),
		),
		fx.Annotate(
			func(config TokenConfig) (*PasetoService, error) {
				return NewPasetoService(config)
			},
			fx.As(new(command.TokenService)),
		),
	)
}
