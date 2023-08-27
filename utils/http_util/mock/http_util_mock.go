package http_util_mock

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"

	"github.com/Moranilt/http_template/utils/mock"
)

const (
	ERR_Unexpected_Url  = "call %q expected url %q, got %q"
	ERR_Unexpected_Data = "call %q expected data %v, got %v"
)

type MockedClient struct {
	history        *mock.MockHistory[*mockClientData]
	wg             *sync.WaitGroup
	expectCallback func()
	actualCallback func()
}

type mockClientData struct {
	url      string
	body     []byte
	response *http.Response
	err      error
}

func NewMockClient(expectCallback func(), actualCallback func()) *MockedClient {
	return &MockedClient{
		history:        mock.NewMockHistory[*mockClientData](),
		wg:             &sync.WaitGroup{},
		expectCallback: expectCallback,
		actualCallback: actualCallback,
	}
}

func (m *MockedClient) ExpectPost(url string, body []byte, err error, response *http.Response) {
	if m.expectCallback != nil {
		m.expectCallback()
	}
	m.history.Push("Post", &mockClientData{
		url:      url,
		body:     body,
		response: response,
		err:      err,
	}, err)
}

func (m *MockedClient) ExpectGet(url string, err error, response *http.Response) {
	if m.expectCallback != nil {
		m.expectCallback()
	}
	m.history.Push("Get", &mockClientData{
		url:      url,
		body:     nil,
		response: response,
		err:      err,
	}, err)
}

func (m *MockedClient) AllExpectationsDone() error {
	return m.history.AllExpectationsDone()
}

func (m *MockedClient) Reset() {
	m.history.Clear()
}

func (m *MockedClient) Post(url string, body []byte) (*http.Response, error) {
	if m.actualCallback != nil {
		m.actualCallback()
	}
	item, err := m.checkCall("Post", url, body)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

func (m *MockedClient) Get(url string) (*http.Response, error) {
	if m.actualCallback != nil {
		m.actualCallback()
	}
	item, err := m.checkCall("Get", url, nil)
	if err != nil {
		return nil, err
	}

	return item.Data.response, item.Data.err
}

func (m *MockedClient) checkCall(name string, url string, body []byte) (*mock.MockHistoryItem[*mockClientData], error) {
	item, err := m.history.Get(name)
	if err != nil {
		return nil, err
	}

	if item.Data.url != url {
		return nil, fmt.Errorf(ERR_Unexpected_Url, name, item.Data.url, url)
	}

	if !reflect.DeepEqual(item.Data.body, body) {
		return nil, fmt.Errorf(ERR_Unexpected_Data, name, string(item.Data.body), string(body))
	}

	return item, nil
}
