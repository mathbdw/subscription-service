package repositories

import (
	"context"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
)

//go:generate mockgen -destination=./../../../mocks/mock_subscription_repository.go -package=mocks -source=./subscription_repository.go

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription entities.Subscription) error
	GetByID(ctx context.Context, id int64) (*entities.Subscription, error)
	List(ctx context.Context, params entities.QueryCriteria) (*entities.ResponseListSubscription, error)
	Update(ctx context.Context, id int64, fields map[string]any) error
	Delete(ctx context.Context, id int64) error
	GetCost(ctx context.Context, params entities.FilterParams) (int64, error)
}
