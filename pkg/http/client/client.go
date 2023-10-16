package client

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type client struct {
	method  string
	url     string
	headers map[string]string
	body    io.Reader
	ctx     context.Context
}

type OptFunc func(*client)

type Response struct {
	data []byte
}

func (r *Response) String() string {
	return string(r.data)
}

func Get(url string, options ...OptFunc) (*Response, error) {
	c := &client{
		method:  http.MethodGet,
		url:     url,
		headers: make(map[string]string),
		ctx:     context.Background(),
	}

	for _, option := range options {
		option(c)
	}

	return exec(c)
}

func Post(url string, body io.Reader, options ...OptFunc) (*Response, error) {
	c := &client{
		method:  http.MethodPost,
		url:     url,
		headers: make(map[string]string),
		body:    body,
		ctx:     context.Background(),
	}

	for _, option := range options {
		option(c)
	}

	return exec(c)
}

func exec(c *client) (*Response, error) {
	req, err := http.NewRequestWithContext(c.ctx, c.method, c.url, c.body)
	if err != nil {
		return nil, err
	}

	for key, value := range c.headers {
		req.Header.Add(key, value)
	}

	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{data: data}, nil
}

func WithHeaders(headers map[string]string) OptFunc {
	return func(c *client) {
		c.headers = headers
	}
}

func WithContext(ctx context.Context) OptFunc {
	return func(c *client) {
		c.ctx = ctx
	}
}

func FormBody(formData url.Values) io.Reader {
	encodedForm := formData.Encode()
	return strings.NewReader(encodedForm)
}
