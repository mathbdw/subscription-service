package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int64     `db:"id"`
	ServiceName string    `db:"service_name"`
	UserId      uuid.UUID `db:"user_id"`
	Price       uint32    `db:"price"`
	StartDate   time.Time `db:"start_date"`
	EndDate     sql.NullTime `db:"end_date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type PaginationInfo struct {
	Page       uint64
	PageSize   uint16
	TotalCount uint64
	TotalPages uint32
}

type ResponseListSubscription struct {
	Data []Subscription
	Info PaginationInfo
}

var SubscriptionUpdateFields = map[string]func(value any) bool{
	"service_name": isString,
	"user_id":      isUUID,
	"price":        isUint32,
	"start_date":   isTime,
	"end_date":     isTime,
	"updated_at":   isTime,
}

func isString(value any) bool {
	_, ok := value.(string)
	return ok
}

func isUUID(value any) bool {
	_, ok := value.(uuid.UUID)
	return ok
}

func isUint32(value any) bool {
	_, ok := value.(uint32)
	return ok
}

func isTime(value any) bool {
	_, ok := value.(time.Time)
	return ok
}
