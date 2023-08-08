package response

import (
	"encoding/json"
	"net/http"
)

type DefaultResponse[T any, E any] struct {
	Error E `json:"error"`
	Body  T `json:"body"`
}

func Default(w http.ResponseWriter, body any, responseErr any, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(DefaultResponse[any, *any]{
		Error: &responseErr,
		Body:  body,
	})
}

func ErrorResponse(w http.ResponseWriter, responseErr error, status int) {
	w.WriteHeader(status)
	err := responseErr.Error()
	json.NewEncoder(w).Encode(DefaultResponse[any, string]{
		Error: err,
		Body:  nil,
	})
}

func SuccessResponse(w http.ResponseWriter, body any, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(DefaultResponse[any, *string]{
		Error: nil,
		Body:  body,
	})
}
