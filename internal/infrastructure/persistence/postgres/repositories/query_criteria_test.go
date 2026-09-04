package repositories

import (
	"fmt"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
)

var builder sq.StatementBuilderType

func init() {
	builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func TestQueryCriteria_ConditionList(t *testing.T) {
	t.Parallel()
	build := builder.Select("*").From("test")

	from := time.Date(2020, time.January, 15, 14, 30, 0, 0, time.UTC)
	to := time.Date(2020, time.January, 15, 15, 30, 0, 0, time.UTC)

	filter := entities.FilterParams{
		ServiceName: "TestService",
		UserID:      uuid.New(),
		StartDate:   entities.DateRange{From: &from, To: &to},
	}
	build = conditionList(build, filter)

	sql, _, err := build.ToSql()
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM test WHERE service_name LIKE $1 AND user_id = $2 AND start_date >= $3 AND start_date <= $4", sql)
}

func TestQueryCriteria_PaginationList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		expectedSQL   string
		expectedPage  uint64
		expectedLimit uint64
		totalCount    uint64
		params        *entities.PaginationParams
	}{
		{
			name:          "first_page_no_full",
			expectedSQL:   "SELECT * FROM test LIMIT 13 OFFSET 0",
			expectedPage:  1,
			expectedLimit: 13,
			totalCount:    13,
			params:        &entities.PaginationParams{Page: 1, Limit: 20},
		},
		{
			name:          "page_gte_total",
			expectedSQL:   "SELECT * FROM test LIMIT 13 OFFSET 0",
			expectedPage:  1,
			expectedLimit: 13,
			totalCount:    13,
			params:        &entities.PaginationParams{Page: 8, Limit: 20},
		},
		{
			name:          "first_page_full",
			expectedSQL:   "SELECT * FROM test LIMIT 20 OFFSET 0",
			expectedPage:  1,
			expectedLimit: 20,
			totalCount:    23,
			params:        &entities.PaginationParams{Page: 1, Limit: 20},
		},
		{
			name:          "page_middle_total",
			expectedSQL:   "SELECT * FROM test LIMIT 20 OFFSET 20",
			expectedPage:  2,
			expectedLimit: 20,
			totalCount:    68,
			params:        &entities.PaginationParams{Page: 2, Limit: 20},
		},
		{
			name:          "page_eq_total",
			expectedSQL:   "SELECT * FROM test LIMIT 20 OFFSET 40",
			expectedPage:  3,
			expectedLimit: 20,
			totalCount:    60,
			params:        &entities.PaginationParams{Page: 3, Limit: 20},
		},
	}

	query := builder.Select("*").From("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tmpQuery := paginationList(query, tt.totalCount, tt.params)

			sql, _, err := tmpQuery.ToSql()

			require.NoError(t, err)
			assert.Equal(t, tt.expectedSQL, sql)
			assert.Equal(t, tt.expectedPage, tt.params.Page)
			assert.Equal(t, tt.expectedLimit, tt.params.Limit)
		})
	}
}

func TestQueryCriteria_SortListEmpty(t *testing.T) {
	t.Parallel()
	build := builder.Select("*").From("test")

	params := entities.SortParams{}

	build = sortList(build, params)
	sql, _, err := build.ToSql()

	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM test", sql)
}

func TestQueryCriteria_SortList(t *testing.T) {
	t.Parallel()
	build := builder.Select("*").From("test")

	params := entities.SortParams{
		SortBy:    entities.SortType(entities.SortTypeID),
		SortOrder: entities.SortOrderTypeDesc,
	}

	build = sortList(build, params)

	sql, _, err := build.ToSql()

	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("SELECT * FROM test ORDER BY %s %s", params.SortBy, params.SortOrder), sql)
}

func TestQueryCriteria_ConditionCost(t *testing.T) {
	t.Parallel()
	from := time.Date(2020, time.January, 15, 14, 30, 0, 0, time.UTC)
	to := time.Date(2020, time.January, 15, 15, 30, 0, 0, time.UTC)

	serviceName := "TestService"
	userID := uuid.New()
	startDate := entities.DateRange{From: &from}
	fullDate := entities.DateRange{From: &from, To: &to}

	tests := []struct {
		name           string
		filter         entities.FilterParams
		exepectedQuery string
	}{
		{
			name:           "empty",
			filter:         entities.FilterParams{},
			exepectedQuery: "SELECT * FROM test",
		},
		{
			name: "WithServiceName",
			filter: entities.FilterParams{
				ServiceName: serviceName,
			},
			exepectedQuery: "SELECT * FROM test WHERE service_name = $1",
		},
		{
			name: "WithServiceNameUserID",
			filter: entities.FilterParams{
				ServiceName: serviceName,
				UserID:      userID,
			},
			exepectedQuery: "SELECT * FROM test WHERE service_name = $1 AND user_id = $2",
		},
		{
			name: "WithServiceNameUserIDDateFrom",
			filter: entities.FilterParams{
				ServiceName: serviceName,
				UserID:      userID,
				StartDate:   startDate,
			},
			exepectedQuery: "SELECT * FROM test WHERE service_name = $1 AND user_id = $2 AND start_date >= $3",
		},
		{
			name: "WithServiceNameUserIDDateFromDateTo",
			filter: entities.FilterParams{
				ServiceName: serviceName,
				UserID:      userID,
				StartDate:   fullDate,
			},
			exepectedQuery: "SELECT * FROM test WHERE service_name = $1 AND user_id = $2 AND start_date >= $3 AND start_date <= $4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			build := builder.Select("*").From("test")
			build = conditionCost(build, tt.filter)
			sql, _, err := build.ToSql()

			require.NoError(t, err)
			require.Equal(t, tt.exepectedQuery, sql)
		})
	}
}
