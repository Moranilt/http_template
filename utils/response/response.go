package response

import (
	"encoding/json"
	"net/http"
)

type DefaultResponse[T any] struct {
	Error *string `json:"error"`
	Body  T       `json:"body"`
}

func ErrorResponse(w http.ResponseWriter, responseErr error, status int) {
	w.WriteHeader(status)
	err := responseErr.Error()
	json.NewEncoder(w).Encode(DefaultResponse[any]{
		Error: &err,
		Body:  nil,
	})
}

func SuccessResponse(w http.ResponseWriter, body any, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(DefaultResponse[any]{
		Error: nil,
		Body:  body,
	})
}
