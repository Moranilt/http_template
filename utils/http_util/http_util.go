package http_util

import (
	"bytes"
	"context"
	"net/http"
	"sync/atomic"
	"time"
)

var timeout atomic.Value

func init() {
	timeout.Store(10 * time.Second)
}

// Set request timeout
func SetTimeout(val time.Duration) {
	timeout.Store(val)
}

// Get request timeout
func Timeout() time.Duration {
	return timeout.Load().(time.Duration)
}

type client struct {
	client *http.Client
}

type Client interface {
	Post(url string, body []byte) (*http.Response, error)
}

func New() Client {
	return &client{
		client: &http.Client{},
	}
}

func (c *client) Post(url string, body []byte) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	request, cancel := c.setRequestTimeout(request)
	defer cancel()
	res, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) Get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request, cancel := c.setRequestTimeout(request)
	defer cancel()
	res, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *client) setRequestTimeout(req *http.Request) (*http.Request, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(req.Context(), Timeout())
	return req.WithContext(ctx), cancel
}
