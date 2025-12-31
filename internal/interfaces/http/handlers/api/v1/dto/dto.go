package dto

import (
	"github.com/google/uuid"
)

type SubscriptionReq struct {
	ServiceName string    `json:"service_name" validate:"required" example:"TestService"`
	UserId      uuid.UUID `json:"user_id"  validate:"required,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Price       uint32    `json:"price"  validate:"required,gte=1,lte=4294967295" example:"100"`
	StartDate   string    `json:"start_date"  validate:"required,datetime=01-2006" example:"12-2001"`
	EndDate     string    `json:"end_date"  validate:"omitempty,datetime=01-2006" example:"12-2002"`
}

type SubscriptionUpdateReq struct {
	ServiceName string    `json:"service_name" validate:"omitempty" example:"TestService"`
	UserId      uuid.UUID `json:"user_id"  validate:"omitempty,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Price       uint32    `json:"price"  validate:"omitempty,gte=1,lte=4294967295" example:"100"`
	StartDate   string    `json:"start_date"  validate:"omitempty,datetime=01-2006" example:"12-2001"`
	EndDate     string    `json:"end_date"  validate:"omitempty,datetime=01-2006" example:"12-2002"`
}

type SubscriptionResp struct {
	ServiceName string    `json:"service_name"`
	UserId      uuid.UUID `json:"user_id"`
	Price       uint32    `json:"price"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

type QueryParamList struct {
	SortBy    string `form:"sort" query:"sort" validate:"omitempty,sort_by" example:"id"`
	SortOrder string `form:"order" query:"order" validate:"omitempty,sort_order" example:"asc"`

	ServiceName string `form:"service_name" query:"service_name" validate:"omitempty,min=1,max=255" example:"TestService"`
	UserId      string `form:"user_id" query:"user_id" validate:"omitempty,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `form:"start_date" query:"start_date" validate:"omitempty,datetime=01-2006" example:"01-2000"`
	EndDate     string `form:"end_date" query:"end_date" validate:"omitempty,datetime=01-2006" example:"01-2000"`

	Page  int `form:"page" query:"page" validate:"omitempty,gte=1"`
	Limit int `form:"page_size" query:"page_size" validate:"omitempty,gte=1,lte=100"`
}

type QueryParamCost struct {
	ServiceName string `form:"service_name" query:"service_name" validate:"omitempty,min=1,max=255" example:"TestService"`
	UserId      string `form:"user_id" query:"user_id" validate:"omitempty,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `form:"start_date" query:"start_date" validate:"omitempty,datetime=01-2006" example:"01-2000"`
	EndDate     string `form:"end_date" query:"end_date" validate:"omitempty,datetime=01-2006" example:"01-2000"`
}
