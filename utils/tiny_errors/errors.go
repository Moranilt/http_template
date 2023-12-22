package tiny_errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	go_json "github.com/goccy/go-json"
)

var errorStorage atomic.Value

func init() {
	errorStorage.Store(make(map[int]string))
}

func Init(errors map[int]string) {
	errorStorage.Store(errors)
}

func ErrorStorage() map[int]string {
	return errorStorage.Load().(map[int]string)
}

type Error struct {
	HTTPStatus  int            `json:"http_status"`
	HTTPMessage string         `json:"http_message"`
	Code        int            `json:"code"`
	Message     string         `json:"message"`
	Details     map[string]any `json:"details"`
	Errors      []*MiniError   `json:"errors"`
}

type MiniError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details"`
}

type ErrorHandler interface {
	JSON() string
	JSONOrigin() string
	Error() string

	GetHTTPStatus() int
	GetHTTPMessage() string
	GetCode() int
	GetMessage() string
}

type PropertySetter interface {
	SetMessage(string, ...any)
	FormatMessage(...any)
	SetCode(int)
	SetDetail(name, data string)
	SetError(*MiniError)
	SetHTTPStatus(int)
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) JSON() string {
	b, _ := go_json.Marshal(e)
	return string(b)
}

func (e *Error) JSONOrigin() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e *Error) GetMessage() string {
	return e.Message
}

func (e *Error) GetCode() int {
	return e.Code
}

func (e *Error) GetHTTPStatus() int {
	return e.HTTPStatus
}
func (e *Error) GetHTTPMessage() string {
	return e.HTTPMessage
}

func (e *Error) SetCode(code int) {
	e.Code = code
}

func (e *Error) SetHTTPStatus(code int) {
	e.HTTPStatus = code
	e.HTTPMessage = http.StatusText(code)
}

func (e *Error) SetMessage(msg string, format ...any) {
	if len(format) > 0 {
		msg = fmt.Sprintf(msg, format)
	}
	e.Message = msg
}

func (e *Error) SetError(err *MiniError) {
	e.Errors = append(e.Errors, err)
}

func (e *Error) SetDetail(name, data string) {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[name] = data
}

func (e *Error) FormatMessage(args ...any) {
	e.Message = fmt.Sprintf(e.Message, args...)
}

type ErrorOption func(PropertySetter)

func Message(message string, format ...any) ErrorOption {
	return func(err PropertySetter) {
		err.SetMessage(message, format...)
	}
}

func MessageArgs(args ...any) ErrorOption {
	return func(err PropertySetter) {
		err.FormatMessage(args...)
	}
}

func Detail(name, data string) ErrorOption {
	return func(err PropertySetter) {
		err.SetDetail(name, data)
	}
}

func HTTPStatus(status int) ErrorOption {
	return func(err PropertySetter) {
		err.SetHTTPStatus(status)
	}
}

func Err(miniErr *MiniError) ErrorOption {
	return func(err PropertySetter) {
		err.SetError(miniErr)
	}
}

func New(code int, options ...ErrorOption) ErrorHandler {
	err := &Error{
		HTTPStatus:  http.StatusBadRequest,
		HTTPMessage: http.StatusText(http.StatusBadRequest),
		Code:        code,
		Errors:      []*MiniError{},
	}

	if text, ok := ErrorStorage()[code]; ok {
		err.Message = text
	}

	for _, opt := range options {
		opt(err)
	}

	return err
}

func NewMini(code int, message string, details map[string]any) *MiniError {
	err := &MiniError{
		Code:    code,
		Details: details,
	}

	if len(message) > 0 {
		err.Message = message
	} else if text, ok := ErrorStorage()[code]; ok {
		err.Message = text
	}

	return err
}
