package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/Moranilt/http_template/logger"
	response "github.com/Moranilt/http_template/utils/reponse"
	"github.com/gorilla/mux"
)

type mockRequest struct {
	Name    string `json:"name,omitempty" mapstructure:"name"`
	Phone   string `json:"phone,omitempty" mapstructure:"phone"`
	Message string `json:"message,omitempty" mapstructure:"message"`
}

type mockMultipartRequest struct {
	Files []*multipart.FileHeader `json:"files" mapstructure:"files"`
	Name  string                  `json:"name" mapstructure:"name"`
	Age   string                  `json:"age" mapstructure:"age"`
}

type mockResponse struct {
	Info string `json:"info"`
}

var (
	successInfo        = "success info"
	errNameRequired    = "name required"
	errPhoneRequired   = "phone required"
	errMessageRequired = "message required"
)

func makeMockedFunction[ReqT any, RespT any](requestValidator func(request ReqT) RespT, err error) CallerFunc[ReqT, RespT] {
	return func(ctx context.Context, request ReqT) (RespT, error) {
		resp := requestValidator(request)

		return resp, err
	}
}

func TestHandler(t *testing.T) {
	logger := logger.New()
	logger.SetOutput(io.Discard)

	t.Run("default handler Run", func(t *testing.T) {
		routePath := "/test-route"

		testHandler := makeMockedFunction(func(request mockRequest) *mockResponse {
			if request.Name == "" {
				return &mockResponse{
					Info: errNameRequired,
				}
			}

			return &mockResponse{
				Info: successInfo,
			}
		}, nil)

		responseValidator := func(t testing.TB, resp *mockResponse) {
			t.Helper()
			if resp.Info != errNameRequired {
				t.Errorf("should return err message %q, got %q", resp.Info, errNameRequired)
				return
			}

			if resp.Info == successInfo {
				t.Errorf("should not be success response. Got %q, want %q", resp.Info, errPhoneRequired)
				return
			}
		}

		test := newTestHandleController(routePath, testHandler, responseValidator, nil, nil, nil, nil)
		test.Run(t, logger)
	})

	t.Run("Run with JSON request", func(t *testing.T) {
		routePath := "/test-route"

		testHandler := makeMockedFunction(func(request mockRequest) *mockResponse {
			if request.Phone != "phone_number" {
				return &mockResponse{
					Info: errPhoneRequired,
				}
			}

			if request.Name == "" {
				return &mockResponse{
					Info: errNameRequired,
				}
			}

			return &mockResponse{
				Info: successInfo,
			}
		}, nil)

		responseValidator := func(t testing.TB, resp *mockResponse) {
			t.Helper()
			if resp.Info == errPhoneRequired {
				t.Error("phone field is empty")
				return
			}

			if resp.Info == errNameRequired {
				t.Error("name field is empty")
				return
			}

			if resp.Info != successInfo {
				t.Errorf("not valid response, got %q, expected %q", resp.Info, successInfo)
				return
			}
		}

		body, err := json.Marshal(mockRequest{
			Name:  "name",
			Phone: "phone_number",
		})
		if err != nil {
			t.Error(err)
			return
		}

		test := newTestHandleController(routePath, testHandler, responseValidator, body, nil, nil, nil)
		test.Run(t, logger)
	})

	t.Run("Run with Vars request", func(t *testing.T) {
		routePath := "/test-route/{phone}"
		vars := map[string]string{
			"{phone}": "phone_number",
		}

		testHandler := makeMockedFunction(func(request mockRequest) *mockResponse {
			if request.Phone != "phone_number" {
				return &mockResponse{
					Info: errPhoneRequired,
				}
			}

			return &mockResponse{
				Info: successInfo,
			}
		}, nil)

		responseValidator := func(t testing.TB, resp *mockResponse) {
			t.Helper()
			if resp.Info == errPhoneRequired {
				t.Error("phone field is empty")
				return
			}

			if resp.Info != successInfo {
				t.Errorf("not valid response, got %q, expected %q", resp.Info, successInfo)
				return
			}
		}

		test := newTestHandleController(routePath, testHandler, responseValidator, nil, vars, nil, nil)
		test.Run(t, logger)
	})

	t.Run("Run with Query request", func(t *testing.T) {
		routePath := "/test-route"

		testHandler := makeMockedFunction(func(request mockRequest) *mockResponse {
			if request.Phone == "" {
				return &mockResponse{
					Info: errPhoneRequired,
				}
			}

			return &mockResponse{
				Info: successInfo,
			}
		}, nil)

		responseValidator := func(t testing.TB, resp *mockResponse) {
			t.Helper()
			if resp.Info == errPhoneRequired {
				t.Error("phone field is empty")
				return
			}

			if resp.Info != successInfo {
				t.Errorf("not valid response, got %q, expected %q", resp.Info, successInfo)
				return
			}
		}

		query := url.Values{}
		query.Set("phone", "phone_number")

		test := newTestHandleController(routePath, testHandler, responseValidator, nil, nil, query, nil)
		test.Run(t, logger)
	})

	t.Run("Run with Json and Vars request", func(t *testing.T) {
		routePath := "/test-route/{phone}"
		vars := map[string]string{
			"{phone}": "phone_number",
		}

		testHandler := makeMockedFunction(func(request mockRequest) *mockResponse {
			if request.Phone != "phone_number" {
				return &mockResponse{
					Info: errPhoneRequired,
				}
			}

			if request.Name == "" {
				return &mockResponse{
					Info: errNameRequired,
				}
			}

			return &mockResponse{
				Info: successInfo,
			}
		}, nil)

		responseValidator := func(t testing.TB, resp *mockResponse) {
			t.Helper()
			if resp.Info == errPhoneRequired {
				t.Error("phone field is empty")
				return
			}

			if resp.Info == errNameRequired {
				t.Error("name field is empty")
				return
			}

			if resp.Info != successInfo {
				t.Errorf("not valid response, got %q, expected %q", resp.Info, successInfo)
				return
			}
		}

		body, err := json.Marshal(mockRequest{
			Name: "name",
		})
		if err != nil {
			t.Error(err)
			return
		}

		test := newTestHandleController(routePath, testHandler, responseValidator, body, vars, nil, nil)
		test.Run(t, logger)
	})

	t.Run("Run with Json, Vars and Query request", func(t *testing.T) {
		routePath := "/test-route/{phone}"
		vars := map[string]string{
			"{phone}": "phone_number",
		}

		testHandler := makeMockedFunction(func(request mockRequest) *mockResponse {
			if request.Phone != "phone_number" {
				return &mockResponse{
					Info: errPhoneRequired,
				}
			}

			if request.Name == "" {
				return &mockResponse{
					Info: errNameRequired,
				}
			}

			if request.Message == "" {
				return &mockResponse{
					Info: errMessageRequired,
				}
			}

			return &mockResponse{
				Info: successInfo,
			}
		}, nil)

		responseValidator := func(t testing.TB, resp *mockResponse) {
			t.Helper()
			if resp.Info == errPhoneRequired {
				t.Error("phone field is empty")
				return
			}

			if resp.Info == errNameRequired {
				t.Error("name field is empty")
				return
			}

			if resp.Info == errMessageRequired {
				t.Error("message field is empty")
				return
			}

			if resp.Info != successInfo {
				t.Errorf("not valid response, got %q, expected %q", resp.Info, successInfo)
				return
			}
		}

		body, err := json.Marshal(mockRequest{
			Name: "name",
		})
		if err != nil {
			t.Error(err)
			return
		}

		query := url.Values{}
		query.Set("message", "message")

		test := newTestHandleController(routePath, testHandler, responseValidator, body, vars, query, nil)
		test.Run(t, logger)
	})

	t.Run("Run with Multipart", func(t *testing.T) {
		routePath := "/test-files"
		mockedMultipart := mockedMultipartData(t)

		testHandler := makeMockedFunction(func(request mockMultipartRequest) *mockResponse {
			if len(request.Files) == 0 {
				return &mockResponse{
					Info: "files field is empty",
				}
			}

			if request.Name == "" {
				return &mockResponse{
					Info: errNameRequired,
				}
			}

			if request.Age == "" {
				return &mockResponse{
					Info: "age required",
				}
			}

			return &mockResponse{
				Info: successInfo,
			}
		}, nil)

		responseValidator := func(t testing.TB, resp *mockResponse) {
			t.Helper()
			if resp.Info != "" && resp.Info != successInfo {
				t.Error(resp.Info)
				return
			}

			if resp.Info != successInfo {
				t.Errorf("not valid response, got %q, expected %q", resp.Info, successInfo)
				return
			}
		}

		test := newTestHandleController(routePath, testHandler, responseValidator, nil, nil, nil, mockedMultipart)

		test.Run(t, logger)
	})
}

