package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
	errs "github.com/mathbdw/subscription-service/internal/errors"
	"github.com/mathbdw/subscription-service/internal/interfaces/observability"
	"github.com/mathbdw/subscription-service/internal/interfaces/repositories"
)

type subscriptionRepository struct {
	querier sqlx.ExtContext
	builder sq.StatementBuilderType

	logger observability.Logger
}

// NewUserRepository - Constructor UserRepository
func NewUserRepository(querier sqlx.ExtContext, builder sq.StatementBuilderType, logger observability.Logger) repositories.SubscriptionRepository {
	return &subscriptionRepository{
		querier: querier,
		builder: builder,

		logger: logger,
	}
}

var (
	table              = "subscription"
	columnsSelect      = []string{"id", "service_name", "user_id", "price", "start_date", "end_date"}
	columnsSelectCount = []string{"COUNT(*)"}
	columnsCost        = []string{"COALESCE(SUM(price), 0)"}
)

// Create - create new row
func (r *subscriptionRepository) Create(ctx context.Context, subs entities.Subscription) error {
	dataMap := SubscriptionToMap(subs)

	query, args, err := r.builder.Insert(table).SetMap(dataMap).ToSql()
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Create: build query")
	}

	res, err := r.querier.ExecContext(ctx, query, args...)
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Create: exec query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Create: get affected rows")
	}

	if rowsAffected != 1 {
		return fmt.Errorf("subscriptionRepositories.Create: expected rowsAffected %d", rowsAffected)
	}

	return nil
}

// GetByID - Returns subscription by ID
func (r *subscriptionRepository) GetByID(ctx context.Context, id int64) (*entities.Subscription, error) {
	query, args, err := r.builder.Select(columnsSelect...).
		From(table).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, errs.Wrap(err, "subscriptionRepositories.getByID: build query")
	}

	sub := &entities.Subscription{}

	err = r.querier.QueryRowxContext(ctx, query, args...).StructScan(sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, errs.Wrap(err, "subscriptionRepositories.getByID: scan query")
	}

	return sub, nil
}

// List - Returns a list of subscription using query criteria
func (r *subscriptionRepository) List(ctx context.Context, params entities.QueryCriteria) (*entities.ResponseListSubscription, error) {
	query := r.builder.Select(columnsSelectCount...).From(table)
	query = conditionList(query, params.Filter)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errs.Wrap(err, "subscriptionRepositories.List: build query count()")
	}

	var totalCount uint64
	err = r.querier.QueryRowxContext(ctx, sql, args...).Scan(&totalCount)
	if err != nil {
		return nil, errs.Wrap(err, "subscriptionRepositories.List: scan query")
	}

	limit := params.Pagination.Limit

	query = r.builder.Select(columnsSelect...).From(table)
	query = conditionList(query, params.Filter)
	query = paginationList(query, totalCount, &params.Pagination)
	query = sortList(query, params.Sort)

	sql, args, err = query.ToSql()
	if err != nil {
		return nil, errs.Wrap(err, "subscriptionRepositories.List: build query")
	}

	rows, err := r.querier.QueryxContext(ctx, sql, args...)
	if err != nil {
		return nil, errs.Wrap(err, "subscriptionRepositories.List: get query")
	}
	defer rows.Close()

	subs := make([]entities.Subscription, 0, params.Pagination.Limit)
	for rows.Next() {
		var sub entities.Subscription
		err = rows.StructScan(&sub)
		if err != nil {
			return nil, errs.Wrap(err, "subscriptionRepositories.List: scan query")
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.Wrap(err, "subscriptionRepositories.List: iteration rows")
	}

	totalPages := totalCount / limit
	if totalPages*limit < totalCount {
		totalPages++
	}

	return &entities.ResponseListSubscription{
		Data: subs,
		Info: entities.PaginationInfo{
			Page:       params.Pagination.Page,
			PageSize:   uint16(params.Pagination.Limit),
			TotalCount: totalCount,
			TotalPages: uint32(totalPages),
		},
	}, nil
}

// Update - Updated the fields
func (r *subscriptionRepository) Update(ctx context.Context, id int64, fields map[string]any) error {
	if err := validateUpdateFields(fields); err != nil {
		r.logger.Error("subscriptionRepositories.Update: validateUpdateFields", fields)

		return errs.Wrap(err, "subscriptionRepositories.Update: validate")
	}

	fields["updated_at"] = time.Now().UTC()

	query, args, err := r.builder.Update(table).
		Where(sq.Eq{"id": id}).
		SetMap(fields).
		ToSql()
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Update: build query")
	}

	res, err := r.querier.ExecContext(ctx, query, args...)
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Update: exec query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Update: get affected rows")
	}

	if rowsAffected != 1 {
		return fmt.Errorf("subscriptionRepositories.Update: expected rowsAffected %d", rowsAffected)
	}

	return nil
}

// Delete - Deleted row with the id
func (r *subscriptionRepository) Delete(ctx context.Context, id int64) error {
	query, args, err := r.builder.Delete(table).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Delete: build query")
	}

	res, err := r.querier.ExecContext(ctx, query, args...)
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Delete: exec query")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errs.Wrap(err, "subscriptionRepositories.Delete: get affected rows")
	}

	if rowsAffected != 1 {
		return fmt.Errorf("subscriptionRepositories.Delete: expected rowsAffected %d", rowsAffected)
	}

	return nil
}

// GetCost - Returns total cost of user subscription
func (r *subscriptionRepository) GetCost(ctx context.Context, params entities.FilterParams) (int64, error) {
	query := r.builder.Select(columnsCost...).From(table)
	query = conditionCost(query, params)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, errs.Wrap(err, "subscriptionRepositories.GetCost: build query")
	}

	var cost int64

	err = r.querier.QueryRowxContext(ctx, sql, args...).Scan(&cost)
	if err != nil {
		return 0, errs.Wrap(err, "subscriptionRepositories.GetCost: scan query")
	}
	return cost, nil
}
