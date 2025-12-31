package entities

import (
	"time"

	"github.com/google/uuid"
)

type (
	SortType    string
	SortOrderType string
)

const (
	SortTypeID          SortType = "id"
	SortTypeServiceName SortType = "service_name"

	SortOrderTypeAsc  SortOrderType = "ASC"
	SortOrderTypeDesc SortOrderType = "DESC"
)

type QueryCriteria struct {
	Filter     FilterParams
	Pagination PaginationParams
	Sort       SortParams
}

type FilterParams struct {
	ServiceName string
	UserId      uuid.UUID
	StartDate   DateRange
}

type DateRange struct {
	From *time.Time
	To   *time.Time
}

type PaginationParams struct {
	Page  uint64
	Limit uint64
}

type SortParams struct {
	SortBy    SortType
	SortOrder SortOrderType
}

var SortByTypes = map[string]bool {
	string(SortTypeID): true,
	string(SortTypeServiceName): true,
}

var SortOrderTypes = map[string]bool {
	string(SortOrderTypeAsc): true,
	string(SortOrderTypeDesc): true,
}
