package repositories

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	
	"github.com/mathbdw/subscription-service/internal/domain/entities"
)

// conditionList - SelectBuilder query condition builder for list
func conditionList(query sq.SelectBuilder, params entities.FilterParams) sq.SelectBuilder {
	if params.ServiceName != "" {
		query = query.Where(sq.Like{"service_name": params.ServiceName})
	}

	if params.UserId != uuid.Nil{
		query = query.Where(sq.Eq{"user_id": params.UserId})
	}

	if params.StartDate.From != nil {
		query = query.Where(sq.GtOrEq{"start_date": params.StartDate.From})
	}

	if params.StartDate.To != nil {
		query = query.Where(sq.LtOrEq{"start_date": params.StartDate.To})
	}

	return query
}

// paginationList - SelectBuilder query pagination builder for list
func paginationList(query sq.SelectBuilder, totalCount uint64, params *entities.PaginationParams) sq.SelectBuilder {
	limit := params.Limit
	page := params.Page
	offset := (page - 1) * limit

	if totalCount < page*limit {
		limit = totalCount % limit
		offset = totalCount - limit

		params.Page = offset/params.Limit + 1
		params.Limit = limit
	}
	query = query.Limit(limit).Offset(offset)

	return query
}

// sortList - SelectBuilder query sort builder
func sortList(query sq.SelectBuilder, params entities.SortParams) sq.SelectBuilder {
	if params.SortBy == "" {
		return query
	}

	return query.OrderBy(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder))
}

// conditionCost - SelectBuilder query condition builder for cost
func conditionCost(query sq.SelectBuilder, params entities.FilterParams) sq.SelectBuilder {
	if params.ServiceName != "" {
		query = query.Where(sq.Eq{"service_name": params.ServiceName})
	}

	if params.UserId != uuid.Nil{
		query = query.Where(sq.Eq{"user_id": params.UserId})
	}

	if params.StartDate.From != nil {
		query = query.Where(sq.GtOrEq{"start_date": params.StartDate.From})
	}

	if params.StartDate.To != nil {
		query = query.Where(sq.LtOrEq{"start_date": params.StartDate.To})
	}

	return query
}