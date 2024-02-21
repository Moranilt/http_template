package response

import (
	"net/http"

	go_json "github.com/goccy/go-json"
)

type DefaultResponse[T any, E any] struct {
	Error E `json:"error"`
	Body  T `json:"body"`
}

func Default[Body any, Err any](w http.ResponseWriter, body Body, responseErr Err, status int) {
	w.WriteHeader(status)
	go_json.NewEncoder(w).Encode(DefaultResponse[Body, Err]{
		Error: responseErr,
		Body:  body,
	})
}

func ErrorResponse[Err any](w http.ResponseWriter, responseErr Err, status int) {
	w.WriteHeader(status)
	go_json.NewEncoder(w).Encode(DefaultResponse[any, Err]{
		Error: responseErr,
		Body:  nil,
	})
}

func SuccessResponse[Body any](w http.ResponseWriter, body Body, status int) {
	w.WriteHeader(status)
	go_json.NewEncoder(w).Encode(DefaultResponse[Body, *string]{
		Error: nil,
		Body:  body,
	})
}
