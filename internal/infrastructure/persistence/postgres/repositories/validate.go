package repositories

import (
	"fmt"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
)

func validateUpdateFields(fields map[string]any) error {
	for key, value := range fields {
		validator, ok := entities.SubscriptionUpdateFields[key]
		if !ok {
			return fmt.Errorf("validateReposiroties.UpdateFields: field not found - %s", key)
		}

		if !validator(value) {
			return fmt.Errorf("validateReposiroties.UpdateFields: invalid value for field %s", key)
		}
	}

	return nil
}
