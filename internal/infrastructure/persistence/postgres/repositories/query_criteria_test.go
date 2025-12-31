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

func init (){
	builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}

func TestQueryCriteria_ConditionList(t *testing.T) {
	build := builder.Select("*").From("test")

	from := time.Date(2020, time.January, 15, 14, 30, 0, 0, time.UTC)
	to := time.Date(2020, time.January, 15, 15, 30, 0, 0, time.UTC)

	filter := entities.FilterParams{
		ServiceName: "TestService",
		UserId:      uuid.New(),
		StartDate:   entities.DateRange{From: &from, To: &to},
	}
	build = conditionList(build, filter)

	sql, _, _ := build.ToSql()

	assert.Equal(t, "SELECT * FROM test WHERE service_name LIKE $1 AND user_id = $2 AND start_date >= $3 AND start_date <= $4", sql)
}

func TestQueryCriteria_PaginationList(t *testing.T) {
	tests := []struct {
		name          string
		expectedSql   string
		expectedPage  uint64
		expectedLimit uint64
		totalCount    uint64
		params        *entities.PaginationParams
	}{
		{
			name:          "first_page_no_full",
			expectedSql:   "SELECT * FROM test LIMIT 13 OFFSET 0",
			expectedPage:  1,
			expectedLimit: 13,
			totalCount:    13,
			params:        &entities.PaginationParams{Page: 1, Limit: 20},
		},
		{
			name:          "page_gte_total",
			expectedSql:   "SELECT * FROM test LIMIT 13 OFFSET 0",
			expectedPage:  1,
			expectedLimit: 13,
			totalCount:    13,
			params:        &entities.PaginationParams{Page: 8, Limit: 20},
		},
		{
			name:          "first_page_full",
			expectedSql:   "SELECT * FROM test LIMIT 20 OFFSET 0",
			expectedPage:  1,
			expectedLimit: 20,
			totalCount:    23,
			params:        &entities.PaginationParams{Page: 1, Limit: 20},
		},
		{
			name:          "page_middle_total",
			expectedSql:   "SELECT * FROM test LIMIT 20 OFFSET 20",
			expectedPage:  2,
			expectedLimit: 20,
			totalCount:    68,
			params:        &entities.PaginationParams{Page: 2, Limit: 20},
		},
		{
			name:          "page_eq_total",
			expectedSql:   "SELECT * FROM test LIMIT 20 OFFSET 40",
			expectedPage:  3,
			expectedLimit: 20,
			totalCount:    60,
			params:        &entities.PaginationParams{Page: 3, Limit: 20},
		},
	}

	query := builder.Select("*").From("test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpQuery := paginationList(query, tt.totalCount, tt.params)

			sql, _, err := tmpQuery.ToSql()

			require.NoError(t, err)
			assert.Equal(t, tt.expectedSql, sql)
			assert.Equal(t, tt.expectedPage, tt.params.Page)
			assert.Equal(t, tt.expectedLimit, tt.params.Limit)
		})
	}
}

func TestQueryCriteria_SortListEmpty(t *testing.T) {
	build := builder.Select("*").From("test")

	params := entities.SortParams{}

	build = sortList(build, params)
	sql, _, err := build.ToSql()

	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM test", sql)
}

func TestQueryCriteria_SortList(t *testing.T) {
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
	from := time.Date(2020, time.January, 15, 14, 30, 0, 0, time.UTC)
	to := time.Date(2020, time.January, 15, 15, 30, 0, 0, time.UTC)

	serviceName := "TestService"
	userId := uuid.New()
	startDate := entities.DateRange{From: &from}
	fullDate := entities.DateRange{From: &from, To: &to}

	filter := entities.FilterParams{}

	tests := []struct {
		name           string
		fn             func()
		exepectedQuery string
	}{
		{
			name:           "empty",
			fn:             func() {},
			exepectedQuery: "SELECT * FROM test",
		},
		{
			name:           "WithServiceName",
			fn:             func() { filter.ServiceName = serviceName},
			exepectedQuery: "SELECT * FROM test WHERE service_name = $1",
		},
		{
			name:           "WithServiceNameUserId",
			fn:             func() { filter.UserId = userId },
			exepectedQuery: "SELECT * FROM test WHERE service_name = $1 AND user_id = $2",
		},
		{
			name: "WithServiceNameUserIdDateFrom",
			fn:   func() { filter.StartDate = startDate },
			exepectedQuery: "SELECT * FROM test WHERE service_name = $1 AND user_id = $2 AND start_date >= $3",
		},
		{
			name: "WithServiceNameUserIdDateFromDateTo",
			fn:   func() { filter.StartDate = fullDate },
			exepectedQuery: "SELECT * FROM test WHERE service_name = $1 AND user_id = $2 AND start_date >= $3 AND start_date <= $4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn()
			build := builder.Select("*").From("test")
			build = conditionCost(build, filter)
			sql, _, _ := build.ToSql()

			require.Equal(t, tt.exepectedQuery, sql)
		})
	}
}
