package fx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlattenErrorsIfAny(t *testing.T) {
	tests := []struct {
		name     string
		input    []error
		expected error
	}{ // cases
		{
			name:     "no input error should return nil",
			input:    []error{},
			expected: nil,
		},
		{
			name:     "array of nil errors should return nil",
			input:    []error{nil, nil},
			expected: nil,
		},
		{
			name:     "single error should return same error",
			input:    []error{errors.New("expected error")},
			expected: errors.New("[1:expected error]"),
		},
		{
			name:     "multiple errors should return single error containing all sub-errors",
			input:    []error{errors.New("error1"), errors.New("error2")},
			expected: errors.New("[2:error2] [1:error1]"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resErr := FlattenErrorsIfAny(tt.input...)
			if tt.expected != nil {
				require.Error(t, resErr)
				require.EqualError(t, resErr, tt.expected.Error())
			} else {
				require.NoError(t, resErr)
			}
		})
	}
}

func TestFlattenErrorsAsStringIfAny(t *testing.T) {
	tests := []struct {
		name     string
		input    []error
		expected string
	}{ // cases
		{
			name:     "no input error should return nil",
			input:    []error{},
			expected: "",
		},
		{
			name:     "array of nil errors should return nil",
			input:    []error{nil, nil},
			expected: "",
		},
		{
			name:     "single error should return same error",
			input:    []error{errors.New("expected error")},
			expected: errors.New("[1:expected error]").Error(),
		},
		{
			name:     "multiple errors should return single error containing all sub-errors",
			input:    []error{errors.New("error1"), errors.New("error2")},
			expected: errors.New("[2:error2] [1:error1]").Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resErr := FlattenErrorsAsStringIfAny(tt.input...)
			if tt.expected != "" {
				require.Equal(t, resErr, tt.expected)
			} else {
				require.Empty(t, resErr)
			}
		})
	}
}
