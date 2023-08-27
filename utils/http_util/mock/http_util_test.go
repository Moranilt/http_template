package http_util_mock

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/Moranilt/http_template/utils/mock"
)

var postClientTests = []struct {
	name          string
	body          []byte
	url           string
	response      *http.Response
	runExpect     bool
	expectedError error
}{
	{
		name: "default call of Post",
		body: []byte(`{"name":"John"}`),
		url:  "http://test.com",
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
		runExpect:     true,
		expectedError: nil,
	},
	{
		name:      "unexpected call of Post",
		body:      []byte(`{"name":"Jane"}`),
		url:       "http://test.com",
		runExpect: false,
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
		expectedError: fmt.Errorf(mock.ERR_Events_Is_Empty, "Post"),
	},
}

var getClientTests = []struct {
	name          string
	url           string
	response      *http.Response
	runExpect     bool
	expectedError error
}{
	{
		name: "default call of Get",
		url:  "http://test.com",
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
		runExpect:     true,
		expectedError: nil,
	},
	{
		name: "unexpected call of Get",
		url:  "http://test.com/users",
		response: &http.Response{
			StatusCode: http.StatusOK,
		},
		runExpect:     false,
		expectedError: fmt.Errorf(mock.ERR_Events_Is_Empty, "Get"),
	},
}

func TestMockHttpUtil(t *testing.T) {

	for _, test := range postClientTests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := NewMockClient(nil, nil)

			if test.runExpect {
				mockClient.ExpectPost(test.url, test.body, nil, test.response)
			}

			resp, err := mockClient.Post(test.url, test.body)
			if err != nil && err.Error() != test.expectedError.Error() {
				t.Errorf("not expected error: %v", err)
			}

			if resp != nil && !reflect.DeepEqual(resp, test.response) {
				t.Errorf("not equal post responses, expect %#v, go %#v", *test.response, *resp)
			}

			if err := mockClient.AllExpectationsDone(); err != nil {
				t.Error(err)
			}
		})
	}

	for _, test := range getClientTests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := NewMockClient(nil, nil)

			if test.runExpect {
				mockClient.ExpectGet(test.url, nil, test.response)
			}

			resp, err := mockClient.Get(test.url)
			if err != nil && err.Error() != test.expectedError.Error() {
				t.Errorf("not expected error: %v", err)
			}

			if resp != nil && !reflect.DeepEqual(resp, test.response) {
				t.Errorf("not equal post responses, expect %#v, go %#v", *test.response, *resp)
			}

			if err := mockClient.AllExpectationsDone(); err != nil {
				t.Error(err)
			}
		})
	}
}

var checkCallTests = []struct {
	name           string
	expectedURL    string
	actualURL      string
	runExpects     bool
	expectedBody   []byte
	unexpectedBody []byte
	expectedError  error
}{
	{
		name:           "not expected call",
		actualURL:      "http://test.com",
		runExpects:     false,
		expectedURL:    "",
		expectedBody:   []byte{},
		unexpectedBody: []byte{},
		expectedError:  fmt.Errorf(mock.ERR_Events_Is_Empty, "Post"),
	},
	{
		name:           "not expected url",
		expectedURL:    "http://test.com",
		actualURL:      "http://test.com/users",
		runExpects:     true,
		expectedBody:   []byte{},
		unexpectedBody: []byte{},
		expectedError:  fmt.Errorf(ERR_Unexpected_Url, "Post", "http://test.com", "http://test.com/users"),
	},
	{
		name:           "not expected body",
		expectedURL:    "http://test.com/users",
		actualURL:      "http://test.com/users",
		runExpects:     true,
		expectedBody:   []byte("expected body"),
		unexpectedBody: []byte("unexpected body"),
		expectedError:  fmt.Errorf(ERR_Unexpected_Data, "Post", "expected body", "unexpected body"),
	},
}

func TestHttpCheckCall(t *testing.T) {
	for _, test := range checkCallTests {
		t.Run(test.name, func(t *testing.T) {
			mockClient := NewMockClient(nil, nil)

			if test.runExpects {
				mockClient.ExpectPost(test.expectedURL, test.expectedBody, nil, &http.Response{
					StatusCode: http.StatusOK,
				})
			}

			item, err := mockClient.checkCall("Post", test.actualURL, test.unexpectedBody)
			if err.Error() != test.expectedError.Error() {
				t.Errorf("got error %q, expected %q", err, test.expectedError)
			}

			if item != nil {
				t.Errorf("expected item to be nil, got %#v", *item)
			}
		})
	}
}

func TestHttpClientReset(t *testing.T) {
	mockClient := NewMockClient(nil, nil)
	mockClient.ExpectGet("http://test.com", nil, &http.Response{StatusCode: http.StatusOK})

	mockClient.Reset()

	if !mockClient.history.IsEmpty() {
		t.Errorf("history is not empty after reset")
	}
}
