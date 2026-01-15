package security

import (
	"go.uber.org/fx"
)

// Module provides security infrastructure dependencies
func Module() fx.Option {
	return fx.Provide(
		newBcryptHasher,
		newPasetoService,
	)
}
