package utils

import (
	"testing"

	"github.com/Moranilt/http-utils/tiny_errors"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequiredFields(t *testing.T) {
	tests := []struct {
		name           string
		input          []RequiredField
		expectedOutput []tiny_errors.ErrorOption
	}{
		{
			name: "All fields valid",
			input: []RequiredField{
				{Name: "field1", Value: "value1"},
				{Name: "field2", Value: 42},
				{Name: "field3", Value: 3.14},
			},
			expectedOutput: nil,
		},
		{
			name: "Nil value",
			input: []RequiredField{
				{Name: "field1", Value: nil},
			},
			expectedOutput: []tiny_errors.ErrorOption{tiny_errors.Detail("field1", "required")},
		},
		{
			name: "Empty string",
			input: []RequiredField{
				{Name: "field1", Value: ""},
			},
			expectedOutput: []tiny_errors.ErrorOption{tiny_errors.Detail("field1", "required")},
		},
		{
			name: "Zero int",
			input: []RequiredField{
				{Name: "field1", Value: 0},
			},
			expectedOutput: []tiny_errors.ErrorOption{tiny_errors.Detail("field1", "required")},
		},
		{
			name: "Zero float64",
			input: []RequiredField{
				{Name: "field1", Value: 0.0},
			},
			expectedOutput: []tiny_errors.ErrorOption{tiny_errors.Detail("field1", "required")},
		},
		{
			name: "Nil pointer",
			input: []RequiredField{
				{Name: "field1", Value: (*string)(nil)},
			},
			expectedOutput: []tiny_errors.ErrorOption{tiny_errors.Detail("field1", "required")},
		},
		{
			name: "Mixed valid and invalid fields",
			input: []RequiredField{
				{Name: "field1", Value: "value1"},
				{Name: "field2", Value: ""},
				{Name: "field3", Value: 42},
				{Name: "field4", Value: 0.0},
			},
			expectedOutput: []tiny_errors.ErrorOption{
				tiny_errors.Detail("field2", "required"),
				tiny_errors.Detail("field4", "required"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := ValidateRequiredFields(tt.input...)
			newErr := tiny_errors.New(1, output...)
			expectedErr := tiny_errors.New(1, tt.expectedOutput...)

			assert.EqualValues(t, expectedErr.GetDetails(), newErr.GetDetails())
		})
	}
}
