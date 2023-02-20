package http

import (
	"net/http"
	"time"
)

type Client struct {
	useragent string
	client    *http.Client
}

// A ClientOption is a functional option.
type ClientOption func(*clientOptions)

type clientOptions struct {
	clientTimeout time.Duration
	useragent     string
}

// WithClientTimeout sets the client timeout.
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(opts *clientOptions) { opts.clientTimeout = timeout }
}

// WithUserAgent sets the client user-agent.
func WithUserAgent(useragent string) ClientOption {
	return func(opts *clientOptions) { opts.useragent = useragent }
}

// NewClient returns a new http.Client with a default timeout.
func NewClient(options ...ClientOption) *Client {
	opts := clientOptions{clientTimeout: 10 * time.Second}

	for i := range options {
		options[i](&opts)
	}
	client := Client{client: &http.Client{Timeout: opts.clientTimeout}}

	if opts.useragent != "" {
		client.useragent = opts.useragent
	}
	return &client
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.useragent != "" {
		req.Header.Set("User-Agent", c.useragent)
	}
	return c.client.Do(req)
}
