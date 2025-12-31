package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
)

func Logger(logger observability.Logger) func(c *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()

		logger.Info("middleware", map[string]any{
			"ip":          ctx.IP(),
			"method":      ctx.Method(),
			"url":         ctx.OriginalURL(),
			"status_code": ctx.Response().StatusCode(),
			"length":      len(ctx.Response().Body()),
		})

		return err
	}
}
