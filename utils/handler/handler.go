package handler

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/Moranilt/http_template/logger"
	"github.com/Moranilt/http_template/utils/response"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

const (
	ErrNotValidBodyFormat = "unable to unmarshal request body "
)

type HandlerMaker[ReqT any, RespT any] struct {
	request     *http.Request
	response    http.ResponseWriter
	requestBody ReqT
	logger      *logrus.Entry
	caller      CallerFunc[ReqT, RespT]
}

type CallerFunc[ReqT any, RespT any] func(ctx context.Context, req ReqT) (RespT, error)

func New[ReqT any, RespT any](w http.ResponseWriter, r *http.Request, logger *logger.Logger, caller CallerFunc[ReqT, RespT]) *HandlerMaker[ReqT, RespT] {
	log := logger.WithRequestInfo(r)
	return &HandlerMaker[ReqT, RespT]{
		logger:   log,
		request:  r,
		caller:   caller,
		response: w,
	}
}

// Request type should include fields with tags of json
//
// Example:
//
//	type YourRequest struct {
//			FieldName string `json:"field_name"`
//	}
func (h *HandlerMaker[ReqT, RespT]) WithJson() *HandlerMaker[ReqT, RespT] {
	if h.request.Method == http.MethodGet {
		return h
	}
	err := json.NewDecoder(h.request.Body).Decode(&h.requestBody)
	if err != nil {
		h.logger.Error(ErrNotValidBodyFormat, err)
		return h
	}
	return h
}

// Request type should include fields with tags of mapstructure
//
// Example:
//
//	type YourRequest struct {
//			FieldName string `mapstructure:"field_name"`
//	}
func (h *HandlerMaker[ReqT, RespT]) WithVars() *HandlerMaker[ReqT, RespT] {
	vars := mux.Vars(h.request)
	err := mapstructure.Decode(vars, &h.requestBody)
	if err != nil {
		h.logger.Error(ErrNotValidBodyFormat, err)
		return h
	}
	return h
}

// Request type should include fields with tags of mapstructure
//
// Example:
//
//	type YourRequest struct {
//			FieldName string `mapstructure:"field_name"`
//	}
func (h *HandlerMaker[ReqT, RespT]) WithQuery() *HandlerMaker[ReqT, RespT] {
	query := h.request.URL.Query()
	if len(query) == 0 {
		return h
	}

	queryVars := make(map[string]any)
	for name, q := range query {
		queryVars[name] = q[0]
	}
	err := mapstructure.Decode(queryVars, &h.requestBody)
	if err != nil {
		h.logger.Error(ErrNotValidBodyFormat, err)
		return h
	}
	return h
}

// Request type should include fields with tags of mapstructure
//
// If field is an array of files you should set tag name as files[] and type []*multipart.FileHeader([mime/multipart.FileHeader])
//
// If field is file and not array of files you should set tag with field name without brackets and type *multipart.FileHeader([mime/multipart.FileHeader])
//
// Other fields should have string type([mime/multipart.Form])
//
// # File types
//
//   - []*multipart.FileHeader -	field with array of files. Should contain square brackets in name
//   - *multipart.FileHeader -	field with single file. Should not contain square brackets in field name
//
// # Example
//
//	type YourRequest struct {
//		MultipleFiles []*multipart.FileHeader `mapstructure:"your_files[]"`
//		SingleFile *multipart.FileHeader 	`mapstructure:"single_file"`
//		Name string `mapstructure:"name"`
//	}
func (h *HandlerMaker[ReqT, RespT]) WithMultipart(maxMemory int64) *HandlerMaker[ReqT, RespT] {
	err := h.request.ParseMultipartForm(maxMemory)
	if err != nil {
		h.logger.Error(ErrNotValidBodyFormat, err)
		return h
	}

	result := make(map[string]any, len(h.request.MultipartForm.Value)+len(h.request.MultipartForm.File))
	for name, value := range h.request.MultipartForm.Value {
		if len(value) > 0 {
			result[name] = value[0]
		}
	}

	for name, value := range h.request.MultipartForm.File {
		if len(value) > 0 {
			fieldName, validName := extractArrayName(name)
			if validName {
				safeValue := make([]*multipart.FileHeader, 0)
				safeValue = append(safeValue, value...)
				result[fieldName] = safeValue
			} else {
				result[name] = value[0]
			}
		}
	}

	err = mapstructure.Decode(result, &h.requestBody)
	if err != nil {
		h.logger.Error(ErrNotValidBodyFormat, err)
		return h
	}

	return h
}

func (h *HandlerMaker[ReqT, RespT]) Run(successStatus, failedStatus int) {
	h.logger.WithField("body", h.requestBody).Debug("request")

	resp, err := h.caller(h.request.Context(), h.requestBody)
	if err != nil {
		h.logger.Error(err)
		response.ErrorResponse(h.response, err, failedStatus)
		return
	}
	response.SuccessResponse(h.response, resp, successStatus)
}
