package response

import "github.com/gofiber/fiber/v2"

type Error struct {
	Error string `json:"error" example:"message"`
}

func ErrorResponse(ctx *fiber.Ctx, code int, msg string) error {
	return ctx.Status(code).JSON(Error{Error: msg})
}
