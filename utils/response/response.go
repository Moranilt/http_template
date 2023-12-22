package response

import (
	"net/http"

	go_json "github.com/goccy/go-json"
)

type DefaultResponse[T any, E any] struct {
	Error E `json:"error"`
	Body  T `json:"body"`
}

func Default(w http.ResponseWriter, body any, responseErr any, status int) {
	w.WriteHeader(status)
	go_json.NewEncoder(w).Encode(DefaultResponse[any, *any]{
		Error: &responseErr,
		Body:  body,
	})
}

func ErrorResponse(w http.ResponseWriter, responseErr any, status int) {
	w.WriteHeader(status)
	go_json.NewEncoder(w).Encode(DefaultResponse[any, any]{
		Error: responseErr,
		Body:  nil,
	})
}

func SuccessResponse(w http.ResponseWriter, body any, status int) {
	w.WriteHeader(status)
	go_json.NewEncoder(w).Encode(DefaultResponse[any, *any]{
		Error: nil,
		Body:  body,
	})
}