func BenchmarkMultipart(b *testing.B) {
	logger := logger.New()
	logger.SetOutput(io.Discard)
	routePath := "/test-files"

	testHandler := makeMockedFunction(func(request mockMultipartRequest) *mockResponse {
		if len(request.Files) == 0 {
			return &mockResponse{
				Info: "files field is empty",
			}
		}

		if request.Name == "" {
			return &mockResponse{
				Info: errNameRequired,
			}
		}

		if request.Age == "" {
			return &mockResponse{
				Info: "age required",
			}
		}

		return &mockResponse{
			Info: successInfo,
		}
	}, nil)

	router := mux.NewRouter()
	router.HandleFunc(routePath, func(w http.ResponseWriter, r *http.Request) {
		newHandler := New(w, r, logger, testHandler)
		newHandler = newHandler.WithMultipart(32 << 20)
		newHandler.Run(http.StatusOK, http.StatusBadRequest)
	}).Methods(http.MethodPost)

	server := httptest.NewServer(router)
	defer server.Close()
	requestURL, _ := url.Parse(server.URL)
	url := requestURL.JoinPath(routePath)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mockedData := mockedMultipartData(b)
		request, err := http.NewRequest(http.MethodPost, url.String(), mockedData.data)
		if err != nil {
			return
		}
		request.Header.Set("Content-Type", mockedData.header)
		b.StartTimer()
		client := http.Client{}
		_, err = client.Do(request)
		if err != nil {
			b.Error(err)
			return
		}
		b.StopTimer()
	}
}

