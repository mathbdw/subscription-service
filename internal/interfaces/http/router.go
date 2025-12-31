package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/mathbdw/subscription-service/config"
	_ "github.com/mathbdw/subscription-service/docs/swagger" // Swagger docs.
	"github.com/mathbdw/subscription-service/internal/interfaces/http/handlers/api/v1"
	"github.com/mathbdw/subscription-service/internal/interfaces/http/middleware"
	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
	uc "github.com/mathbdw/subscription-service/internal/usecases/subscription"
)

// NewRouter -.
// Swagger spec:
// @title       Subscription API
// @description Subscription service API methods
// @version     1.0
// @host        localhost:8080
// @BasePath    /api/v1
func NewRouter(app *fiber.App, cfg *config.Rest, uc uc.SubscriptionUsecase, logger observability.Logger) {
	// Options
	app.Use(middleware.Logger(logger))
	app.Use(middleware.Recovery(logger))

	// Swagger
	if cfg.Swagger {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// Routers
	apiV1Group := app.Group("/api/v1")
	{
		v1.NewHandler(apiV1Group, validator.New(validator.WithRequiredStructEnabled()), uc, logger)
	}
}
