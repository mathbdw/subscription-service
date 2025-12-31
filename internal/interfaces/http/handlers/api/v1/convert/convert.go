package convert

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
	"github.com/mathbdw/subscription-service/internal/interfaces/http/handlers/api/v1/dto"
)

func SubscriptionRequestToEntity(req dto.SubscriptionReq) (entities.Subscription, error) {
	sub := entities.Subscription{
		ServiceName: req.ServiceName,
		UserId:      req.UserId,
		Price:       req.Price,
	}
	tmpDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return entities.Subscription{}, fmt.Errorf("StartDate parse - %s", req.StartDate)
	}
	sub.StartDate = tmpDate

	if req.EndDate != "" {
		tmpDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return entities.Subscription{}, fmt.Errorf("EndDate parse - %s", req.EndDate)
		}

		sub.EndDate = sql.NullTime{
			Time:  tmpDate,
			Valid: true,
		}
	}

	return sub, nil
}

func SubscriptionRequestToMap(req dto.SubscriptionUpdateReq) (map[string]any, error) {
	dataMap := map[string]any{
		"service_name": req.ServiceName,
		"user_id":      req.UserId,
		"price":        req.Price,
	}
	tmpDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("StartDate parse - %s", req.StartDate)
	}
	dataMap["start_date"] = tmpDate

	if req.EndDate != "" {
		tmpDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("EndDate parse - %s", req.EndDate)
		}
		dataMap["end_date"] = tmpDate
	}

	return dataMap, nil
}

func SubscriptionEntityToResponse(entity entities.Subscription) dto.SubscriptionResp {

	resp := dto.SubscriptionResp{
		ServiceName: entity.ServiceName,
		UserId:      entity.UserId,
		Price:       entity.Price,
		StartDate:   entity.StartDate.Format("01-2006"),
	}

	if entity.EndDate.Valid {
		resp.EndDate = entity.EndDate.Time.Format("01-2006")
	}

	return resp
}

func SubscriptionListToResponse(subs []entities.Subscription) []dto.SubscriptionResp {
	resp := make([]dto.SubscriptionResp, 0, len(subs))
	for _, sub := range subs {
		resp = append(resp, SubscriptionEntityToResponse(sub))
	}

	return resp
}

func SubscriptionQueryParamsToQueryCriteria(params dto.QueryParamList) (*entities.QueryCriteria, error) {
	var (
		queryCriteria entities.QueryCriteria
		tmpUUID       uuid.UUID
		err           error
	)
	if params.SortBy == "" {
		queryCriteria.Sort.SortBy = entities.SortTypeID
	} else {
		queryCriteria.Sort.SortBy = entities.SortType(params.SortBy)
	}

	if params.SortOrder == "" {
		queryCriteria.Sort.SortOrder = entities.SortOrderTypeAsc
	} else {
		queryCriteria.Sort.SortOrder = entities.SortOrderType(params.SortOrder)
	}

	if params.Page == 0 {
		queryCriteria.Pagination.Page = uint64(1)
	} else {
		queryCriteria.Pagination.Page = uint64(params.Page)
	}

	if params.Limit == 0 {
		queryCriteria.Pagination.Limit = uint64(20)
	} else {
		queryCriteria.Pagination.Limit = uint64(params.Limit)
	}

	if params.ServiceName != "" {
		queryCriteria.Filter.ServiceName = params.ServiceName
	}

	if params.UserId != "" {
		tmpUUID, err = uuid.Parse(params.UserId)
		if err != nil {
			return nil, fmt.Errorf("UUID parse - %s", params.UserId)
		}

		queryCriteria.Filter.UserId = tmpUUID
	}

	if params.StartDate != "" {
		tmpTime, err := time.Parse("01-2006", params.StartDate)
		if err != nil {
			return nil, fmt.Errorf("StartDate parse - %s", params.StartDate)
		}

		queryCriteria.Filter.StartDate.From = &tmpTime
	}
	if params.EndDate != "" {
		tmpTime, err := time.Parse("01-2006", params.EndDate)
		if err != nil {
			return nil, fmt.Errorf("EndDate parse - %s", params.EndDate)
		}

		queryCriteria.Filter.StartDate.To = &tmpTime
	}

	return &queryCriteria, nil
}

func SubscriptionQueryParamsCostToFilterParam(params dto.QueryParamCost) (entities.FilterParams, error) {
	var (
		filter entities.FilterParams
		tmpUUID       uuid.UUID
		err           error
	)

	if params.ServiceName != "" {
		filter.ServiceName = params.ServiceName
	}

	if params.UserId != "" {
		tmpUUID, err = uuid.Parse(params.UserId)
		if err != nil {
			return entities.FilterParams{}, fmt.Errorf("UUID parse - %s", params.UserId)
		}

		filter.UserId = tmpUUID
	}

	if params.StartDate != "" {
		tmpTime, err := time.Parse("01-2006", params.StartDate)
		if err != nil {
			return entities.FilterParams{}, fmt.Errorf("StartDate parse - %s", params.StartDate)
		}

		filter.StartDate.From = &tmpTime
	}
	if params.EndDate != "" {
		tmpTime, err := time.Parse("01-2006", params.EndDate)
		if err != nil {
			return entities.FilterParams{}, fmt.Errorf("EndDate parse - %s", params.EndDate)
		}

		filter.StartDate.To = &tmpTime
	}

	return filter, nil
}
