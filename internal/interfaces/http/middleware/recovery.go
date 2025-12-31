package middleware

import (
	"github.com/gofiber/fiber/v2"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
)

func logPanic(logger observability.Logger) func(c *fiber.Ctx, err any) {
	return func(ctx *fiber.Ctx, err any) {
		logger.Error("PANIC DETECTED", map[string]any{
			"ip":          ctx.IP(),
			"method":      ctx.Method(),
			"url":         ctx.OriginalURL(),
			"status_code": ctx.Response().StatusCode(),
			"length":      len(ctx.Response().Body()),
			"err":         err,
		})
	}
}

func Recovery(logger observability.Logger) func(c *fiber.Ctx) error {
	return fiberRecover.New(fiberRecover.Config{
		EnableStackTrace:  true,
		StackTraceHandler: logPanic(logger),
	})
}
