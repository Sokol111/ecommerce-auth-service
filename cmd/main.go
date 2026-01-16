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
	"github.com/Sokol111/ecommerce-commons/pkg/modules"
	"github.com/Sokol111/ecommerce-commons/pkg/swaggerui"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var AppModules = fx.Options(
	// Infrastructure - Core
	modules.NewCoreModule(),
	modules.NewPersistenceModule(),
	modules.NewHTTPModule(),
	modules.NewObservabilityModule(),
	modules.NewMessagingModule(),
	swaggerui.NewSwaggerModule(swaggerui.SwaggerConfig{OpenAPIContent: httpapi.OpenAPIDoc}),

	// Application
	application.Module(),

	// Infrastructure
	permissions.Module(),
	security.Module(),
	mongo.Module(),
	seeder.Module(),

	// HTTP
	http.NewHttpHandlerModule(),
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
