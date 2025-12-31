package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
	errs "github.com/mathbdw/subscription-service/internal/errors"
	"github.com/mathbdw/subscription-service/internal/interfaces/http/handlers/api/v1/convert"
	"github.com/mathbdw/subscription-service/internal/interfaces/http/handlers/api/v1/dto"
	"github.com/mathbdw/subscription-service/internal/interfaces/http/handlers/api/v1/response"
	"github.com/mathbdw/subscription-service/internal/interfaces/http/middleware"
	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
	uc "github.com/mathbdw/subscription-service/internal/usecases/subscription"
)

type HandlerSubscription struct {
	r         fiber.Router
	validator *validator.Validate
	uc        uc.SubscriptionUsecase

	logger observability.Logger
}

func NewHandler(apiV1Group fiber.Router, validator *validator.Validate, uc uc.SubscriptionUsecase, logger observability.Logger) {
	router := HandlerSubscription{
		uc:        uc,
		validator: validator,
		logger:    logger,
	}

	subscriptionGroup := apiV1Group.Group("/subscription")
	{
		subscriptionGroup.Post("/create", router.create)
		subscriptionGroup.Get("/list", middleware.ValidatedQueryParamsMiddleware(logger), router.list)
		subscriptionGroup.Get("/cost", middleware.ValidatedQueryParamsCostMiddleware(logger), router.cost)

		subscriptionGroup.Get("/:id", middleware.ValidatedQueryIdMiddleware(logger), router.getId)
		subscriptionGroup.Delete("/:id", middleware.ValidatedQueryIdMiddleware(logger), router.delete)
		subscriptionGroup.Patch("/:id", middleware.ValidatedQueryIdMiddleware(logger), router.update)
	}

}

