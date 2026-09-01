package subscription

import (
	"context"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
	"github.com/mathbdw/subscription-service/internal/errors"
	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
	"github.com/mathbdw/subscription-service/internal/interfaces/repositories"
)

type Usecase struct {
	repo   repositories.SubscriptionRepository
	logger observability.Logger
}

// NewUsecase - Constructor Usecase.
func NewUsecase(repo repositories.SubscriptionRepository, logger observability.Logger) Usecase {
	return Usecase{repo: repo, logger: logger}
}

// Create - Adds new subscription.
func (uc *Usecase) Create(ctx context.Context, sub entities.Subscription) error {
	err := uc.repo.Create(ctx, sub)
	if err != nil {
		return errors.Wrap(err, "Usecase.Create: repo exec")
	}

	return nil
}

// GetByID - Returns subscription by ID.
func (uc *Usecase) GetByID(ctx context.Context, id int64) (*entities.Subscription, error) {
	sub, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Usecase.GetByID: repo exec")
	}

	return sub, nil
}

// List - Returns slice subscriptions by Query Criteria.
func (uc *Usecase) List(ctx context.Context, params entities.QueryCriteria) (*entities.ResponseListSubscription, error) {
	resp, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, errors.Wrap(err, "Usecase.List: repo exec")
	}

	return resp, nil
}

// Update - Updated fields of subscription by ID.
func (uc *Usecase) Update(ctx context.Context, id int64, fields map[string]any) error {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Usecase.Update: repo getById")
	}

	err = uc.repo.Update(ctx, id, fields)
	if err != nil {
		return errors.Wrap(err, "Usecase.Update: repo exec")
	}

	return nil
}

// Delete - Deleted subscription by ID.
func (uc *Usecase) Delete(ctx context.Context, id int64) error {
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Usecase.Delete: repo getById")
	}

	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Usecase.Delete: repo exec")
	}

	return nil
}

// GetCost - Returns total cost of subscriptions by FilterParams.
func (uc *Usecase) GetCost(ctx context.Context, params entities.FilterParams) (int64, error) {
	cost, err := uc.repo.GetCost(ctx, params)
	if err != nil {
		return 0, errors.Wrap(err, "Usecase.GetCost: repo exec")
	}

	return cost, nil
}
