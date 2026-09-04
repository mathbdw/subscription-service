package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/mathbdw/subscription-service/internal/domain/entities"
)

func TestUser_SubscriptionToMap(t *testing.T) {
	t.Parallel()
	subsTest := entities.Subscription{
		ServiceName: "test service",
		UserID:      uuid.New(),
		Price:       100,
		StartDate:   time.Now(),
	}

	tests := []struct {
		name      string
		endTime   sql.NullTime
		expectLen int
	}{
		{
			name:      "withoutEndTime",
			expectLen: 4,
		},
		{
			name:      "withEndTime",
			expectLen: 5,
			endTime:   sql.NullTime{Time: time.Now(), Valid: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sub := subsTest

			if tt.expectLen > 4 {
				sub.EndDate = tt.endTime
			}

			resMap := SubscriptionToMap(sub)

			require.Len(t, resMap, tt.expectLen)
		})
	}
}
