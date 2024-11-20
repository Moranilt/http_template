package utils

import (
	"reflect"

	"github.com/Moranilt/http-utils/tiny_errors"
)

type RequiredField struct {
	Name  string
	Value any
}

func NewRequiredField(name string, value any) RequiredField {
	return RequiredField{
		Name:  name,
		Value: value,
	}
}

func ValidateRequiredFields(data ...RequiredField) []tiny_errors.ErrorOption {
	var options []tiny_errors.ErrorOption
	for _, field := range data {
		if isEmptyValue(field.Value) {
			options = append(options, tiny_errors.Detail(field.Name, "required"))
		}
	}
	return options
}

func isEmptyValue(v any) bool {
	switch value := v.(type) {
	case nil:
		return true
	case string:
		return value == ""
	case int:
		return value == 0
	case float64:
		return value == 0.0
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Ptr, reflect.Interface:
			return rv.IsNil()
		case reflect.Slice, reflect.Map:
			return rv.Len() == 0
		}
	}
	return false
}
