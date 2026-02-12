package main

import (
	"context"

	"github.com/Sokol111/ecommerce-auth-service-api/gen/httpapi"
	"github.com/Sokol111/ecommerce-auth-service/internal/application"
	http "github.com/Sokol111/ecommerce-auth-service/internal/http"
	"github.com/Sokol111/ecommerce-auth-service/internal/infrastructure/permissions"
	"github.com/Sokol111/ecommerce-auth-service/internal/infrastructure/persistence/mongo"
	"github.com/Sokol111/ecommerce-auth-service/internal/infrastructure/security"
	"github.com/Sokol111/ecommerce-auth-service/internal/infrastructure/seeder"
	commons_core "github.com/Sokol111/ecommerce-commons/pkg/core"
	commons_http "github.com/Sokol111/ecommerce-commons/pkg/http"
	commons_observability "github.com/Sokol111/ecommerce-commons/pkg/observability"
	commons_persistence "github.com/Sokol111/ecommerce-commons/pkg/persistence"
	commons_swaggerui "github.com/Sokol111/ecommerce-commons/pkg/swaggerui"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppModules = fx.Options(
	// Commons
	commons_core.NewCoreModule(),
	commons_persistence.NewPersistenceModule(),
	commons_http.NewHTTPModule(),
	commons_observability.NewObservabilityModule(),
	commons_swaggerui.NewSwaggerModule(commons_swaggerui.SwaggerConfig{OpenAPIContent: httpapi.OpenAPIDoc}),

	// Application
	application.Module(),

	// Infrastructure
	permissions.Module(),
	security.Module(),
	mongo.Module(),
	seeder.Module(),

	// HTTP
	http.NewHTTPHandlerModule(),
)

func main() {
	app := fx.New(
		AppModules,
		fx.Invoke(func(lc fx.Lifecycle, log *zap.Logger) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					log.Info("Application stopping...")
					return nil
				},
			})
		}),
	)
	app.Run()
}
