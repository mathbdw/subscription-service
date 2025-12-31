package repositories

import "github.com/mathbdw/subscription-service/internal/domain/entities"

// SubscriptionToMap - convert struct Subscription to map
func SubscriptionToMap(subs entities.Subscription) map[string]any {
	data := map[string]any{
		"service_name": subs.ServiceName,
		"user_id":      subs.UserId,
		"price":        subs.Price,
		"start_date":   subs.StartDate,
	}

	if subs.EndDate.Valid {
		data["end_date"] = subs.EndDate.Time.Format("2006-01-02")
	}

	return data
}
