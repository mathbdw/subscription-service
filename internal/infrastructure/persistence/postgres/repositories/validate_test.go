package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate_UpdateFields(t *testing.T) {
	tests := []struct {
		name                string
		expectedError       bool
		expectedErrorString string
		fields              map[string]any
	}{
		{
			name:                "name_field_error",
			expectedError:       true,
			expectedErrorString: "field not found",
			fields: map[string]any{
				"not_found": 1,
			},
		},
		{
			name:                "value_field_error",
			expectedError:       true,
			expectedErrorString: "invalid value for field",
			fields: map[string]any{
				"service_name": uint32(4),
			},
		},
		{
			name:                "value_success",
			expectedError:       false,
			expectedErrorString: "",
			fields: map[string]any{
				"service_name": "Test Service",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUpdateFields(tt.fields)

			if tt.expectedError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorString)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
