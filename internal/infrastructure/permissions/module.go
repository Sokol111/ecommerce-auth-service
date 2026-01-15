package permissions

import "go.uber.org/fx"

// Module provides permission configuration loading
func Module() fx.Option {
	return fx.Provide(
		newConfig,
		newRegistry,
	)
}
