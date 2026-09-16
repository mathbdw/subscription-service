package state

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/mathbdw/subscription-service/config"
	"github.com/mathbdw/subscription-service/internal/usecases/subscription"
)

func Liveness() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func Readiness(uc subscription.Usecase, t time.Duration) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		reqCtx, cancel := context.WithTimeout(ctx.Context(), t)
		defer cancel()

		if err := uc.Check(reqCtx); err != nil {

			return ctx.SendStatus(fiber.StatusServiceUnavailable)
		}
		return ctx.SendStatus(fiber.StatusOK)
	}
}

func Version(cfg *config.Config) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		data := map[string]interface{}{
			"name":    cfg.Project.Name,
			"debug":   cfg.Project.Debug,
			"version": cfg.Project.Version,
		}

		return ctx.Status(fiber.StatusOK).JSON(data)
	}
}