type multipartData struct {
	data   *bytes.Buffer
	header string
}
type testHandleFuncController[ReqT any, RespT any] struct {
	routePath         string
	handler           CallerFunc[ReqT, RespT]
	responseValidator func(t testing.TB, resp RespT)
	jsonRequest       []byte
	vars              map[string]string
	query             url.Values
	multipartData     *multipartData
}

func newTestHandleController[ReqT any, RespT any](
	routePath string,
	handler CallerFunc[ReqT, RespT],
	responseValidator func(t testing.TB, resp RespT),
	jsonRequest []byte,
	vars map[string]string,
	query url.Values,
	multipartData *multipartData,
) *testHandleFuncController[ReqT, RespT] {
	return &testHandleFuncController[ReqT, RespT]{
		routePath:         routePath,
		handler:           handler,
		responseValidator: responseValidator,
		jsonRequest:       jsonRequest,
		vars:              vars,
		query:             query,
		multipartData:     multipartData,
	}
}

func (cntr *testHandleFuncController[ReqT, RespT]) Run(t testing.TB, logger *logger.Logger) {
	router := mux.NewRouter()
	requestPath := cntr.routePath
	if cntr.vars != nil {
		for key, value := range cntr.vars {
			requestPath = strings.Replace(requestPath, key, value, 1)
		}
	}
	t.Log(cntr.routePath, requestPath)
	router.HandleFunc(cntr.routePath, func(w http.ResponseWriter, r *http.Request) {
		newHandler := New(w, r, logger, cntr.handler)
		if cntr.query != nil {
			newHandler = newHandler.WithQuery()
		}

		if cntr.jsonRequest != nil {
			newHandler = newHandler.WithJson()
		}
		if cntr.vars != nil {
			newHandler = newHandler.WithVars()
		}

		if cntr.multipartData != nil {
			newHandler = newHandler.WithMultipart(32 << 20)
		}

		newHandler.Run(http.StatusOK, http.StatusBadRequest)
	}).Methods(http.MethodPost)

	client := http.Client{}
	server := httptest.NewServer(router)
	defer server.Close()
	requestURL, _ := url.Parse(server.URL)

	url := requestURL.JoinPath(requestPath)
	if cntr.query != nil {
		url.RawQuery = cntr.query.Encode()
	}

	var request *http.Request
	var err error

	if cntr.multipartData != nil {
		request, err = http.NewRequest(http.MethodPost, url.String(), cntr.multipartData.data)
		if err != nil {
			t.Error(err)
			return
		}
		request.Header.Set("Content-Type", cntr.multipartData.header)
	} else {
		request, err = http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(cntr.jsonRequest))
		if err != nil {
			t.Error(err)
			return
		}
	}

	resp, err := client.Do(request)
	if err != nil {
		t.Error(err)
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		t.Error("route not found")
		return
	}

	var mockResp response.DefaultResponse[RespT]
	err = json.NewDecoder(resp.Body).Decode(&mockResp)
	if err != nil {
		t.Error(err)
		return
	}

	cntr.responseValidator(t, mockResp.Body)

}

func mockedMultipartData(t testing.TB) *multipartData {
	t.Helper()
	fields := []struct {
		name  string
		value []byte
	}{
		{
			name:  "name",
			value: []byte("John"),
		},
		{
			name:  "age",
			value: []byte("20"),
		},
	}

	var requestBody bytes.Buffer
	w := multipart.NewWriter(&requestBody)
	defer w.Close()
	for _, field := range fields {
		fw, err := w.CreateFormField(field.name)
		if err != nil {
			t.Errorf("create form field error: %v", err)
			return nil
		}
		fw.Write([]byte(field.value))
	}

	files := []struct {
		name  string
		value string
	}{
		{
			name:  "test_file_1.json",
			value: "{\"name\": \"Elizabeth\"}",
		},
		{
			name:  "test_file_2.json",
			value: "{\"job\": \"Pizzamaker\"}",
		},
	}

	for _, file := range files {
		newFile, err := os.Create(file.name)
		if err != nil {
			t.Errorf("create file: %v", err)
			return nil
		}

		newFile.WriteString(file.value)
		defer os.Remove(newFile.Name())

		fw, err := w.CreateFormFile("files[]", newFile.Name())
		if err != nil {
			t.Errorf("create form file: %v", err)
			return nil
		}
		io.Copy(fw, newFile)
	}

	return &multipartData{
		data:   &requestBody,
		header: w.FormDataContentType(),
	}
}
