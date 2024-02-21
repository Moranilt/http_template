package tiny_errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	t.Run("code and message", func(t *testing.T) {
		err := New(2, Message("error message"))

		assert.Equal(t, 2, err.GetCode())
		assert.Equal(t, "{\"code\":2,\"message\":\"error message\",\"details\":null,\"errors\":[]}", err.JSON())
	})

	t.Run("code, message and details", func(t *testing.T) {
		err := New(2, Message("error message"), Detail("name", "John"))

		assert.Equal(t, 2, err.GetCode())
		assert.Equal(t, "{\"code\":2,\"message\":\"error message\",\"details\":{\"name\":\"John\"},\"errors\":[]}", err.JSON())
	})

	t.Run("error message", func(t *testing.T) {
		err := New(2, Message("error message"), Detail("name", "John"))

		assert.Equal(t, "error message", err.Error())
	})

	t.Run("details map is nil", func(t *testing.T) {
		err := &Error{Code: 1, Message: "error message"}
		option := Detail("name", "John")
		option(err)

		assert.Equal(t, "error message", err.Error())
	})

	t.Run("extra error", func(t *testing.T) {
		secondError := NewMini(3, "error message 2", nil)
		err := New(2, Message("error message"), Err(secondError))

		assert.Equal(t, "error message", err.Error())
		assert.Equal(t, "{\"code\":2,\"message\":\"error message\",\"details\":null,\"errors\":[{\"code\":3,\"message\":\"error message 2\",\"details\":null}]}", err.JSON())
	})

	t.Run("init array of errors", func(t *testing.T) {
		var (
			ErrCodeBodyRequired = 1
			ErrCodeValidation   = 2

			errors = map[int]string{
				ErrCodeBodyRequired: "body required",
				ErrCodeValidation:   "not valid field %s",
			}
		)
		Init(errors)

		t.Run("without message option", func(t *testing.T) {
			err := New(ErrCodeBodyRequired)
			expected := fmt.Sprintf(
				"{\"code\":%d,\"message\":\"%s\",\"details\":null,\"errors\":[]}",
				ErrCodeBodyRequired, errors[ErrCodeBodyRequired],
			)

			assert.Equal(t, expected, err.JSON())
		})

		t.Run("with message option", func(t *testing.T) {
			err := New(ErrCodeBodyRequired, Message("message option"))
			expected := fmt.Sprintf(
				"{\"code\":%d,\"message\":\"message option\",\"details\":null,\"errors\":[]}",
				ErrCodeBodyRequired,
			)

			assert.Equal(t, expected, err.JSON())
		})

		t.Run("mini error without message", func(t *testing.T) {
			err := NewMini(ErrCodeBodyRequired, "", nil)

			assert.Equal(t, ErrCodeBodyRequired, err.Code)
			assert.Equal(t, errors[ErrCodeBodyRequired], err.Message)
		})

		t.Run("mini error with message", func(t *testing.T) {
			err := NewMini(ErrCodeBodyRequired, "error message", nil)

			assert.Equal(t, ErrCodeBodyRequired, err.Code)
			assert.Equal(t, "error message", err.Message)
		})

		t.Run("message args", func(t *testing.T) {
			err := New(ErrCodeValidation, MessageArgs("fieldname"))
			expected := fmt.Sprintf(
				"{\"code\":%d,\"message\":\"%s\",\"details\":null,\"errors\":[]}",
				ErrCodeValidation, fmt.Sprintf(errors[ErrCodeValidation], "fieldname"),
			)

			assert.Equal(t, expected, err.JSON())
		})
	})

}

func BenchmarkHandlerError_GoJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := &Error{
			Code:    400,
			Message: "Bad Request",
			Details: map[string]any{"field": "value"},
			Errors: []*MiniError{
				{
					Code:    400,
					Message: "Bad Request",
					Details: map[string]any{"field": "value"},
				},
			},
		}
		_ = err.JSON()
	}
}

func BenchmarkHandlerError_OriginalJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := &Error{
			Code:    400,
			Message: "Bad Request",
			Details: map[string]any{"field": "value"},
			Errors: []*MiniError{
				{
					Code:    400,
					Message: "Bad Request",
					Details: map[string]any{"field": "value"},
				},
			},
		}
		_ = err.JSONOrigin()
	}
}