// @Summary     Create subscription
// @Description Create new subscription
// @ID          SubscriptionCreate
// @Tags  	    Subscription
// @Accept      json
// @Produce     json
// @Param       request body dto.SubscriptionReq true "Data subscription"
// @Success     201
// @Failure     400 {object} response.Error
// @Failure     422 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /subscription/create [post]
func (h *HandlerSubscription) create(ctx *fiber.Ctx) error {
	var body dto.SubscriptionReq
	if err := ctx.BodyParser(&body); err != nil {
		h.logger.Error("subscriptionV1.Create: parse body", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := h.validator.Struct(body); err != nil {
		h.logger.Error("subscriptionV1.Create: validate struct", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	sub, err := convert.SubscriptionRequestToEntity(body)
	if err != nil {
		h.logger.Error("subscriptionV1.Create: convert", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := h.uc.Create(ctx.UserContext(), sub); err != nil {
		h.logger.Error("subscriptionV1.Create: usecase exec", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
	}

	return ctx.SendStatus(http.StatusCreated)
}

// @Summary     get subscription by ID
// @Description Returns subscription by ID
// @ID          SubscriptionGetByID
// @Tags  	    Subscription
// @Accept      json
// @Produce     json
// @Param       id   path      int  true  "Subscription ID"
// @Success     204
// @Failure     400 {object} response.Error
// @Failure     422 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /subscription/{id} [get]
func (h *HandlerSubscription) getId(ctx *fiber.Ctx) error {
	subID, ok := ctx.Locals("query_id").(int64)
	if !ok {
		h.logger.Error("subscriptionV1.GetId: get query_id", nil)

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "invalid request body")
	}

	sub, err := h.uc.GetByID(ctx.UserContext(), subID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.logger.Error("subscriptionV1.GetId: not found row", map[string]any{"id": subID})

			return response.ErrorResponse(ctx, http.StatusNotFound, "Not found")
		}
		h.logger.Error("subscriptionV1.GetId: usecase exec", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
	}

	subResp := convert.SubscriptionEntityToResponse(*sub)
	return ctx.Status(http.StatusOK).JSON(subResp)
}

// @Summary     get list subscriptions
// @Description Returns list subscriptions
// @ID          SubscriptionList
// @Tags  	    Subscription
// @Accept      json
// @Produce     json
// @Param       query query dto.QueryParamList true "Query Criteria"
// @Success     204
// @Failure     400 {object} response.Error
// @Failure     422 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /subscription/list [get]
func (h *HandlerSubscription) list(ctx *fiber.Ctx) error {
	params, ok := ctx.Locals("query_params").(dto.QueryParamList)
	if !ok {
		h.logger.Error("subscriptionV1.List: get query_params", nil)

		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	queryCriteria, err := convert.SubscriptionQueryParamsToQueryCriteria(params)
	if err != nil {
		h.logger.Error("subscriptionV1.List: convert", map[string]any{"err": err})

		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	pespList, err := h.uc.List(ctx.UserContext(), *queryCriteria)
	if err != nil {
		h.logger.Error("subscriptionV1.List: usecase exec", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
	}

	subsResp := convert.SubscriptionListToResponse(pespList.Data)
	addPaginationHeaders(ctx, pespList.Info)

	return ctx.Status(http.StatusOK).JSON(subsResp)
}

// @Summary     delete subscription by ID
// @Description delete subscription by ID
// @ID          SubscriptionDelete
// @Tags  	    Subscription
// @Accept      json
// @Produce     json
// @Success     204
// @Param       id   path      int  true  "Subscription ID"
// @Failure     400 {object} response.Error
// @Failure     422 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /subscription/{id} [delete]
func (h *HandlerSubscription) delete(ctx *fiber.Ctx) error {
	subID, ok := ctx.Locals("query_id").(int64)
	if !ok {
		h.logger.Error("subscriptionV1.GetId: get query_id", nil)

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "invalid request body")
	}

	err := h.uc.Delete(ctx.UserContext(), subID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.logger.Error("subscriptionV1.GetId: not found row", map[string]any{"id": subID})

			return response.ErrorResponse(ctx, http.StatusNotFound, "Not found")
		}
		h.logger.Error("subscriptionV1.Delete: usecase exec", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
	}

	return ctx.SendStatus(http.StatusNoContent)
}

// @Summary     update subscription by ID
// @Description update subscription by ID
// @ID          SubscriptionUpdate
// @Tags  	    Subscription
// @Accept      json
// @Produce     json
// @Param       id   path      int  true  "Subscription ID"
// @Param       request body dto.SubscriptionUpdateReq true "Data subscription"
// @Success     200
// @Failure     400 {object} response.Error
// @Failure     422 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /subscription/{id} [patch]
func (h *HandlerSubscription) update(ctx *fiber.Ctx) error {
	subID, ok := ctx.Locals("query_id").(int64)
	if !ok {
		h.logger.Error("subscriptionV1.GetId: get query_id", nil)

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "invalid request body")
	}

	var body dto.SubscriptionUpdateReq
	if err := ctx.BodyParser(&body); err != nil {
		h.logger.Error("subscriptionV1.Update: parse body", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := h.validator.Struct(body); err != nil {
		h.logger.Error("subscriptionV1.Update: validate struct", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	fieldsMap, err := convert.SubscriptionRequestToMap(body)
	if err != nil {
		h.logger.Error("subscriptionV1.Update: convert", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := h.uc.Update(ctx.UserContext(), subID, fieldsMap); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.logger.Error("subscriptionV1.GetId: not found row", map[string]any{"id": subID})

			return response.ErrorResponse(ctx, http.StatusNotFound, "Not found")
		}
		h.logger.Error("subscriptionV1.Update: usecase exec", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
	}

	return ctx.SendStatus(http.StatusOK)
}

// @Summary     get cost subscriptions
// @Description Returns cost subscriptions
// @ID          SubscriptionCost
// @Tags  	    Subscription
// @Accept      json
// @Produce     json
// @Param       query query dto.QueryParamCost true "Filter params"
// @Success     204
// @Failure     400 {object} response.Error
// @Failure     422 {object} response.Error
// @Failure     500 {object} response.Error
// @Router      /subscription/cost [get]
func (h *HandlerSubscription) cost(ctx *fiber.Ctx) error {
	params, ok := ctx.Locals("query_cost").(dto.QueryParamCost)
	if !ok {
		h.logger.Error("subscriptionV1.Cost: get query_cost", nil)

		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	filter, err := convert.SubscriptionQueryParamsCostToFilterParam(params)
	if err != nil {
		h.logger.Error("subscriptionV1.Cost: convert", map[string]any{"err": err})

		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	cost, err := h.uc.GetCost(ctx.UserContext(), filter)
	if err != nil {
		h.logger.Error("subscriptionV1.Cost: usecase exec", map[string]any{"err": err})

		return response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(cost)
}

// addPaginationHeaders - sets response headers pagination params for list
func addPaginationHeaders(ctx *fiber.Ctx, info entities.PaginationInfo) {
	ctx.Set("X-Page", strconv.FormatUint(info.Page, 10))
	ctx.Set("X-Page-Size", strconv.FormatUint(uint64(info.PageSize), 10))
	ctx.Set("X-Total-Count", strconv.FormatUint(info.TotalCount, 10))
	ctx.Set("X-Total-Pages", strconv.FormatUint(uint64(info.TotalPages), 10))

	hasNext := info.Page < uint64(info.TotalPages)
	hasPrev := info.Page > 1

	ctx.Set("X-Has-Next-Page", strconv.FormatBool(hasNext))
	ctx.Set("X-Has-Prev-Page", strconv.FormatBool(hasPrev))
}
