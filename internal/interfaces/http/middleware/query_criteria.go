package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
	"github.com/mathbdw/subscription-service/internal/interfaces/http/handlers/api/v1/dto"
	"github.com/mathbdw/subscription-service/internal/interfaces/http/handlers/api/v1/response"
	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
)

var validate = validator.New()

func init() {
	_ = validate.RegisterValidation("sort_by", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()

		if _, ok := entities.SortByTypes[val]; !ok {
			return false
		}

		return true
	})
	_ = validate.RegisterValidation("sort_order", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		val = strings.ToUpper(val)

		if _, ok := entities.SortOrderTypes[val]; !ok {
			return false
		}

		return true
	})
	_ = validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()

		if _, err := uuid.Parse(val); err != nil {
			return false
		}

		return true
	})

}

// ValidatedQueryParamsMiddleware - middleware parse and validate params query for List
func ValidatedQueryParamsMiddleware(logger observability.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var queryParams dto.QueryParamList

		if err := ctx.QueryParser(&queryParams); err != nil {
			logger.Error("middaleware.ValidatedQueryParamsMiddleware: invalid query parameters", map[string]any{"err": err})

			return response.ErrorResponse(ctx, http.StatusBadRequest, "Invalid query parameters")
		}

		if err := validate.Struct(&queryParams); err != nil {
			logger.Error("middaleware.ValidatedQueryParamsMiddleware: validate", map[string]any{"err": err})

			return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
		}

		ctx.Locals("query_params", queryParams)
		return ctx.Next()
	}
}

// ValidatedQueryParamsCostMiddleware - middleware parse and validate params query for Goat
func ValidatedQueryParamsCostMiddleware(logger observability.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var queryParams dto.QueryParamCost

		if err := ctx.QueryParser(&queryParams); err != nil {
			logger.Error("middaleware.ValidatedQueryParamsCostMiddleware: invalid query parameters", map[string]any{"err": err})

			return response.ErrorResponse(ctx, http.StatusBadRequest, "Invalid query parameters")
		}

		if err := validate.Struct(&queryParams); err != nil {
			logger.Error("middaleware.ValidatedQueryParamsCostMiddleware: validate", map[string]any{"err": err})

			return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
		}

		ctx.Locals("query_cost", queryParams)
		return ctx.Next()
	}
}

// ValidatedQueryIdMiddleware - middleware parse and validate params query ID subscription
func ValidatedQueryIdMiddleware(logger observability.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		subStrID := ctx.Params("id")
		subID, err := strconv.ParseInt(subStrID, 10, 64)
		if err != nil {
			logger.Error("middaleware.ValidatedQueryIdMiddleware: parse param", map[string]any{"subStrID": subStrID, "err": err})

			return response.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
		}

		if subID < 1 {
			logger.Error("middaleware.ValidatedQueryIdMiddleware: validate params", map[string]any{"subID": subID})

			return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, "Invalid id: must be a positive integer")
		}

		ctx.Locals("query_id", subID)
		return ctx.Next()
	}
}
