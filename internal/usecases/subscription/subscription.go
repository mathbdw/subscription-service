package subscription

import (
	"context"
	"github.com/mathbdw/subscription-service/internal/domain/entities"
	"github.com/mathbdw/subscription-service/internal/errors"
	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
	"github.com/mathbdw/subscription-service/internal/interfaces/repositories"
)

type SubscriptionUsecase struct {
	repo   repositories.SubscriptionRepository
	logger observability.Logger
}

// NewSubscriptionUsecase - Constructor SubscriptionUsecase
func NewSubscriptionUsecase(repo repositories.SubscriptionRepository, logger observability.Logger) SubscriptionUsecase {
	return SubscriptionUsecase{repo: repo, logger: logger}
}

// Create - Adds new subscription
func (uc *SubscriptionUsecase) Create(ctx context.Context, sub entities.Subscription) error {
	err := uc.repo.Create(ctx, sub)
	if err != nil {
		return errors.Wrap(err, "SubscriptionUsecase.Create: repo exec")
	}

	return nil
}

// GetByID - Returns subscription by ID
func (uc *SubscriptionUsecase) GetByID(ctx context.Context, id int64) (*entities.Subscription, error) {
	sub, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "SubscriptionUsecase.GetByID: repo exec")
	}

	return sub, nil
}

// List - Returns slice subscriptions by Query Criteria
func (uc *SubscriptionUsecase) List(ctx context.Context, params entities.QueryCriteria) (*entities.ResponseListSubscription, error) {
	resp, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "SubscriptionUsecase.List: repo exec")
	}

	return resp, nil
}

// Update - Updated fields of subscription by ID
func (uc *SubscriptionUsecase) Update(ctx context.Context, id int64, fields map[string]any) error {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "SubscriptionUsecase.Update: repo getById")
	}

	err = uc.repo.Update(ctx, id, fields)
	if err != nil {
		return errors.Wrap(err, "SubscriptionUsecase.Update: repo exec")
	}

	return nil
}

// Delete - Deleted subscription by ID
func (uc *SubscriptionUsecase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "SubscriptionUsecase.Delete: repo getById")
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "SubscriptionUsecase.Delete: repo exec")
	}

	return nil
}

// GetCost - Returns total cost of subscriptions by FilterParams
func (uc *SubscriptionUsecase) GetCost(ctx context.Context, params entities.FilterParams) (int64, error) {
	cost, err := uc.repo.GetCost(ctx, params)
	if err != nil {
		return 0, errors.Wrap(err, "SubscriptionUsecase.GetCost: repo exec")
	}

	return cost, nil
}
